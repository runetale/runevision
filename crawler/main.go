package main

import (
	"github.com/runetale/runevision/crawler/driver"
)

func main() {
	err := driver.Crawl()
	if err != nil {
		panic(err)
	}
}
