package hack

import (
	"context"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/runetale/runevision/utility"
)

type NucleiRunner struct {
	runner      *nuclei.NucleiEngine
	targetHosts []string
	logger      *utility.Logger
}

type NucleiParams struct {
	TemplateTags []string
	TargetHosts  []string
}

func NewNuclei(logger *utility.Logger) (*NucleiRunner, error) {
	return &NucleiRunner{
		logger: logger,
	}, nil
}

func (s *NucleiRunner) Start() error {
	// load targets and optionally probe non http/https targets
	s.runner.LoadTargets(s.targetHosts, false)
	err := s.runner.ExecuteWithCallback(nil)
	defer s.runner.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *NucleiRunner) SetParams(params NucleiParams) error {
	ne, err := nuclei.NewNucleiEngineCtx(
		context.Background(),
		nuclei.WithTemplateFilters(nuclei.TemplateFilters{Tags: params.TemplateTags}),
		nuclei.EnableStatsWithOpts(nuclei.StatsOptions{MetricServerPort: 6064}), // optionally enable metrics server for better observability
	)
	if err != nil {
		return err
	}

	s.runner = ne
	s.targetHosts = params.TargetHosts
	return nil
}
