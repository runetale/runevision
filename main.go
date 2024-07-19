package main

import (
	"os"

	"github.com/runetale/runevision/api_router"
	"github.com/runetale/runevision/database"
	"github.com/runetale/runevision/domain/config"
	"github.com/runetale/runevision/utility"
)

func main() {
	// api
	cfg := config.MustLoad()
	log, err := utility.NewLogger(os.Stdout, utility.JsonFmtStr, utility.DebugLevelStr)
	if err != nil {
		panic(err)
	}
	db, err := database.NewPostgresFromConfig(log, cfg.Postgres)
	if err != nil {
		panic(err)
	}
	err = db.MigrateUp("migrations")
	if err != nil {
		panic(err)
	}

	r := api_router.NewAPIRouter(cfg)
	r.Start()
}
