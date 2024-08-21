package domaincrawler_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/projectdiscovery/alterx"

	httpx "github.com/projectdiscovery/httpx/runner"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

type Domains struct {
	URL   string `json:"url"`
	Label string `json:"label"`
	Class string `json:"class"`
}

type SubDomainPrefixTrend struct {
	URL       string `json:"url"`
	SubDomain string `json:"sub_domain"`
	Type      string `json:"type"`
}

func TestDomainAutoCrawler(t *testing.T) {
	dslfilePath := "dsl_output.txt"
	subDomainsfilePath := "sub_domains.txt"
	// todo: (enka) ここはクローラーで無限生成
	targetHost := "runetale.com"

	// ** subdomain
	subfinderOpts := &runner.Options{
		Threads:            10, // Thread controls the number of threads to use for active enumerations
		Timeout:            30, // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
	}

	log.SetFlags(0)

	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		t.Fatalf("failed to create subfinder runner: %v", err)
	}

	subfinderOutput := &bytes.Buffer{}
	if err = subfinder.EnumerateSingleDomainWithCtx(context.Background(), targetHost, []io.Writer{subfinderOutput}); err != nil {
		t.Fatalf("failed to enumerate single domain: %v", err)
	}

	subdomainReader := io.Reader(subfinderOutput)
	var sdOutput bytes.Buffer
	io.Copy(&sdOutput, subdomainReader)

	sdFiles, err := os.Create(subDomainsfilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer sdFiles.Close()

	_, err = sdOutput.WriteTo(sdFiles)
	if err != nil {
		t.Fatal(err)
	}

	sd, err := os.Open(subDomainsfilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer sd.Close()

	var targetHostWithSubDomains []string
	scanner := bufio.NewScanner(sd)
	for scanner.Scan() {
		line := scanner.Text()
		targetHostWithSubDomains = append(targetHostWithSubDomains, line)
	}

	// ** alterx
	// 存在するドメイン+サブドメインを使用して、予測ドメインを算出
	// mlに使用できそ
	opts := &alterx.Options{
		// todo: (enka) ここはクローラーでも無限生成
		Domains: targetHostWithSubDomains,
		MaxSize: math.MaxInt,
	}

	m, err := alterx.New(opts)
	if err != nil {
		t.Fatal(err)
	}

	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	m.ExecuteWithWriter(writer)

	reader := io.Reader(&buffer)

	var output bytes.Buffer
	io.Copy(&output, reader)

	file, err := os.Create(dslfilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = output.WriteTo(file)
	if err != nil {
		t.Fatal(err)
	}

	dslfile, err := os.Open(dslfilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer dslfile.Close()

	// ** sub domainのprefixを取得
	// 今後の学習に活かせそうなので、jsonとしてアウトプットしておく
	subDomainsPrefix, err := getSubdomainPrefix(dslfilePath)
	if err != nil {
		t.Fatal(err)
	}

	var trendSubDomains []SubDomainPrefixTrend
	for rootDomain, subs := range subDomainsPrefix {
		for _, sub := range subs {
			// todo: (enka) このprefixをいくつか用意してもらって、サブドメインのサービスタイプの傾向を学習できそう
			if strings.Contains(sub, "api") {
				trendSubDomains = append(trendSubDomains, SubDomainPrefixTrend{
					URL:       "https://" + sub + "." + rootDomain,
					SubDomain: sub,
					Type:      "api server",
				})
			}

			if strings.Contains(sub, "sql") {
				trendSubDomains = append(trendSubDomains, SubDomainPrefixTrend{
					URL:       "https://" + sub + "." + rootDomain,
					SubDomain: sub,
					Type:      "database",
				})
			}
		}
	}
	trendSubDomainFilePath := "trend_subdomains.json"
	trendoutput, err := os.Create(trendSubDomainFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer trendoutput.Close()

	encoder := json.NewEncoder(trendoutput)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(trendSubDomains); err != nil {
		t.Fatal(err)
	}

	// 存在するドメインの情報を収集し、jsonで保存
	err = saveDomainsInfoJson(targetHostWithSubDomains, "mark")
	if err != nil {
		t.Fatal(err)
	}

	// 収集したドメインの情報を元にドメインの性質を算出
	err = createDomainsJson()
	if err != nil {
		t.Fatal(err)
	}
}

func saveDomainsInfoJson(hosts []string, prefix string) error {
	options := httpx.Options{
		Methods:         "GET",
		InputTargetHost: hosts,
		OnResult: func(r httpx.Result) {
			outputFilePath := prefix + "_" + r.Host + "_" + "domain_spec.json"
			outputFile, err := os.Create(outputFilePath)
			if err != nil {
				return
			}
			defer outputFile.Close()

			encoder := json.NewEncoder(outputFile)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(r); err != nil {
				return
			}
		},
	}

	if err := options.ValidateOptions(); err != nil {
		return err
	}

	httpxRunner, err := httpx.New(&options)
	if err != nil {
		return err
	}
	defer httpxRunner.Close()

	httpxRunner.RunEnumeration()
	return nil
}

func getSubdomainPrefix(path string) (map[string][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	subdomains := make(map[string][]string)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		domain := scanner.Text()
		parts := strings.Split(domain, ".")
		// A valid subdomain should have more than 2 parts (e.g., sub.domain.com)
		if len(parts) > 2 {
			subdomain := strings.Join(parts[:len(parts)-2], ".")
			rootDomain := strings.Join(parts[len(parts)-2:], ".")
			subdomains[rootDomain] = append(subdomains[rootDomain], subdomain)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return subdomains, nil
}

func createDomainsJson() error {
	var results []httpx.Result
	var domains []Domains

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), "mark") && filepath.Ext(info.Name()) == ".json" {
			fileData, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var result httpx.Result
			err = json.Unmarshal(fileData, &result)
			if err != nil {
				return err
			}

			results = append(results, result)
		}
		return nil
	})

	if err != nil {
		return err
	}

	for _, result := range results {
		domains = append(domains, Domains{
			URL:   result.URL,
			Label: buildLabel(result),
			Class: buildClass(result),
		})
	}

	domainsFilePath := "domains.json"
	domainsFile, err := os.Create(domainsFilePath)
	if err != nil {
		return err
	}
	defer domainsFile.Close()

	encoder := json.NewEncoder(domainsFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(domains); err != nil {
		return err
	}
	return nil
}

// todo (enka) labelのデータはbeautiful soupなどから取ってきた方がいい？
// labelのデータ生成方法が一番むずいかも
func buildLabel(r httpx.Result) string {
	return "unknown"
}

// todo (enka) 改善の余地あり
// ここのバリエーションをもっと増やす
func buildClass(r httpx.Result) string {
	var apiServerRef = 0
	var webApplicationRef = 0

	if r.CDNName == "aws" && r.CDNType == "cloud" {
		apiServerRef++
	}
	if r.WebServer == "envoy" {
		apiServerRef++
	}
	if r.WebServer == "awselb/2.0" {
		apiServerRef++
	}
	if r.WebServer == "Vercel" {
		webApplicationRef++
	}
	if r.ContentType == "text/plain" {
		webApplicationRef++
	}

	scores := map[string]int{
		"api server":      apiServerRef,
		"web application": webApplicationRef,
	}

	max := 0
	label := "unknown"

	for key, num := range scores {
		if num > max {
			max = num
			label = key
		}
	}

	return label
}
