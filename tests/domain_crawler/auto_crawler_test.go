package domaincrawler_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/projectdiscovery/alterx"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

type Domains struct {
	URL   string `json:"url"`
	Label string `json:"label"`
	Class string `json:"class"`
}

func TestDomainAutoCrawler(t *testing.T) {
	dslfilePath := "dsl_output.txt"
	subDomainsfilePath := "sub_domains.txt"
	targetHost := "runetale.com"

	subfinderOpts := &runner.Options{
		Threads:            10, // Thread controls the number of threads to use for active enumerations
		Timeout:            30, // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
	}

	// disable timestamps in logs / configure logger
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

	dnsfile, err := os.Open(dslfilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer dnsfile.Close()

	var domains []Domains
	scanner = bufio.NewScanner(dnsfile)
	for scanner.Scan() {
		line := scanner.Text()
		// todo: (enka) curlなどでrequest投げて、ドメインの生存が確認できたドメインだけにする。自動化

		// todo: (enka) 色々なprefixに対応
		// "api"を含むドメインをカウント
		if strings.Contains(line, "api") {
			domain := Domains{
				URL:   "https://" + line,
				Label: "VPN Server", // httpxなどでそれっぽい値を取得して自動化？
				Class: "APIサーバー",
			}
			domains = append(domains, domain)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	outputFilePath := "domains.json"

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(domains); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("output json files => %s\n", outputFilePath)
}
