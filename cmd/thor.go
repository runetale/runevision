package main

import (
	"fmt"
	"os"

	"github.com/runetale/thor/cmd/thor"
)

func main() {
	if err := thor.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
