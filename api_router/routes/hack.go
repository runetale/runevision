package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/runetale/thor/di"
	"github.com/runetale/thor/domain/config"
)

func RegisterHackRoutes(g *echo.Group, cfg config.Config) {
	h := di.InitializeHackHandler(cfg.Postgres, cfg.Log)
	g.POST("/hack", h.Scan)
}
