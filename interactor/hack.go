package interactor

import (
	"github.com/google/uuid"
	"github.com/runetale/runevision/database"
	"github.com/runetale/runevision/domain/entity"
	"github.com/runetale/runevision/domain/requests"
	"github.com/runetale/runevision/interfaces"
	"github.com/runetale/runevision/types"
	"github.com/runetale/runevision/vsengine"
)

type HackInteractor struct {
	db *database.Postgres
}

func NewHackInteractor(
	db *database.Postgres,
) interfaces.HackInteractor {
	return &HackInteractor{
		db: db,
	}
}

func (i *HackInteractor) Scan(request *requests.HackDoScanRequest) (*entity.HackHistory, error) {
	// todo (snt) refactor
	// API(local backend)経由で実行するようにする
	isDebug := true
	engine, err := vsengine.NewEngine(isDebug)
	if err != nil {
		return nil, err
	}

	// todo (snt)
	// redisでidをcacheする、終わればメモリを解放
	sequentialID := uuid.New().String()
	err = engine.Reconfig(types.SequenceID(sequentialID), request)
	if err != nil {
		return nil, err
	}

	status := engine.GetStatus(types.SequenceID(sequentialID))
	hh := entity.NewHackHistory(request.Name, sequentialID, string(status))

	return hh, nil
}
