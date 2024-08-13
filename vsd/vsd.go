// vst (short for "vision daemon")
package vsd

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/runetale/runevision/hack"
	"github.com/runetale/runevision/types"
)

type VisionSystem struct {
	Dsl             SubSystem[*hack.DslRunner]
	HostFinder      SubSystem[*hack.HostFinder]
	Httpx           SubSystem[*hack.HttpxRunner]
	Katana          SubSystem[*hack.KatanaRunner]
	Nuclei          SubSystem[*hack.NucleiRunner]
	PortScanner     SubSystem[*hack.PortScanner]
	SubDomainFinder SubSystem[*hack.SubfinderRunner]

	// root権限かどうか
	// port scanなどのパフォーマンスが向上する
	IsPrivileged bool
}

func (s *VisionSystem) Set(v any) {
	switch v := v.(type) {
	case *hack.DslRunner:
		s.Dsl.Set(v)
	case *hack.HostFinder:
		s.HostFinder.Set(v)
	case *hack.HttpxRunner:
		s.Httpx.Set(v)
	case *hack.KatanaRunner:
		s.Katana.Set(v)
	case *hack.NucleiRunner:
		s.Nuclei.Set(v)
	case *hack.PortScanner:
		s.PortScanner.Set(v)
	case *hack.SubfinderRunner:
		s.SubDomainFinder.Set(v)
	default:
		panic(fmt.Sprintf("unknown type %T", v))
	}
}

func (s *VisionSystem) Run(v any) error {
	switch v := v.(type) {
	case *hack.DslRunner:
		dsl, ok := s.Dsl.GetOK()
		if !ok {
			return errors.New("cannot get dsl")
		}
		err := dsl.Start()
		if err != nil {
			return err
		}
	case *hack.HostFinder:
		hf, ok := s.HostFinder.GetOK()
		if !ok {
			return errors.New("cannot get dsl")
		}
		err := hf.Start()
		if err != nil {
			return err
		}
	case *hack.HttpxRunner:
		hx, ok := s.Httpx.GetOK()
		if !ok {
			return errors.New("cannot get dsl")
		}
		err := hx.Start()
		if err != nil {
			return err
		}
	case *hack.KatanaRunner:
		ka, ok := s.Katana.GetOK()
		if !ok {
			return errors.New("cannot get dsl")
		}
		err := ka.Start()
		if err != nil {
			return err
		}
	case *hack.NucleiRunner:
		nu, ok := s.Nuclei.GetOK()
		if !ok {
			return errors.New("cannot get dsl")
		}
		err := nu.Start()
		if err != nil {
			return err
		}
	case *hack.PortScanner:
		po, ok := s.PortScanner.GetOK()
		if !ok {
			return errors.New("cannot get dsl")
		}
		err := po.Start(s.IsPrivileged)
		if err != nil {
			return err
		}
	case *hack.SubfinderRunner:
		sf, ok := s.SubDomainFinder.GetOK()
		if !ok {
			return errors.New("cannot get dsl")
		}
		err := sf.Start()
		if err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("unknown type %T", v))
	}
	return nil
}

// note: (snt) Should check that it is Set before running AllRun
// Must asynchronous run
func (s *VisionSystem) AllRun(chan types.SequenceID) {
	dsl := s.Dsl.Get()
	go dsl.Start()

	// 有料なのでオプション
	// 	hf, ok := s.HostFinder.GetOK()
	// 	if ok {
	// 		go hf.Start()
	// 	}

	hx := s.Httpx.Get()
	go hx.Start()

	k := s.Katana.Get()
	go k.Start()

	n := s.Nuclei.Get()
	go n.Start()

	ps := s.PortScanner.Get()
	go ps.Start(s.IsPrivileged)

	sdf := s.SubDomainFinder.Get()
	go sdf.Start()
}

type SubSystem[T any] struct {
	set bool
	v   T
}

func (p *SubSystem[T]) Set(v T) {
	if p.set {
		var oldVal any = p.v
		var newVal any = v
		if oldVal == newVal {
			return
		}

		var z *T
		panic(fmt.Sprintf("%v is already set", reflect.TypeOf(z).Elem().String()))
	}
	p.v = v
	p.set = true
}

// Get returns the value of p, panicking if it hasn't been set.
func (p *SubSystem[T]) Get() T {
	if !p.set {
		var z *T
		panic(fmt.Sprintf("%v is not set", reflect.TypeOf(z).Elem().String()))
	}
	return p.v
}

// GetOK returns the value of p (if any) and whether it's been set.
func (p *SubSystem[T]) GetOK() (_ T, ok bool) {
	return p.v, p.set
}
