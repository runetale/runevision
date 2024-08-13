package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/runetale/runevision/api_router"
	"github.com/runetale/runevision/di"
	"github.com/runetale/runevision/domain/config"
)

func main() {
	cfg := config.MustLoad()

	// api
	//
	db := di.InitializePostgres(cfg.Postgres, cfg.Log)
	err := db.MigrateUp("migrations")
	if err != nil {
		panic(err)
	}

	r := api_router.NewAPIRouter(cfg)
	go r.Start()

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
}
