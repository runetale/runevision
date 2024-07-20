package main

import (
	"github.com/runetale/runevision/api_router"
	"github.com/runetale/runevision/di"
	"github.com/runetale/runevision/domain/config"
)

func main() {
	// api
	cfg := config.MustLoad()
	db := di.InitializePostgres(cfg.Postgres, cfg.Log)
	err := db.MigrateUp("migrations")
	if err != nil {
		panic(err)
	}

	r := api_router.NewAPIRouter(cfg)
	r.Start()
}
