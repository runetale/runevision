package interactor

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/runetale/thor/database"
	"github.com/runetale/thor/domain/entity"
	"github.com/runetale/thor/domain/requests"
	"github.com/runetale/thor/interfaces"
	"github.com/runetale/thor/localclient"
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

func (i *HackInteractor) Scan(request *requests.HackDoScanRequest, c echo.Context) (*entity.ScanResponse, error) {
	// todo (snt)
	// redisでidをcacheする
	sequentialID := uuid.New().String()
	return i.lc.DoScan(sequentialID, request, c)
}
