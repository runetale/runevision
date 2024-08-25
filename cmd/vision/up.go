package vision

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/peterbourgon/ff/ffcli"
	"github.com/runetale/runevision/backend"
	"github.com/runetale/runevision/safesocket"
	"github.com/runetale/runevision/utility"
	"github.com/runetale/runevision/vsengine"
)

var upCmd = &ffcli.Command{
	Name: "up",
	Exec: execUp,
}

func execUp(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init local backend server
	logger, err := utility.NewLogger(os.Stdout, "json", "debug")
	if err != nil {
		fmt.Printf("failed to initialze logger: %v", err)
		return err
	}

	ln, err := safesocket.Listen(safesocket.VisionSocketPath())
	if err != nil {
		fmt.Printf("failed to listen safe socket: %v", err)
		return err
	}

	bs := backend.New(logger)
	go bs.Run(ctx, ln)

	// init vision engine
	engine, err := vsengine.NewEngine(false)
	if err != nil {
		fmt.Printf("failed to listen safe socket: %v", err)
		return err
	}
	lb := backend.NewLocalBackend(engine, logger)
	bs.SetLocalBackend(lb)

	logger.Logger.Debug("started vision daemon")

	ch := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGINT,
		)
		select {
		case <-c:
			close(ch)
		}
	}()
	<-ch

	return nil
}
