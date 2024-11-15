package api_router

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/runetale/thor/api_router/routes"
	"github.com/runetale/thor/domain/config"
)

type APIRouter struct {
	cfg    config.Config
	engine *echo.Echo
}

func NewAPIRouter(cfg config.Config) *APIRouter {
	return &APIRouter{
		cfg:    cfg,
		engine: echo.New(),
	}
}

func (r *APIRouter) Start() {
	r.setupEngine()
	r.setupEndpoints()
	fmt.Println(r.engine.Start(fmt.Sprintf("%s:%d", r.cfg.App.Host, r.cfg.App.Port)))
}

func (r *APIRouter) setupEngine() {
	r.engine.Use(middleware.Recover())
}

func (r *APIRouter) setupEndpoints() {
	apiGroup := r.engine.Group("/api")
	routes.RegisterDashboardRoutes(apiGroup, r.cfg)
	routes.RegisterHackRoutes(apiGroup, r.cfg)
}
