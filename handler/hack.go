package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/runetale/runevision/domain/requests"
	"github.com/runetale/runevision/interfaces"
)

type HackHandler struct {
	hackInteractor interfaces.HackInteractor
}

func NewHackHandler(
	hackInteractor interfaces.HackInteractor,
) interfaces.HackHandler {
	return &HackHandler{
		hackInteractor: hackInteractor,
	}
}

func (h *HackHandler) Scan(c echo.Context) error {
	var req requests.HackDoScanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	res, err := h.hackInteractor.Scan(&req)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, res)
	return nil
}
