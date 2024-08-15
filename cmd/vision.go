package main

import (
	"fmt"
	"os"

	"github.com/runetale/runevision/cmd/vision"
)

func main() {
	if err := vision.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
