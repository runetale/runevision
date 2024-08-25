package driver

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// ファイルからリストを読み込む関数
func readLines(filename string, limit int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
		if limit > 0 && len(lines) >= limit {
			break
		}
	}

	return lines, scanner.Err()
}

func Crawl() error {
	// コマンドライン引数をパース
	n := flag.Int("n", 0, "Limit the number of items to process from each file")
	flag.Parse()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 各ファイルからリストを読み込む
	subdomains, err := readLines("./words/deepmagic.com-prefixes-top500.txt", *n)
	if err != nil {
		fmt.Println("Error reading subdomains:", err)
		return err
	}

	domains, err := readLines("./words/services-names.txt", *n)
	if err != nil {
		fmt.Println("Error reading domains:", err)
		return err
	}

	tlds, err := readLines("./words/tlds.txt", *n)
	if err != nil {
		fmt.Println("Error reading TLDs:", err)
		return err
	}

	// 結果をCSVファイルに保存する準備
	file, err := os.Create("./output/domain_status.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーを書き込む
	writer.Write([]string{"url", "status_code"})

	for _, subdomain := range subdomains {
		for _, domain := range domains {
			for _, tld := range tlds {
				// URLを組み立てる
				url := fmt.Sprintf("https://%s.%s%s", subdomain, domain, tld)

				// HTTPリクエストを送信
				resp, err := client.Get(url)
				statusCode := 0
				if err == nil {
					statusCode = resp.StatusCode
					resp.Body.Close()
				}

				// 結果をCSVに書き込む
				writer.Write([]string{url, fmt.Sprintf("%d", statusCode)})
				fmt.Printf("Checked %s: %d\n", url, statusCode)
			}
		}
	}
	return nil
}
