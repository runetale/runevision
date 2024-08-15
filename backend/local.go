// local backendはlocal backend serverからapi経由で
// vision engineの実行や、cvemap binaryの実行などを行います
package backend

import (
	"sync"

	"github.com/runetale/runevision/utility"
	"github.com/runetale/runevision/vsengine"
)

type LocalBackend struct {
	engine *vsengine.Engine
	mu     sync.Mutex
	logger *utility.Logger
}

func (b *LocalBackend) Shutdown() {

}
