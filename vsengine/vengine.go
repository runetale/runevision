// thor engine
// attack surface management engine
// ここからattackerなどを操作する
package vsengine

import (
	"io"
	"os"
	"sync"

	"github.com/runetale/thor/domain/requests"
	"github.com/runetale/thor/hack"
	"github.com/runetale/thor/types"
	"github.com/runetale/thor/utility"
	"github.com/runetale/thor/utility/privileges"
	"github.com/runetale/thor/vsd"
)

// custom error hooks
type closeOnErrorPool []func()

func (p *closeOnErrorPool) add(c io.Closer) { *p = append(*p, func() { c.Close() }) }
func (p closeOnErrorPool) closeAllIfError(errp *error) {
	if *errp != nil {
		for _, closeFn := range p {
			closeFn()
		}
	}
}

type visionEngine struct {
	isDebug bool

	attacker *Attacker

	vsd *vsd.VisionSystem

	attackerStatus map[types.SequenceID]types.AttackStatus

	mu sync.Mutex

	waitCh chan struct{}
}

func newVisionEngine(isDebug bool) *visionEngine {
	return &visionEngine{
		attackerStatus: make(map[types.SequenceID]types.AttackStatus),
		isDebug:        isDebug,
	}
}

func NewEngine(isDebug bool) (_ Engine, reterr error) {
	var closePool closeOnErrorPool
	defer closePool.closeAllIfError(&reterr)

	e := newVisionEngine(isDebug)
	e.attacker = NewAttacker()
	closePool.add(e)
	closePool.add(e.attacker)

	// set vsd system
	vsys, err := initVisionSystem()
	if err != nil {
		return
	}
	e.vsd = vsys

	// set callback function
	updateAttackStatusFn := func(sid types.SequenceID, status types.AttackStatus) {
		e.mu.Lock()
		e.attackerStatus[sid] = status
		e.mu.Unlock()
	}
	e.attacker.SetUpdateAttackStatusFn(updateAttackStatusFn)

	// starting attacker routine
	go e.attacker.Start()

	go func() {
		select {
		case <-e.waitCh:
			// continue
		}
	}()

	return e, nil
}

func initVisionSystem() (sys *vsd.VisionSystem, err error) {
	logger, err := utility.NewLogger(os.Stdout, "json", "debug")
	if err != nil {
		return
	}

	// init vsd
	sys = new(vsd.VisionSystem)
	sys.IsPrivileged = privileges.IsPrivileged

	// init host finder
	hostRun, err := hack.NewHostFinder(logger)
	if err != nil {
		panic(err)
	}
	sys.Set(hostRun)

	// init dsl
	dslRun, err := hack.NewDslRunner(logger)
	if err != nil {
		panic(err)
	}
	sys.Set(dslRun)

	// init port scan
	pscan, err := hack.NewPortScanner(logger)
	if err != nil {
		panic(err)
	}
	sys.Set(pscan)

	// init subfinder
	sfRun, err := hack.NewSubFinder(logger)
	if err != nil {
		panic(err)
	}
	sys.Set(sfRun)

	// init httpx
	hxRun, err := hack.NewHttpx(logger)
	if err != nil {
		panic(err)
	}
	sys.Set(hxRun)

	// init katana
	katanaRun, err := hack.NewKatana(logger)
	if err != nil {
		panic(err)
	}
	sys.Set(katanaRun)

	// init nuclei
	nuRun, err := hack.NewNuclei(logger)
	if err != nil {
		panic(err)
	}
	sys.Set(nuRun)

	return
}

func (e *visionEngine) Reconfig(sequentialID types.SequenceID, request *requests.HackDoScanRequest) error {
	err := e.attacker.Reconfig(sequentialID, e.vsd, request)
	if err != nil {
		return err
	}

	return nil
}

func (e *visionEngine) Close() error {
	close(e.waitCh)
	return nil
}

// from attacker 'updateAttackStatusFn'
func (e *visionEngine) GetStatus(sid types.SequenceID) types.AttackStatus {
	return e.attackerStatus[sid]
}
