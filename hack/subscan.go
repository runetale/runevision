// find sub domain
package hack

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"github.com/runetale/runevision/utility"
)

type SubfinderRunner struct {
	runner     *runner.Runner
	targetHost []string
	logger     *utility.Logger
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
		log.Println(output.String())
	}

	return nil
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
