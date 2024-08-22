package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Info struct {
	Name           string         `yaml:"name"`
	Author         string         `yaml:"author"`
	Severity       string         `yaml:"severity"`
	Description    string         `yaml:"description"`
	Impact         string         `yaml:"impact"`
	Remediation    string         `yaml:"remediation"`
	Classification Classification `yaml:"classification"`
	Metadata       Metadata       `yaml:"metadata"`
	Tags           string         `yaml:"tags"`
}

type Classification struct {
	CvssMetrics    string  `yaml:"cvss-metrics"`
	CvssScore      float64 `yaml:"cvss-score"`
	CveID          string  `yaml:"cve-id"`
	CweID          string  `yaml:"cwe-id"`
	EpssScore      float64 `yaml:"epss-score"`
	EpssPercentile float64 `yaml:"epss-percentile"`
	Cpe            string  `yaml:"cpe"`
}

type Metadata struct {
	MaxRequest int    `yaml:"max-request"`
	Vendor     string `yaml:"vendor"`
	Product    string `yaml:"product"`
}

type HTTP struct {
	Method            string    `yaml:"method"`
	Path              []string  `yaml:"path"`
	MatchersCondition string    `yaml:"matchers-condition"`
	Matchers          []Matcher `yaml:"matchers"`
}

type Matcher struct {
	Type   string   `yaml:"type"`
	Regex  []string `yaml:"regex"`
	Status []int    `yaml:"status"`
}

type CVE struct {
	ID   string `yaml:"id"`
	Info Info   `yaml:"info"`
	HTTP []HTTP `yaml:"http"`
}

type Result struct {
	Name string `json:"name"`
}

func main() {
	var result []string

	dir := "./target_yaml"

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			result = parseYAML(path, result)
		}
		return nil
	})

	formattedName := strings.Join(result, "\n")
	file, err := os.Create("output.txt")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// 改行して保存
	_, err = file.WriteString(formattedName)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("Successfully saved to output.txt")
}

func parseYAML(file string, result []string) []string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", file, err)
	}

	var cve CVE
	err = yaml.Unmarshal(data, &cve)
	if err != nil {
		log.Fatalf("Failed to parse YAML file %s: %v", file, err)
	}

	// id と info.name を取得
	fmt.Printf("File: %s, ID: %s, Name: %s\n", file, cve.ID, cve.Info.Name)
	result = append(result, toSnakeCase(cve.Info.Name))

	return result
}

func toSnakeCase(str string) string {
	re := regexp.MustCompile(`[^\w\s]`)
	str = re.ReplaceAllString(str, "")
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ReplaceAll(str, "__", "_")
	return strings.ToUpper(str)
}
