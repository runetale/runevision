// this package port scanning by naabu
package hack

import (
	"context"
	"fmt"

	"github.com/projectdiscovery/naabu/v2/pkg/result"
	"github.com/projectdiscovery/naabu/v2/pkg/runner"
	"github.com/runetale/thor/utility"
)

type ScanType string

const (
	SynScan     ScanType = "s"
	ConnectScan ScanType = "c"
)

type ScanParams struct {
	TargetHost  []string
	ScanType    ScanType
	TargetPorts string
}

type PortScanner struct {
	runner       *runner.Runner
	logger       *utility.Logger
	isPrivileged bool
}

func NewPortScanner(logger *utility.Logger) (*PortScanner, error) {
	return &PortScanner{
		logger: logger,
	}, nil
}

func (s *PortScanner) Start(isPrivileged bool) error {
	s.isPrivileged = isPrivileged
	defer s.runner.Close()

	ctx := context.TODO()
	err := s.runner.RunEnumeration(ctx)
	if err != nil {
		return err
	}
	return nil
}

// todo (snt) fix
func (s *PortScanner) SetParams(params ScanParams) error {
	options := runner.Options{
		Host:     params.TargetHost,
		ScanType: string(params.ScanType),
		JSON:     true,
		OnResult: func(hr *result.HostResult) {
			fmt.Println(hr.Host, hr.Ports)
		},
		ScanAllIPS: true,
		Ports:      params.TargetPorts,
	}

	naabuRunner, err := runner.NewRunner(&options)
	if err != nil {
		return err
	}

	s.runner = naabuRunner
	return nil
}
