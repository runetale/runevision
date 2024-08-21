// DSLを使用したサブドメイン・ワードリスト・ジェネレータ
package hack

import (
	"bytes"
	"io"
	"math"
	"os"

	"github.com/projectdiscovery/alterx"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/runetale/runevision/utility"
)

type DslRunner struct {
	runner *alterx.Mutator
	logger *utility.Logger
}

func NewDslRunner(logger *utility.Logger) (*DslRunner, error) {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)
	return &DslRunner{
		logger: logger,
	}, nil
}

func (s *DslRunner) Start() error {
	var buffer bytes.Buffer
	writer := io.Writer(&buffer)
	s.runner.ExecuteWithWriter(writer)

	reader := io.Reader(&buffer)

	var output bytes.Buffer
	io.Copy(&output, reader)

	file, err := os.Create("dsl_output.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = output.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

func (s *DslRunner) SetParams(params ScanParams) error {
	opts := &alterx.Options{
		Domains: params.TargetHost,
		MaxSize: math.MaxInt,
	}

	m, err := alterx.New(opts)
	if err != nil {
		return err
	}

	s.runner = m
	return nil
}
