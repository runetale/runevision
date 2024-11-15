// find url paths

package hack

import (
	"fmt"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/httpx/runner"
	"github.com/runetale/thor/utility"
)

type HttpxRunner struct {
	runner *runner.Runner
	logger *utility.Logger
}

// todo: 後でパラメーターよくみる
// https://github.com/projectdiscovery/httpx/tree/main
type HttpxParams struct {
	Methods     string
	TargetHosts []string
}

func NewHttpx(logger *utility.Logger) (*HttpxRunner, error) {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose) // increase the verbosity (optional)

	return &HttpxRunner{
		logger: logger,
	}, nil
}

func (s *HttpxRunner) Start() error {
	defer s.runner.Close()
	s.runner.RunEnumeration()
	return nil
}

func (s *HttpxRunner) SetParams(params HttpxParams) error {
	options := runner.Options{
		Methods:         params.Methods,
		InputTargetHost: params.TargetHosts,
		OnResult: func(r runner.Result) {
			// handle error
			if r.Err != nil {
				fmt.Printf("[Err] %s: %s\n", r.Input, r.Err)
				return
			}
			fmt.Println(r)
		},
	}

	if err := options.ValidateOptions(); err != nil {
		return err
	}

	httpxRunner, err := runner.New(&options)
	if err != nil {
		return err
	}

	s.runner = httpxRunner
	return nil
}
