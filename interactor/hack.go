package interactor

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/runetale/runevision/database"
	"github.com/runetale/runevision/domain/entity"
	"github.com/runetale/runevision/domain/requests"
	"github.com/runetale/runevision/interfaces"
	"github.com/runetale/runevision/localclient"
)

type HackInteractor struct {
	db *database.Postgres
	lc *localclient.LocalClient
}

func NewHackInteractor(
	db *database.Postgres,
	lc *localclient.LocalClient,
) interfaces.HackInteractor {
	return &HackInteractor{
		db: db,
		lc: lc,
	}
}

func (i *HackInteractor) Scan(request *requests.HackDoScanRequest, c echo.Context) (*entity.HackHistory, error) {
	// todo (snt)
	// redisでidをcacheする
	sequentialID := uuid.New().String()
	return i.lc.DoScan(sequentialID, request, c)
}
