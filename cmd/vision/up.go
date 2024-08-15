package vision

import (
	"context"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/ffcli"
	"github.com/runetale/runevision/backend"
	"github.com/runetale/runevision/safesocket"
	"github.com/runetale/runevision/utility"
)

// // local backend server
// logger, err := utility.NewLogger(os.Stdout, "json", "debug")
// if err != nil {
// 	fmt.Printf("failed to initialze logger: %v", err)
// 	return
// }

// ln, err := safesocket.Listen(safesocket.VisionSocketPath())
// if err != nil {
// 	fmt.Printf("failed to listen safe socket: %v", err)
// 	return
// }

// bs := backend.New(logger)
// err = bs.Run(context.Background(), ln)
// if err != nil {
// 	fmt.Printf("failed to start backend server: %v", err)
// 	return
// }

var upCmd = &ffcli.Command{
	Name: "up",
	Exec: execUp,
}

func execUp(args []string) error {
	// // local backend server
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
	err = bs.Run(context.Background(), ln)
	if err != nil {
		fmt.Printf("failed to start backend server: %v", err)
		return err
	}

	return nil
}
