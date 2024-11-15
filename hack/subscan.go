// find sub domain
package hack

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"github.com/runetale/thor/utility"
)

type SubfinderRunner struct {
	runner                   *runner.Runner
	targetHost               []string
	logger                   *utility.Logger
	targetHostWithSubDomains []string
}

type SubfinderParams struct {
	TargetHost         []string
	Threads            int
	Timeout            int
	MaxEnumerationTime int
}

func NewSubFinder(logger *utility.Logger) (*SubfinderRunner, error) {
	return &SubfinderRunner{
		logger: logger,
	}, nil
}

func (s *SubfinderRunner) Start() error {
	// todo (snt) fix
	for _, host := range s.targetHost {
		output := &bytes.Buffer{}
		if err := s.runner.EnumerateSingleDomainWithCtx(context.Background(), host, []io.Writer{output}); err != nil {
			return err
		}
	}

	return nil
}

func (s *SubfinderRunner) StartWithOutput(out chan<- struct{}) error {
	defer close(out)
	for _, host := range s.targetHost {
		buf := &bytes.Buffer{}
		if err := s.runner.EnumerateSingleDomainWithCtx(context.Background(), host, []io.Writer{buf}); err != nil {
			return err
		}

		s.targetHostWithSubDomains = strings.Split(buf.String(), "\n")
	}

	out <- struct{}{}

	return nil
}

func (s *SubfinderRunner) GetTargetSubDomains() []string {
	return s.targetHostWithSubDomains
}

func (s *SubfinderRunner) SetParams(params SubfinderParams) error {
	s.targetHost = params.TargetHost
	subfinderOpts := &runner.Options{
		Threads:            params.Threads,            // Thread controls the number of threads to use for active enumerations
		Timeout:            params.Timeout,            // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: params.MaxEnumerationTime, // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
	}

	// disable timestamps in logs / configure logger
	log.SetFlags(0)

	subfinder, err := runner.NewRunner(subfinderOpts)
	if err != nil {
		return err
	}

	s.runner = subfinder

	return nil
}
