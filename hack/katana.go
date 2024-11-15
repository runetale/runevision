// ウェブクローリングとスパイダーリング
// ウェブアプリケーションを探索し、エンドポイント、リンク、およびアセットを発見する

package hack

import (
	"math"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/katana/pkg/engine/standard"
	"github.com/projectdiscovery/katana/pkg/output"
	"github.com/projectdiscovery/katana/pkg/types"
	"github.com/runetale/thor/utility"
)

type KatanaRunner struct {
	crawlerOptions *types.CrawlerOptions
	targetHost     []string
	logger         *utility.Logger
}

type KatanaParams struct {
	TargetHost  []string
	MaxDepth    int
	FieldScope  string
	Timeout     int
	Concurrency int
	Parallelism int
	Delay       int
	RateLimit   int
	Strategy    string
}

// todo: パラメーターよくみる
// https://github.com/projectdiscovery/katana
func NewKatana(logger *utility.Logger) (*KatanaRunner, error) {
	gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose) // increase the verbosity (optional)
	return &KatanaRunner{
		logger: logger,
	}, nil
}

func (s *KatanaRunner) Start() error {
	defer s.crawlerOptions.Close()
	crawler, err := standard.New(s.crawlerOptions)
	if err != nil {
		return err
	}

	defer crawler.Close()
	for _, host := range s.targetHost {
		err = crawler.Crawl(host)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *KatanaRunner) SetParams(params KatanaParams) error {
	options := &types.Options{
		MaxDepth:     params.MaxDepth,
		FieldScope:   params.FieldScope,
		BodyReadSize: math.MaxInt,
		Timeout:      params.Timeout,
		Concurrency:  params.Concurrency,
		Parallelism:  params.Parallelism,
		Delay:        params.Delay,
		RateLimit:    params.RateLimit,
		Strategy:     params.Strategy,
		OnResult: func(result output.Result) {
			gologger.Info().Msg(result.Request.URL)
		},
	}

	crawlerOptions, err := types.NewCrawlerOptions(options)
	if err != nil {
		return err
	}
	s.crawlerOptions = crawlerOptions
	s.targetHost = params.TargetHost

	return nil
}
