// cve database servers
package hack

import (
	"github.com/projectdiscovery/cvemap/pkg/runner"
	"github.com/projectdiscovery/gologger"
)

func NewCveMapServer() {
	options := runner.ParseOptions()
	runner, err := runner.New(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}
	runner.Run()
}
