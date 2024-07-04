package main

import (
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/projectdiscovery/cdncheck"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/ipranger"
	"github.com/projectdiscovery/networkpolicy"
	"github.com/runetale/runevision/router"
	"golang.org/x/exp/maps"
	"golang.org/x/net/proxy"
)

func isPrivileged() bool {
	return os.Geteuid() == 0
}

// Protocol
type Protocol int

const (
	TCP Protocol = iota
	UDP
	ARP
)

func (p Protocol) String() string {
	switch p {
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	case ARP:
		return "arp"
	default:
		panic("uknown type")
	}
}

func (p Protocol) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.String() + `"`), nil
}

// Pkg Result
type PkgResult struct {
	ipv4 string
	ipv6 string
	port *Port
}

// Port
type Port struct {
	Port     int      `json:"port"`
	Protocol Protocol `json:"protocol"`
	TLS      bool     `json:"tls"`
}

func (p *Port) String() string {
	return fmt.Sprintf("%d-%d-%v", p.Port, p.Protocol, p.TLS)
}

// ListenHandler
var (
	PkgRouter router.Router
)

type ListenHandler struct {
	Busy                                   bool
	Phase                                  *Phase
	SourceHW                               net.HardwareAddr
	SourceIp4                              net.IP
	SourceIP6                              net.IP
	Port                                   int
	TcpConn4, UdpConn4, TcpConn6, UdpConn6 *net.IPConn
	TcpChan, UdpChan, HostDiscoveryChan    chan *PkgResult
}

// Host Results
type ResultFn func(*HostResult)

type HostResult struct {
	Host  string
	IP    string
	Ports []*Port
}

// Scan Result
type Result struct {
	sync.RWMutex
	ipPorts map[string]map[string]*Port
	ips     map[string]struct{}
	skipped map[string]struct{}
}

func NewResult() *Result {
	ipPorts := make(map[string]map[string]*Port)
	ips := make(map[string]struct{})
	skipped := make(map[string]struct{})
	return &Result{ipPorts: ipPorts, ips: ips, skipped: skipped}
}

func (r *Result) GetIPs() chan string {
	r.Lock()

	out := make(chan string)

	go func() {
		defer close(out)
		defer r.Unlock()

		for ip := range r.ips {
			out <- ip
		}
	}()

	return out
}

func (r *Result) HasIPS() bool {
	r.RLock()
	defer r.RUnlock()

	return len(r.ips) > 0
}

// GetIpsPorts returns the ips and ports
func (r *Result) GetIPsPorts() chan *HostResult {
	r.RLock()

	out := make(chan *HostResult)

	go func() {
		defer close(out)
		defer r.RUnlock()

		for ip, ports := range r.ipPorts {
			if r.HasSkipped(ip) {
				continue
			}
			out <- &HostResult{IP: ip, Ports: maps.Values(ports)}
		}
	}()

	return out
}

func (r *Result) HasIPsPorts() bool {
	r.RLock()
	defer r.RUnlock()

	return len(r.ipPorts) > 0
}

// AddPort to a specific ip
func (r *Result) AddPort(ip string, p *Port) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.ipPorts[ip]; !ok {
		r.ipPorts[ip] = make(map[string]*Port)
	}

	r.ipPorts[ip][p.String()] = p
	r.ips[ip] = struct{}{}
}

// SetPorts for a specific ip
func (r *Result) SetPorts(ip string, ports []*Port) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.ipPorts[ip]; !ok {
		r.ipPorts[ip] = make(map[string]*Port)
	}

	for _, p := range ports {
		r.ipPorts[ip][p.String()] = p
	}
	r.ips[ip] = struct{}{}
}

// IPHasPort checks if an ip has a specific port
func (r *Result) IPHasPort(ip string, p *Port) bool {
	r.RLock()
	defer r.RUnlock()

	ipPorts, hasports := r.ipPorts[ip]
	if !hasports {
		return false
	}
	_, hasport := ipPorts[p.String()]

	return hasport
}

// AddIp adds an ip to the results
func (r *Result) AddIp(ip string) {
	r.Lock()
	defer r.Unlock()

	r.ips[ip] = struct{}{}
}

// HasIP checks if an ip has been seen
func (r *Result) HasIP(ip string) bool {
	r.RLock()
	defer r.RUnlock()

	_, ok := r.ips[ip]
	return ok
}

func (r *Result) IsEmpty() bool {
	return r.Len() == 0
}

func (r *Result) Len() int {
	r.RLock()
	defer r.RUnlock()

	return len(r.ips)
}

// GetPortCount returns the number of ports discovered for an ip
func (r *Result) GetPortCount(host string) int {
	r.RLock()
	defer r.RUnlock()

	return len(r.ipPorts[host])
}

// AddSkipped adds an ip to the skipped list
func (r *Result) AddSkipped(ip string) {
	r.Lock()
	defer r.Unlock()

	r.skipped[ip] = struct{}{}
}

// HasSkipped checks if an ip has been skipped
func (r *Result) HasSkipped(ip string) bool {
	r.RLock()
	defer r.RUnlock()

	_, ok := r.skipped[ip]
	return ok
}

