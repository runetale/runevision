package vision

import (
	"flag"
	"strings"

	"github.com/peterbourgon/ff/ffcli"
)

func Run(args []string) error {
	if len(args) == 1 && (args[0] == "-V" || args[0] == "--version" || args[0] == "-v") {
		args = []string{"version"}
	}

	fs := flag.NewFlagSet("vision", flag.ExitOnError)

	cmd := &ffcli.Command{
		Name:      "",
		Usage:     "",
		ShortHelp: "",
		LongHelp: strings.TrimSpace(`
`),
		Subcommands: []*ffcli.Command{
			upCmd,
		},
		FlagSet: fs,
		Exec:    func(args []string) error { return flag.ErrHelp },
	}

	if err := cmd.Run(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	return nil
}
