package vsengine

import (
	"fmt"

	"github.com/runetale/thor/domain/requests"
	"github.com/runetale/thor/hack"
	"github.com/runetale/thor/types"
	"github.com/runetale/thor/utility/atomics"
	"github.com/runetale/thor/vsd"
)

type Attacker struct {
	vsd        map[types.SequenceID]*vsd.VisionSystem
	isVsd      atomics.AtomicValue[func(types.SequenceID) bool]
	sequenceID types.SequenceID

	updateAttackStatusFn func(sid types.SequenceID, status types.AttackStatus)

	// Attackerのルーチンを完全に終了
	closeCh chan struct{}

	// vsdが使用するid毎の攻撃完了通知チャネル
	doneAttackCh map[types.SequenceID]chan struct{}
}

func NewAttacker() *Attacker {
	ata := &Attacker{
		vsd: make(map[types.SequenceID]*vsd.VisionSystem, 0),
	}
	return ata
}

func (a *Attacker) SetUpdateAttackStatusFn(updateStatusFn func(sid types.SequenceID, status types.AttackStatus)) {
	a.updateAttackStatusFn = updateStatusFn
}

func (a *Attacker) Start() {
	a.isVsd.Store(func(s types.SequenceID) bool { return false })
	go func() {
		select {
		case <-a.closeCh:
			// continue
		}
	}()
}

func (a *Attacker) Close() error {
	close(a.closeCh)
	return nil
}

func (a *Attacker) setSequentialID(sid types.SequenceID) func(s types.SequenceID) bool {
	return func(s types.SequenceID) bool {
		a.sequenceID = sid
		return true
	}
}

type attackParams struct {
	scanParams hack.ScanParams
	hfParams   hack.HostFinderParams
	subParams  hack.SubfinderParams
	hxparams   hack.HttpxParams
	kaParams   hack.KatanaParams
	nuParams   hack.NucleiParams
}

func (a *Attacker) parseAttackerParamsFromRequest(req *requests.HackDoScanRequest) *attackParams {
	return &attackParams{
		scanParams: hack.ScanParams{
			TargetHost:  req.TargetHost,
			ScanType:    hack.ScanType(req.PortScanType),
			TargetPorts: req.TargetPorts,
		},
		hfParams: hack.HostFinderParams{
			Queries:  req.Queries,
			Limit:    req.RateLimit,
			MaxRetry: req.Retry,
			Timeout:  req.Timeout,
		},
		subParams: hack.SubfinderParams{
			TargetHost:         req.TargetHost,
			Threads:            req.Threads,
			Timeout:            req.Timeout,
			MaxEnumerationTime: requests.MaxEnumerationTime,
		},
		hxparams: hack.HttpxParams{
			Methods:     req.HTTPMethods,
			TargetHosts: req.TargetHost,
		},
		kaParams: hack.KatanaParams{
			TargetHost:  req.TargetHost,
			MaxDepth:    req.MaxDepth,
			FieldScope:  string(req.FieldScope),
			Timeout:     requests.Timeout,
			Concurrency: requests.Threads,
			Parallelism: req.Parallelism,
			Delay:       req.Delay,
			RateLimit:   req.RateLimit,
			Strategy:    string(req.Strategy),
		},
		nuParams: hack.NucleiParams{
			TemplateTags: req.TemplateTags,
			TargetHosts:  req.TargetHost,
		},
	}
}

func (a *Attacker) setVsdAttackerParams(vsd *vsd.VisionSystem, params *attackParams) error {
	// dsl
	dsl := vsd.Dsl.Get()
	err := dsl.SetParams(params.scanParams)
	if err != nil {
		return err
	}

	// hostfiner
	// hostFinder := i.vsd.HostFinder.Get()
	// hostFinder.SetParams(hfParams)
	// hostFinder.Start()

	// port scan
	ps := vsd.PortScanner.Get()
	err = ps.SetParams(params.scanParams)
	if err != nil {
		return err
	}

	// subfinder
	sf := vsd.SubDomainFinder.Get()
	err = sf.SetParams(params.subParams)
	if err != nil {
		return err
	}

	// httpx
	hx := vsd.Httpx.Get()
	err = hx.SetParams(params.hxparams)
	if err != nil {
		return err
	}

	// katana
	ka := vsd.Katana.Get()
	err = ka.SetParams(params.kaParams)
	if err != nil {
		return err
	}

	// nuclei
	nu := vsd.Nuclei.Get()
	err = nu.SetParams(params.nuParams)
	if err != nil {
		return err
	}
	return nil
}

func (a *Attacker) Reconfig(seqID types.SequenceID, vsd *vsd.VisionSystem, requests *requests.HackDoScanRequest) error {
	a.isVsd.Store(a.setSequentialID(seqID))

	attackerParams := a.parseAttackerParamsFromRequest(requests)
	err := a.setVsdAttackerParams(vsd, attackerParams)
	if err != nil {
		return err
	}

	a.vsd[seqID] = vsd

	a.updateAttackStatusFn(seqID, types.STARTING)

	a.attack(seqID)

	return nil
}

func (a *Attacker) attack(seqID types.SequenceID) {
	if _, ok := a.isVsd.LoadOk(); !ok {
		fmt.Println("thor system has no set.")
	}

	// todo (snt) 終わったことを通知させる
	for se, v := range a.vsd {
		a.updateAttackStatusFn(se, types.SCANNING)
		fmt.Printf("scanning %s\n", se)
		v.AllRun()
	}

	go func() {
		select {
		case sequenceID := <-a.doneAttackCh[seqID]:
			fmt.Printf("completed scan, sequence number is `%s`", sequenceID)
			a.updateAttackStatusFn(seqID, types.COMPLETED)
		}
	}()

}