type TCPSequencer struct {
	current uint32
}

// NewTCPSequencer creates a new linear tcp sequenc enumber generator
func NewTCPSequencer() *TCPSequencer {
	// Start the sequence with math.MaxUint32, which will then wrap around
	// when incremented starting the sequence with 0 as desired.
	return &TCPSequencer{current: math.MaxUint32}
}

// Next returns the next number in the sequence of tcp sequence numbers
func (t *TCPSequencer) Next() uint32 {
	value := atomic.AddUint32(&t.current, 1)
	return value
}

// Scanner
type Scanner struct {
	retries       int
	rate          int
	portThreshold int
	timeout       time.Duration
	proxyDialer   proxy.Dialer

	Ports    []*Port
	IPRanger *ipranger.IPRanger

	HostDiscoveryResults *Result
	ScanResults          *Result
	NetworkInterface     *net.Interface
	cdn                  *cdncheck.Client
	tcpsequencer         *TCPSequencer
	stream               bool
	ListenHandler        *ListenHandler
	OnReceive            ResultFn
}

// Scan Options
type ScanOptions struct {
	Timeout       time.Duration
	Retries       int
	Rate          int
	PortThreshold int
	ExcludeCdn    bool
	OutputCdn     bool
	ExcludedIps   []string
	Proxy         string
	ProxyAuth     string
	Stream        bool
	OnReceive     ResultFn
	ScanType      string
}

// State determines the internal scan state
type State int

const (
	maxRetries     = 10
	sendDelayMsec  = 10
	chanSize       = 1000  //nolint
	packetSendSize = 2500  //nolint
	snaplen        = 65536 //nolint
	readtimeout    = 1500  //nolint
)

const (
	Init State = iota
	HostDiscovery
	Scan
	Done
	Guard
)

type Phase struct {
	sync.RWMutex
	State
}

func (phase *Phase) Is(state State) bool {
	phase.RLock()
	defer phase.RUnlock()

	return phase.State == state
}

func (phase *Phase) Set(state State) {
	phase.Lock()
	defer phase.Unlock()

	phase.State = state
}

func NewListenHandler() *ListenHandler {
	return &ListenHandler{Phase: &Phase{}}
}

func Acquire(options *ScanOptions) (*ListenHandler, error) {
	// always grant to unprivileged scans or connect scan
	fmt.Println("start scan")
	if PkgRouter == nil || !isPrivileged() || options.ScanType == "c" {
		return NewListenHandler(), nil
	}

	return nil, errors.New("no free handlers")
}

func NewScanner(options *ScanOptions) (*Scanner, error) {
	iprang, err := ipranger.New()
	if err != nil {
		return nil, err
	}

	var nPolicyOptions networkpolicy.Options
	nPolicyOptions.DenyList = append(nPolicyOptions.DenyList, options.ExcludedIps...)
	nPolicy, err := networkpolicy.New(nPolicyOptions)
	if err != nil {
		return nil, err
	}
	iprang.Np = nPolicy
	scanner := &Scanner{
		timeout:       options.Timeout,
		retries:       options.Retries,
		rate:          options.Rate,
		portThreshold: options.PortThreshold,
		tcpsequencer:  NewTCPSequencer(),
		IPRanger:      iprang,
		OnReceive:     options.OnReceive,
	}

	scanner.HostDiscoveryResults = NewResult()
	scanner.ScanResults = NewResult()
	if options.ExcludeCdn || options.OutputCdn {
		scanner.cdn = cdncheck.New()
	}

	var auth *proxy.Auth = nil

	if options.ProxyAuth != "" && strings.Contains(options.ProxyAuth, ":") {
		credentials := strings.SplitN(options.ProxyAuth, ":", 2)
		var user, password string
		user = credentials[0]
		if len(credentials) == 2 {
			password = credentials[1]
		}
		auth = &proxy.Auth{User: user, Password: password}
	}

	if options.Proxy != "" {
		proxyDialer, err := proxy.SOCKS5("tcp", options.Proxy, auth, &net.Dialer{Timeout: options.Timeout})
		if err != nil {
			return nil, err
		}
		scanner.proxyDialer = proxyDialer
	}

	scanner.stream = options.Stream
acquire:
	if handler, err := Acquire(options); err != nil {
		// automatically fallback to connect scan
		if err != nil && options.ScanType == "s" {
			gologger.Info().Msgf("syn scan is not possible, falling back to connect scan")
			options.ScanType = "c"
			goto acquire
		}
		return scanner, err
	} else {
		scanner.ListenHandler = handler
	}

	return scanner, err
}

func init() {
	if r, err := router.New(); err != nil {
		gologger.Error().Msgf("could not initialize router: %s\n", err)
	} else {
		PkgRouter = r
	}
}

func main() {
	options := ParseOptions()
	fmt.Println(options)

	// ctx, cancel := context.WithCancel(context.TODO())
	// defer cancel()

	if isPrivileged() {
		fmt.Println("root mode")
	} else {
		fmt.Println("please root")
		return
	}

}
