package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App      App      `required:"true" envconfig:"APP"`
	Postgres Postgres `required:"true" envconfig:"POSTGRES"`
	Log      Log      `required:"true" envconfig:"LOG"`
}

type App struct {
	Host string `required:"true" envconfig:"HOST"`
	Port uint   `required:"true" envconfig:"PORT"`
}

type Postgres struct {
	UserName     string `required:"true" envconfig:"USER_NAME"`
	Password     string `required:"true" envconfig:"PASSWORD"`
	Host         string `required:"true" envconfig:"HOST"`
	Port         uint   `required:"true" envconfig:"PORT"`
	DatabaseName string `required:"true" envconfig:"DATABASE_NAME"`
}

type Log struct {
	Format string `required:"true" envconfig:"FORMAT"`
	Level  string `required:"true" envconfig:"LEVEL"`
}

func MustLoad() Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		panic(err)
	}
	return c
}
