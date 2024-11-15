// インターネット上で公開されているホストを検索
// 現状はcensysを使って検索している
package hack

import (
	"context"
	"fmt"

	"github.com/projectdiscovery/uncover"
	"github.com/projectdiscovery/uncover/sources"
	"github.com/runetale/thor/utility"
)

type HostFinderParams struct {
	Queries  []string
	Limit    int
	MaxRetry int
	Timeout  int
}

type HostFinder struct {
	runner *uncover.Service
	logger *utility.Logger
}

func NewHostFinder(logger *utility.Logger) (*HostFinder, error) {
	opts := uncover.Options{
		Agents: []string{"censys"},
	}

	u, err := uncover.New(&opts)
	if err != nil {
		return nil, err
	}

	return &HostFinder{
		runner: u,
		logger: logger,
	}, nil
}

func (s *HostFinder) Start() error {
	result := func(result sources.Result) {
		fmt.Println("HostFinder: Get IPPort ------")
		fmt.Println(result.IpPort())
		fmt.Println("HostFinder: Get Source ------")
		fmt.Println(result.Source)
		fmt.Println("HostFinder: Get JSON ------")
		fmt.Println(result.JSON())
	}

	if err := s.runner.ExecuteWithCallback(context.TODO(), result); err != nil {
		return err
	}

	return nil
}

func (s *HostFinder) SetParams(params HostFinderParams) {
	s.runner.Options.Queries = params.Queries
	s.runner.Options.Limit = params.Limit
	s.runner.Options.MaxRetry = params.MaxRetry
	s.runner.Options.Timeout = params.Timeout
}
