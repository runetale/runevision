package api_router

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/runetale/runevision/api_router/routes"
	"github.com/runetale/runevision/database"
	"github.com/runetale/runevision/domain/config"
	"github.com/runetale/runevision/interfaces"
	"github.com/runetale/runevision/utility"
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
	db, err := database.NewPostgresFromConfig(&utility.Logger{}, r.cfg.Postgres)
	if err != nil {
		panic(err)
	}
	r.setupEngine()
	r.setupEndpoints(db)
	fmt.Println(r.engine.Start(fmt.Sprintf("%s:%d", r.cfg.App.Host, r.cfg.App.Port)))
}

func (r *APIRouter) setupEngine() {
	r.engine.Use(middleware.Recover())
}

func (r *APIRouter) setupEndpoints(db interfaces.SQLExecuter) {
	apiGroup := r.engine.Group("/api")
	routes.RegisterDashboardRoutes(apiGroup, db)
}
