package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/runetale/runevision/di"
	"github.com/runetale/runevision/domain/config"
)

func RegisterHackRoutes(g *echo.Group, cfg config.Config) {
	h := di.InitializeHackHandler(cfg.Postgres, cfg.Log)
	g.POST("/hack", h.Scan)
}
