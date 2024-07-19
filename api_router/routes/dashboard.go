package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/runetale/runevision/di"
	"github.com/runetale/runevision/interfaces"
)

func RegisterDashboardRoutes(g *echo.Group, db interfaces.SQLExecuter) {
	h := di.InitializeDashboardHandler(db)
	g.GET("/dashboard", h.Get)
	g.POST("/dashboard", h.Add)
}
