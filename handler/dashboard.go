package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/runetale/thor/domain/entity"
	"github.com/runetale/thor/domain/requests"
	"github.com/runetale/thor/domain/responses"
	"github.com/runetale/thor/interfaces"
)

type DashboardHandler struct {
	dashboardInteractor interfaces.DashboardInteractor
}

func NewDashboardHandler(
	dashboardInteractor interfaces.DashboardInteractor,
) interfaces.DashboardHandler {
	return &DashboardHandler{
		dashboardInteractor: dashboardInteractor,
	}
}

func (h *DashboardHandler) Get(c echo.Context) error {
	histories, err := h.dashboardInteractor.Get()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return err
	}
	c.JSON(http.StatusOK, responses.NewGetDashboardResponse(1, 10, histories))
	return nil
}

func (h *DashboardHandler) Add(c echo.Context) error {
	var req requests.AddDashboardRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	err := h.dashboardInteractor.Add(entity.NewDashboardHistoryFromRequest(req))
	c.JSON(http.StatusOK, responses.NewOkResponse(err == nil))
	return err
}
