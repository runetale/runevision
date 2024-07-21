package repository_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/runetale/runevision/database"
	"github.com/runetale/runevision/di"
	"github.com/runetale/runevision/domain/config"
)

var (
	cfg config.Config
	db  *database.Postgres
)

func TestMain(m *testing.M) {
	cfg = config.MustLoad()
	db = di.InitializePostgres(cfg.Postgres, cfg.Log)
	err := db.MigrateUp("../../migrations")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Do tests
	//
	code := m.Run()

	// Clean
	//
	os.Exit(code)
}

func setup() {
	// 全テーブルのデータを削除
	if db != nil {
		err := db.Exec(`
				DO $$ DECLARE
					r RECORD;
				BEGIN
					FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
						EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
					END LOOP;
				END $$;
			`)
		if err != nil {
			log.Fatalf("Failed to truncate tables: %v", err)
		}
	}
}
