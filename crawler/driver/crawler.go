package driver

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// ファイルからリストを読み込む関数
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	return lines, scanner.Err()
}

func GetDomains(output chan<- string) error {
	defer close(output)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	domains, err := readLines("./words/services-names.txt")
	if err != nil {
		fmt.Println("Error reading domains:", err)
		return nil
	}

	tlds, err := readLines("./words/tlds.txt")
	if err != nil {
		fmt.Println("Error reading TLDs:", err)
		return nil
	}

	// 結果をCSVファイルに保存する準備
	// file, err := os.OpenFile("./output/domain_status.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	fmt.Println("Error creating CSV file:", err)
	// 	return nil, err
	// }
	// defer file.Close()

	// writer := csv.NewWriter(file)
	// defer writer.Flush()

	// ヘッダーを書き込む
	// writer.Write([]string{"url", "status_code"})

	for _, domain := range domains {
		for _, tld := range tlds {
			// URLを組み立てる
			url := fmt.Sprintf("https://%s%s", domain, tld)

			// HTTPリクエストを送信
			resp, err := client.Get(url)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					output <- domain + tld
				}
			}

			// 結果をCSVに書き込む
			// writer.Write([]string{url, fmt.Sprintf("%d", statusCode)})
			// fmt.Printf("Checked %s: %d\n", url, statusCode)
		}
	}

	return nil
}
