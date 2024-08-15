// このパッケージはvision engineやvision database(cvemap)などにアクセスするためのAPIを提供する
package backend

import (
	"context"
	"encoding/json"
	"flag"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/runetale/runevision/utility"
)

var (
	currentCvemapBinary = flag.String("current", "./cvemap", "Current Branch Cvemap Binary")
)

// type CPE struct {
// 	CPE     string `json:"cpe"`
// 	Vendor  string `json:"vendor"`
// 	Product string `json:"product"`
// }

// type EPSS struct {
// 	EPSSScore      float64 `json:"epss_score"`
// 	EPSSPercentile float64 `json:"epss_percentile"`
// }

// type HackerOne struct {
// 	Rank  int `json:"rank"`
// 	Count int `json:"count"`
// }

// type CVSS2 struct {
// 	Score    float64 `json:"score"`
// 	Vector   string  `json:"vector"`
// 	Severity string  `json:"severity"`
// }

// type CVSSMetrics struct {
// 	CVSS2 CVSS2 `json:"cvss2"`
// }

// type Weakness struct {
// 	CWEID   string `json:"cwe_id"`
// 	CWEName string `json:"cwe_name"`
// }

// type CVE struct {
// 	CPE            CPE         `json:"cpe"`
// 	EPSS           EPSS        `json:"epss"`
// 	CVEID          string      `json:"cve_id"`
// 	IsOSS          bool        `json:"is_oss"`
// 	IsPOC          bool        `json:"is_poc"`
// 	Assignee       string      `json:"assignee"`
// 	Severity       string      `json:"severity"`
// 	HackerOne      HackerOne   `json:"hackerone"`
// 	IsRemote       bool        `json:"is_remote"`
// 	Reference      []string    `json:"reference"`
// 	CVSSScore      float64     `json:"cvss_score"`
// 	UpdatedAt      CustomTime  `json:"updated_at"`
// 	Weaknesses     []Weakness  `json:"weaknesses"`
// 	AgeInDays      int         `json:"age_in_days"`
// 	IsTemplate     bool        `json:"is_template"`
// 	VulnStatus     string      `json:"vuln_status"`
// 	CVSSMetrics    CVSSMetrics `json:"cvss_metrics"`
// 	IsExploited    bool        `json:"is_exploited"`
// 	PublishedAt    CustomTime  `json:"published_at"`
// 	VulnerableCPE  []string    `json:"vulnerable_cpe"`
// 	CVEDescription string      `json:"cve_description"`
// 	VendorAdvisory string      `json:"vendor_advisory"`
// }

// type CustomTime struct {
// 	time.Time
// }

// func (ct *CustomTime) UnmarshalJSON(b []byte) error {
// 	s := string(b[1 : len(b)-1])
// 	layout := "2006-01-02T15:04:05.000"
// 	t, err := time.Parse(layout, s)
// 	if err != nil {
// 		return err
// 	}
// 	ct.Time = t
// 	return nil
// }
// func main() {
// 	err := execute()
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func execute() error {
// 	currentOutput, err := testutils.RunCvemapBinaryAndGetResults(*currentCvemapBinary, true, []string{"-id", "CVE-1999-0027", "-j", "-silent"})
// 	if err != nil {
// 		return errors.Wrap(err, "could not run cvemap test")
// 	}

// 	if len(currentOutput) == 0 {
// 		return errors.New("no output from cvemap")
// 	}

// 	result := strings.Join(currentOutput[1:len(currentOutput)-1], "\n")

// 	// var data CVEData
// 	var data CVE
// 	err = json.Unmarshal([]byte(result), &data)
// 	if err != nil {
// 		return errors.Wrap(err, "error unmarshaling JSON")
// 	}

// 	fmt.Printf("%+v\n", data)

// 	return nil
// }

type LocalBackendServer struct {
	lb atomic.Pointer[LocalBackend]

	backendWaiter backendWaiter

	mu sync.Mutex

	logger *utility.Logger
}

func New(logger *utility.Logger) *LocalBackendServer {
	return &LocalBackendServer{
		logger: logger,
	}
}

type backendWaiter utility.HandleSet[context.CancelFunc]

func (waiter *backendWaiter) add(mu *sync.Mutex, ctx context.Context) (ready <-chan struct{}, cleanup func()) {
	ctx, cancel := context.WithCancel(ctx)
	hs := (*utility.HandleSet[context.CancelFunc])(waiter)
	mu.Lock()
	h := hs.Add(cancel)
	mu.Unlock()
	return ctx.Done(), func() {
		mu.Lock()
		delete(*hs, h)
		mu.Unlock()
		cancel()
	}
}

func (waiter backendWaiter) wakeAll() {
	for _, cancel := range waiter {
		cancel()
	}
}

func (s *LocalBackendServer) awaitBackend(ctx context.Context) (_ *LocalBackend, ok bool) {
	lb := s.lb.Load()
	if lb != nil {
		return lb, true
	}

	ready, cleanup := s.backendWaiter.add(&s.mu, ctx)
	defer cleanup()

	lb = s.lb.Load()
	if lb != nil {
		return lb, true
	}

	<-ready
	lb = s.lb.Load()
	return lb, lb != nil
}

func (s *LocalBackendServer) serveServerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var res struct {
		Error string `json:"error,omitempty"`
	}

	lb := s.lb.Load()
	if lb == nil {
		res.Error = "backend not ready"
	}
	json.NewEncoder(w).Encode(res)
}

func (s *LocalBackendServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && r.URL.Path == "/server-status" {
		s.serveServerStatus(w, r)
		return
	}

	ctx := r.Context()
	_, ok := s.awaitBackend(ctx)
	if !ok {
		http.Error(w, "no backend", http.StatusServiceUnavailable)
		return
	}

	// if strings.HasPrefix(r.URL.Path, "/localapi/") {
	// 	lah := hashigo.NewHandler(hb)
	// 	lah.ServeHTTP(w, r)
	// 	return
	// }

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	io.WriteString(w, "<html><title>RuneVision</title><body><h1>RuneVision</h1>hi, i'm Vision backend daemon server.\n")
}

func (s *LocalBackendServer) Run(ctx context.Context, ln net.Listener) error {
	defer func() {
		if hb := s.lb.Load(); hb != nil {
			hb.Shutdown()
		}
	}()

	runDone := make(chan struct{})
	defer close(runDone)

	go func() {
		select {
		case <-ctx.Done():
		case <-runDone:
		}
		ln.Close()
	}()

	hs := &http.Server{
		Handler:     http.HandlerFunc(s.serveHTTP),
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		IdleTimeout: 6 * time.Second,
	}

	if err := hs.Serve(ln); err != nil {
		if err := ctx.Err(); err != nil {
			return err
		}
		return err
	}
	return nil

}

func (s *LocalBackendServer) SetLocalBackend(lb *LocalBackend) {
	if lb == nil {
		panic("nil LocalBackend")
	}

	if !s.lb.CompareAndSwap(nil, lb) {
		panic("already set")
	}

	s.mu.Lock()
	s.backendWaiter.wakeAll()
	s.mu.Unlock()

}
