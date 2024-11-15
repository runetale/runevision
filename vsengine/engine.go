package vsengine

import (
	"github.com/runetale/thor/domain/requests"
	"github.com/runetale/thor/types"
)

type Engine interface {
	// reconfig attacker paramters
	Reconfig(sequentialID types.SequenceID, request *requests.HackDoScanRequest) error
	// closing vengine
	Close() error
	// get attacker status by sequence ID
	GetStatus(sid types.SequenceID) types.AttackStatus
}
