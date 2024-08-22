package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

type Visionary struct {
	Name           string         `json:"name"`
	Tags           string         `json:"tags"`
	Severity       string         `json:"severity"`
	Description    string         `json:"description"`
	Impact         string         `json:"impact"`
	Remediation    string         `json:"remediation"`
	Classification Classification `json:"classification"`
	Metadata       Metadata       `json:"metadata"`
	Vendor         string         `json:"vendor"`
	Product        string         `json:"product"`
	CvssMetrics    string         `json:"cvss-metrics"`
	CvssScore      float64        `json:"cvss-score"`
	CveID          string         `json:"id"`
	CweID          string         `json:"cwe-id"`
	EpssScore      float64        `json:"epss-score"`
	EpssPercentile float64        `json:"epss-percentile"`
	Cpe            string         `json:"cpe"`
}

type AttackType struct {
	Type string `json:"type"`
}

func main() {
	var visionary []*Visionary
	var attackType []AttackType
	var tags []string

	// walking dir for yaml
	dir := "./target_yaml"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			visionary, tags = parseYAMLToVisionary(path, visionary, tags)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	saveAttackTypeJson(attackType, tags)

	saveVisonaryJson(visionary)
}

func saveAttackTypeJson(at []AttackType, tags []string) {
	tagsFile := "tags.txt"
	formattedName := strings.Join(tags, "\n")
	file, err := os.Create("tags.txt")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(formattedName)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	// 重複したtagを消して、attack_type.jsonを出力する
	file, err = os.Open(tagsFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content string
	for scanner.Scan() {
		content += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	items := strings.Split(content, ",")

	seen := make(map[string]bool)
	var uniqueItems []string

	for _, item := range items {
		trimmedItem := strings.TrimSpace(item)
		if _, exists := seen[trimmedItem]; !exists {
			seen[trimmedItem] = true
			uniqueItems = append(uniqueItems, trimmedItem)
		}
	}

	for _, o := range uniqueItems {
		at = append(at, AttackType{
			Type: o,
		})
	}

	filePath := "attack_types.json"
	output, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(at); err != nil {
		panic(err)
	}

	fmt.Println("Successfully saved to attack_types.json")
}

func saveVisonaryJson(v []*Visionary) {
	filePath := "visonary.json"
	output, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		panic(err)
	}

	fmt.Println("Successfully saved to visonary.json")
}

func parseYAMLToVisionary(file string, v []*Visionary, tags []string) ([]*Visionary, []string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file %s: %v", file, err)
	}

	var cve CVE
	err = yaml.Unmarshal(data, &cve)
	if err != nil {
		log.Fatalf("Failed to parse YAML file %s: %v", file, err)
	}

	fmt.Printf("File: %s, ID: %s, Name: %s\n", file, cve.ID, cve.Info.Name)
	v = append(v, &Visionary{
		Name:           cve.Info.Name,
		Tags:           cve.Info.Tags,
		Severity:       cve.Info.Severity,
		Description:    cve.Info.Description,
		Impact:         cve.Info.Impact,
		Remediation:    cve.Info.Remediation,
		Classification: cve.Info.Classification,
		Metadata:       cve.Info.Metadata,
		CveID:          cve.ID,
		Vendor:         cve.Info.Metadata.Vendor,
		Product:        cve.Info.Metadata.Product,
		CvssMetrics:    cve.Info.Classification.CvssMetrics,
		CvssScore:      cve.Info.Classification.CvssScore,
		CweID:          cve.Info.Classification.CweID,
		EpssScore:      cve.Info.Classification.EpssScore,
		EpssPercentile: cve.Info.Classification.EpssPercentile,
		Cpe:            cve.Info.Classification.Cpe,
	})

	tags = append(tags, cve.Info.Tags)

	return v, tags
}
