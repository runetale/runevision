// local backendはlocal backend serverからapi経由で
// vision engineの実行や、cvemap binaryの実行などを行います
package backend

import (
	"sync"

	"github.com/runetale/runevision/domain/requests"
	"github.com/runetale/runevision/types"
	"github.com/runetale/runevision/utility"
	"github.com/runetale/runevision/vsengine"
)

type LocalBackend struct {
	engine vsengine.Engine
	mu     sync.Mutex
	logger *utility.Logger
}

func NewLocalBackend(engine vsengine.Engine, logger *utility.Logger) *LocalBackend {
	return &LocalBackend{
		engine: engine,
		logger: logger,
	}
}

func (b *LocalBackend) Shutdown() {

}
func (vb *LocalBackend) Scan(sequentialID types.SequenceID, request *requests.HackDoScanRequest) error {
	return vb.engine.Reconfig(sequentialID, request)
}

func (vb *LocalBackend) GetStatus(sequentialID types.SequenceID) types.AttackStatus {
	return vb.engine.GetStatus(sequentialID)
}

func (vb *LocalBackend) Ping() (string, error) {
	return "ping", nil
}
