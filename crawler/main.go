package main

import (
	"context"
	"fmt"
	"math"

	"github.com/projectdiscovery/alterx"
	"github.com/runetale/runevision/crawler/driver"
	"github.com/runetale/runevision/hack"
	"github.com/runetale/runevision/utility"
)

func main() {
	// main domainsを取得する
	domainChan := make(chan string)
	go driver.GetDomains(domainChan)

	for domain := range domainChan {
		runner, err := hack.NewSubFinder(&utility.Logger{})
		if err != nil {
			panic(err)
		}
		runner.SetParams(hack.SubfinderParams{
			TargetHost:         []string{domain},
			Threads:            1,
			MaxEnumerationTime: 10,
		})

		runner.SetParams(hack.SubfinderParams{
			TargetHost:         []string{"www.net"},
			Threads:            4,
			MaxEnumerationTime: 10,
		})

		subDomainChan := make(chan struct{})
		go runner.StartWithOutput(subDomainChan)
		<-subDomainChan

		opts := &alterx.Options{
			Domains: runner.GetTargetSubDomains(),
			MaxSize: math.MaxInt,
		}

		m, err := alterx.New(opts)
		if err != nil {
			panic(err)
		}
		results := m.Execute(context.Background())
		for result := range results {
			fmt.Println(result)
		}
	}

}
