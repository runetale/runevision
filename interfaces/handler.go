package interfaces

import "github.com/labstack/echo/v4"

type DashboardHandler interface {
	Get(c echo.Context) error
	Add(c echo.Context) error
}

type HackHandler interface {
	Scan(c echo.Context) error
}
