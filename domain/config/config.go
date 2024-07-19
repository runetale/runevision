package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App      App      `required:"true" envconfig:"APP"`
	Postgres Postgres `required:"true" envconfig:"POSTGRES"`
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

func MustLoad() Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		panic(err)
	}
	return c
}
