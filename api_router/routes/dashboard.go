package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/runetale/thor/di"
	"github.com/runetale/thor/domain/config"
)

func RegisterDashboardRoutes(g *echo.Group, cfg config.Config) {
	h := di.InitializeDashboardHandler(cfg.Postgres, cfg.Log)
	g.GET("/dashboard", h.Get)
	g.POST("/dashboard", h.Add)
}
