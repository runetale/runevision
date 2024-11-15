// このファイルはvisionary.jsonとattack_types.jsonを吐き出すツール
// visionary.jsonはcveなどの脆弱性のデータを持った巨大なjson
// attack_types.jsonは膨大な脆弱性のデータからどのようなテクノロジーを使用し、どのようかサービスカテゴリかを記録したjson
// 後にこれらはwhite hat moduleのmlとして学習されていくデータ群kA
// *サービスカテゴリーにおいてどのような攻撃手法が的確を判断する

// 使われ方としては、ターゲットのドメインがどのようなサービスの性質かをmlでfilterして
// そのターゲットドメインがこのような形でリクエストが送られてくるので、そのattack_typeに応じた攻撃手法を返す
// “
//
//	{
//	    "url": "https://caterpie.runetale.com",
//	    "label": "unknown",
//	   	"class": "api server",
//		"attack_type": "ai"
//	}
//
// ```
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/runetale/thor/tools/visonary/types"
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

type Attack struct {
	TechStack      types.TechType       `json:"tech_stack"`
	Category       types.AttackCategory `json:"categry"`
	CveID          string               `json:"cve_id"`
	RunetaleImpact bool                 `json:"rt_impact"`
}

func main() {
	var visionary []*Visionary

	// walking dir for yaml
	dir := "./target_yaml"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			visionary = parseYAMLToVisionary(path, visionary)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	saveAttackTypeJson(visionary)
	saveVisonaryJson(visionary)
}

func saveAttackTypeJson(visionary []*Visionary) {
	var attackType []*Attack
	for _, v := range visionary {
		items := strings.Split(v.Tags, ",")

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
			attackType = append(attackType, &Attack{
				TechStack:      types.TechType(o),
				Category:       types.GetAttackCategory(types.TechType(o)),
				CveID:          v.CveID,
				RunetaleImpact: types.GetRunetaleImpact(types.GetAttackCategory(types.TechType(o))),
			})
		}
	}

	filePath := "attack_types.json"
	output, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(attackType); err != nil {
		panic(err)
	}

	fmt.Println("Successfully saved to attack_types.json")
	return
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

func parseYAMLToVisionary(file string, v []*Visionary) []*Visionary {
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

	return v
}
