package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Mzack9999/gcache"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/projectdiscovery/cdncheck"
	"github.com/projectdiscovery/clistats"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/ipranger"
	"github.com/projectdiscovery/networkpolicy"
	"github.com/remeh/sizedwaitgroup"
	"github.com/runetale/runevision/router"
	"golang.org/x/exp/maps"
	"golang.org/x/net/proxy"
)

func init() {
	if r, err := router.New(); err != nil {
		gologger.Error().Msgf("could not initialize router: %s\n", err)
	} else {
		PkgRouter = r
	}
}

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
type ScanResult struct {
	sync.RWMutex
	ipPorts map[string]map[string]*Port
	ips     map[string]struct{}
	skipped map[string]struct{}
}

func NewScanResult() *ScanResult {
	ipPorts := make(map[string]map[string]*Port)
	ips := make(map[string]struct{})
	skipped := make(map[string]struct{})
	return &ScanResult{ipPorts: ipPorts, ips: ips, skipped: skipped}
}

func (r *ScanResult) GetIPs() chan string {
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

func (r *ScanResult) HasIPS() bool {
	r.RLock()
	defer r.RUnlock()

	return len(r.ips) > 0
}

// GetIpsPorts returns the ips and ports
func (r *ScanResult) GetIPsPorts() chan *HostResult {
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

func (r *ScanResult) HasIPsPorts() bool {
	r.RLock()
	defer r.RUnlock()

	return len(r.ipPorts) > 0
}

// AddPort to a specific ip
func (r *ScanResult) AddPort(ip string, p *Port) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.ipPorts[ip]; !ok {
		r.ipPorts[ip] = make(map[string]*Port)
	}

	r.ipPorts[ip][p.String()] = p
	r.ips[ip] = struct{}{}
}

// SetPorts for a specific ip
func (r *ScanResult) SetPorts(ip string, ports []*Port) {
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
func (r *ScanResult) IPHasPort(ip string, p *Port) bool {
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
func (r *ScanResult) AddIp(ip string) {
	r.Lock()
	defer r.Unlock()

	r.ips[ip] = struct{}{}
}

// HasIP checks if an ip has been seen
func (r *ScanResult) HasIP(ip string) bool {
	r.RLock()
	defer r.RUnlock()

	_, ok := r.ips[ip]
	return ok
}

func (r *ScanResult) IsEmpty() bool {
	return r.Len() == 0
}

func (r *ScanResult) Len() int {
	r.RLock()
	defer r.RUnlock()

	return len(r.ips)
}

// GetPortCount returns the number of ports discovered for an ip
func (r *ScanResult) GetPortCount(host string) int {
	r.RLock()
	defer r.RUnlock()

	return len(r.ipPorts[host])
}

// AddSkipped adds an ip to the skipped list
func (r *ScanResult) AddSkipped(ip string) {
	r.Lock()
	defer r.Unlock()

	r.skipped[ip] = struct{}{}
}

// HasSkipped checks if an ip has been skipped
func (r *ScanResult) HasSkipped(ip string) bool {
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

	HostDiscoveryResults *ScanResult
	ScanResults          *ScanResult
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

	scanner.HostDiscoveryResults = NewScanResult()
	scanner.ScanResults = NewScanResult()
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

func (s *Scanner) Close() {
	s.ListenHandler.Busy = false
	s.ListenHandler = nil
}

const portListStrParts = 2

// List of default ports
const (
	Full        = "1-65535"
	NmapTop100  = "7,9,13,21-23,25-26,37,53,79-81,88,106,110-111,113,119,135,139,143-144,179,199,389,427,443-445,465,513-515,543-544,548,554,587,631,646,873,990,993,995,1025-1029,1110,1433,1720,1723,1755,1900,2000-2001,2049,2121,2717,3000,3128,3306,3389,3986,4899,5000,5009,5051,5060,5101,5190,5357,5432,5631,5666,5800,5900,6000-6001,6646,7070,8000,8008-8009,8080-8081,8443,8888,9100,9999-10000,32768,49152-49157"
	NmapTop1000 = "1,3-4,6-7,9,13,17,19-26,30,32-33,37,42-43,49,53,70,79-85,88-90,99-100,106,109-111,113,119,125,135,139,143-144,146,161,163,179,199,211-212,222,254-256,259,264,280,301,306,311,340,366,389,406-407,416-417,425,427,443-445,458,464-465,481,497,500,512-515,524,541,543-545,548,554-555,563,587,593,616-617,625,631,636,646,648,666-668,683,687,691,700,705,711,714,720,722,726,749,765,777,783,787,800-801,808,843,873,880,888,898,900-903,911-912,981,987,990,992-993,995,999-1002,1007,1009-1011,1021-1100,1102,1104-1108,1110-1114,1117,1119,1121-1124,1126,1130-1132,1137-1138,1141,1145,1147-1149,1151-1152,1154,1163-1166,1169,1174-1175,1183,1185-1187,1192,1198-1199,1201,1213,1216-1218,1233-1234,1236,1244,1247-1248,1259,1271-1272,1277,1287,1296,1300-1301,1309-1311,1322,1328,1334,1352,1417,1433-1434,1443,1455,1461,1494,1500-1501,1503,1521,1524,1533,1556,1580,1583,1594,1600,1641,1658,1666,1687-1688,1700,1717-1721,1723,1755,1761,1782-1783,1801,1805,1812,1839-1840,1862-1864,1875,1900,1914,1935,1947,1971-1972,1974,1984,1998-2010,2013,2020-2022,2030,2033-2035,2038,2040-2043,2045-2049,2065,2068,2099-2100,2103,2105-2107,2111,2119,2121,2126,2135,2144,2160-2161,2170,2179,2190-2191,2196,2200,2222,2251,2260,2288,2301,2323,2366,2381-2383,2393-2394,2399,2401,2492,2500,2522,2525,2557,2601-2602,2604-2605,2607-2608,2638,2701-2702,2710,2717-2718,2725,2800,2809,2811,2869,2875,2909-2910,2920,2967-2968,2998,3000-3001,3003,3005-3007,3011,3013,3017,3030-3031,3052,3071,3077,3128,3168,3211,3221,3260-3261,3268-3269,3283,3300-3301,3306,3322-3325,3333,3351,3367,3369-3372,3389-3390,3404,3476,3493,3517,3527,3546,3551,3580,3659,3689-3690,3703,3737,3766,3784,3800-3801,3809,3814,3826-3828,3851,3869,3871,3878,3880,3889,3905,3914,3918,3920,3945,3971,3986,3995,3998,4000-4006,4045,4111,4125-4126,4129,4224,4242,4279,4321,4343,4443-4446,4449,4550,4567,4662,4848,4899-4900,4998,5000-5004,5009,5030,5033,5050-5051,5054,5060-5061,5080,5087,5100-5102,5120,5190,5200,5214,5221-5222,5225-5226,5269,5280,5298,5357,5405,5414,5431-5432,5440,5500,5510,5544,5550,5555,5560,5566,5631,5633,5666,5678-5679,5718,5730,5800-5802,5810-5811,5815,5822,5825,5850,5859,5862,5877,5900-5904,5906-5907,5910-5911,5915,5922,5925,5950,5952,5959-5963,5987-5989,5998-6007,6009,6025,6059,6100-6101,6106,6112,6123,6129,6156,6346,6389,6502,6510,6543,6547,6565-6567,6580,6646,6666-6669,6689,6692,6699,6779,6788-6789,6792,6839,6881,6901,6969,7000-7002,7004,7007,7019,7025,7070,7100,7103,7106,7200-7201,7402,7435,7443,7496,7512,7625,7627,7676,7741,7777-7778,7800,7911,7920-7921,7937-7938,7999-8002,8007-8011,8021-8022,8031,8042,8045,8080-8090,8093,8099-8100,8180-8181,8192-8194,8200,8222,8254,8290-8292,8300,8333,8383,8400,8402,8443,8500,8600,8649,8651-8652,8654,8701,8800,8873,8888,8899,8994,9000-9003,9009-9011,9040,9050,9071,9080-9081,9090-9091,9099-9103,9110-9111,9200,9207,9220,9290,9415,9418,9485,9500,9502-9503,9535,9575,9593-9595,9618,9666,9876-9878,9898,9900,9917,9929,9943-9944,9968,9998-10004,10009-10010,10012,10024-10025,10082,10180,10215,10243,10566,10616-10617,10621,10626,10628-10629,10778,11110-11111,11967,12000,12174,12265,12345,13456,13722,13782-13783,14000,14238,14441-14442,15000,15002-15004,15660,15742,16000-16001,16012,16016,16018,16080,16113,16992-16993,17877,17988,18040,18101,18988,19101,19283,19315,19350,19780,19801,19842,20000,20005,20031,20221-20222,20828,21571,22939,23502,24444,24800,25734-25735,26214,27000,27352-27353,27355-27356,27715,28201,30000,30718,30951,31038,31337,32768-32785,33354,33899,34571-34573,35500,38292,40193,40911,41511,42510,44176,44442-44443,44501,45100,48080,49152-49161,49163,49165,49167,49175-49176,49400,49999-50003,50006,50300,50389,50500,50636,50800,51103,51493,52673,52822,52848,52869,54045,54328,55055-55056,55555,55600,56737-56738,57294,57797,58080,60020,60443,61532,61900,62078,63331,64623,64680,65000,65129,65389"
)

// ParsePorts parses the list of ports and creates a port map
func ParsePorts(options *Options) ([]*Port, error) {
	var portsFileMap, portsCLIMap, topPortsCLIMap, portsConfigList []*Port

	// If the user has specfied a ports file, use it
	if options.PortsFile != "" {
		data, err := os.ReadFile(options.PortsFile)
		if err != nil {
			return nil, fmt.Errorf("could not read ports: %s", err)
		}
		ports, err := parsePortsList(string(data))
		if err != nil {
			return nil, fmt.Errorf("could not read ports: %s", err)
		}
		portsFileMap, err = excludePorts(options, ports)
		if err != nil {
			return nil, fmt.Errorf("could not read ports: %s", err)
		}
	}

	// If the user has specfied top ports, use them as well
	if options.TopPorts != "" {
		switch strings.ToLower(options.TopPorts) {
		case "full": // If the user has specfied full ports, use them
			var err error
			ports, err := parsePortsList(Full)
			if err != nil {
				return nil, fmt.Errorf("could not read ports: %s", err)
			}
			topPortsCLIMap, err = excludePorts(options, ports)
			if err != nil {
				return nil, fmt.Errorf("could not read ports: %s", err)
			}
		case "100": // If the user has specfied 100, use them
			ports, err := parsePortsList(NmapTop100)
			if err != nil {
				return nil, fmt.Errorf("could not read ports: %s", err)
			}
			topPortsCLIMap, err = excludePorts(options, ports)
			if err != nil {
				return nil, fmt.Errorf("could not read ports: %s", err)
			}
		case "1000": // If the user has specfied 1000, use them
			ports, err := parsePortsList(NmapTop1000)
			if err != nil {
				return nil, fmt.Errorf("could not read ports: %s", err)
			}
			topPortsCLIMap, err = excludePorts(options, ports)
			if err != nil {
				return nil, fmt.Errorf("could not read ports: %s", err)
			}
		default:
			return nil, errors.New("invalid top ports option")
		}
	}

	// If the user has specfied ports option, use them too
	if options.Ports != "" {
		// "-" equals to all ports
		if options.Ports == "-" {
			// Parse the custom ports list provided by the user
			options.Ports = "1-65535"
		}
		ports, err := parsePortsList(options.Ports)
		if err != nil {
			return nil, fmt.Errorf("could not read ports: %s", err)
		}
		portsCLIMap, err = excludePorts(options, ports)
		if err != nil {
			return nil, fmt.Errorf("could not read ports: %s", err)
		}
	}

	// merge all the specified ports (meaningless if "all" is used)
	ports := merge(portsFileMap, portsCLIMap, topPortsCLIMap, portsConfigList)

	// By default scan top 100 ports only
	if len(ports) == 0 {
		portsList, err := parsePortsList(NmapTop100)
		if err != nil {
			return nil, fmt.Errorf("could not read ports: %s", err)
		}
		m, err := excludePorts(options, portsList)
		if err != nil {
			return nil, err
		}
		return m, nil
	}

	return ports, nil
}

// excludePorts excludes the list of ports from the exclusion list
func excludePorts(options *Options, ports []*Port) ([]*Port, error) {
	if options.ExcludePorts == "" {
		return ports, nil
	}

	var filteredPorts []*Port

	// Exclude the ports specified by the user in exclusion list
	excludedPortsCLI, err := parsePortsList(options.ExcludePorts)
	if err != nil {
		return nil, fmt.Errorf("could not read exclusion ports: %s", err)
	}

	for _, port := range ports {
		found := false
		for _, excludedPort := range excludedPortsCLI {
			if excludedPort.Port == port.Port && excludedPort.Protocol == port.Protocol {
				found = true
				break
			}
		}
		if !found {
			filteredPorts = append(filteredPorts, port)
		}
	}
	return filteredPorts, nil
}

func parsePortsSlice(ranges []string) ([]*Port, error) {
	var ports []*Port
	for _, r := range ranges {
		r = strings.TrimSpace(r)

		portProtocol := TCP
		if strings.HasPrefix(r, "u:") {
			portProtocol = UDP
			r = strings.TrimPrefix(r, "u:")
		}

		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != portListStrParts {
				return nil, fmt.Errorf("invalid port selection segment: '%s'", r)
			}

			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid port number: '%s'", parts[0])
			}

			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid port number: '%s'", parts[1])
			}

			if p1 > p2 || p2 > 65535 {
				return nil, fmt.Errorf("invalid port range: %d-%d", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				port := &Port{Port: i, Protocol: portProtocol}
				ports = append(ports, port)
			}
		} else {
			portNumber, err := strconv.Atoi(r)
			if err != nil || portNumber > 65535 {
				return nil, fmt.Errorf("invalid port number: '%s'", r)
			}
			port := &Port{Port: portNumber, Protocol: portProtocol}
			ports = append(ports, port)
		}
	}

	// dedupe ports
	seen := make(map[string]struct{})
	var dedupedPorts []*Port
	for _, port := range ports {
		if _, ok := seen[port.String()]; ok {
			continue
		}
		seen[port.String()] = struct{}{}
		dedupedPorts = append(dedupedPorts, port)
	}

	return dedupedPorts, nil
}

func parsePortsList(data string) ([]*Port, error) {
	return parsePortsSlice(strings.Split(data, ","))
}

func merge(slices ...[]*Port) []*Port {
	var result []*Port
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

type Target struct {
	Ip   string
	Cidr string
	Fqdn string
	Port string
}

type Runner struct {
	options       *Options
	targetsFile   string
	scanner       *Scanner
	wgscan        sizedwaitgroup.SizedWaitGroup
	dnsclient     *dnsx.DNSX
	stats         *clistats.Statistics
	streamChannel chan Target

	unique gcache.Cache[string, struct{}]
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func ReadFile(filename string) (chan string, error) {
	if !FileExists(filename) {
		return nil, errors.New("file doesn't exist")
	}
	out := make(chan string)
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		if err != nil {
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}
func (r *Runner) ShowScanResultOnExit() {
	r.handleOutput(r.scanner.ScanResults)
	// err := r.handleNmap()
	// if err != nil {
	// 	gologger.Fatal().Msgf("Could not run enumeration: %s\n", err)
	// }
}

func FolderExists(foldername string) bool {
	info, err := os.Stat(foldername)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (r *Runner) handleOutput(scanResults *ScanResult) {
	var (
		file   *os.File
		err    error
		output string
	)

	// In case the user has given an output file, write all the found
	// ports to the output file.
	if r.options.Output != "" {
		output = r.options.Output

		// create path if not existing
		outputFolder := filepath.Dir(output)
		if FolderExists(outputFolder) {
			mkdirErr := os.MkdirAll(outputFolder, 0700)
			if mkdirErr != nil {
				gologger.Error().Msgf("Could not create output folder %s: %s\n", outputFolder, mkdirErr)
				return
			}
		}

		file, err = os.Create(output)
		if err != nil {
			gologger.Error().Msgf("Could not create file %s: %s\n", output, err)
			return
		}
		defer file.Close()
	}
	csvFileHeaderEnabled := true

	switch {
	case scanResults.HasIPsPorts():
		for hostResult := range scanResults.GetIPsPorts() {
			dt, err := r.scanner.IPRanger.GetHostsByIP(hostResult.IP)
			if err != nil {
				continue
			}

			if !ipMatchesIpVersions(hostResult.IP, r.options.IPVersion...) {
				continue
			}

			// recover hostnames from ip:port combination
			for _, p := range hostResult.Ports {
				ipPort := net.JoinHostPort(hostResult.IP, fmt.Sprint(p.Port))
				if dtOthers, ok := r.scanner.IPRanger.Hosts.Get(ipPort); ok {
					if otherName, _, err := net.SplitHostPort(string(dtOthers)); err == nil {
						// replace bare ip:port with host
						for idx, ipCandidate := range dt {
							if IsIP(ipCandidate) {
								dt[idx] = otherName
							}
						}
					}
				}
			}

			buffer := bytes.Buffer{}
			for _, host := range dt {
				buffer.Reset()
				if host == "ip" {
					host = hostResult.IP
				}
				isCDNIP, cdnName, _ := r.scanner.CdnCheck(hostResult.IP)
				gologger.Info().Msgf("Found %d ports on host %s (%s)\n", len(hostResult.Ports), host, hostResult.IP)
				// file output
				if file != nil {
					if r.options.JSON {
						err = WriteJSONOutput(host, hostResult.IP, hostResult.Ports, r.options.OutputCDN, isCDNIP, cdnName, file)
					} else if r.options.CSV {
						err = WriteCsvOutput(host, hostResult.IP, hostResult.Ports, r.options.OutputCDN, isCDNIP, cdnName, csvFileHeaderEnabled, file)
					} else {
						err = WriteHostOutput(host, hostResult.Ports, r.options.OutputCDN, cdnName, file)
					}
					if err != nil {
						gologger.Error().Msgf("Could not write results to file %s for %s: %s\n", output, host, err)
					}
				}

				if r.options.OnResult != nil {
					r.options.OnResult(&HostResult{Host: host, IP: hostResult.IP, Ports: hostResult.Ports})
				}
			}
			csvFileHeaderEnabled = false
		}
	case scanResults.HasIPS():
		for hostIP := range scanResults.GetIPs() {
			dt, err := r.scanner.IPRanger.GetHostsByIP(hostIP)
			if err != nil {
				continue
			}
			if !ipMatchesIpVersions(hostIP, r.options.IPVersion...) {
				continue
			}

			buffer := bytes.Buffer{}
			writer := csv.NewWriter(&buffer)
			for _, host := range dt {
				buffer.Reset()
				if host == "ip" {
					host = hostIP
				}
				isCDNIP, cdnName, _ := r.scanner.CdnCheck(hostIP)
				gologger.Info().Msgf("Found alive host %s (%s)\n", host, hostIP)
				// console output
				if r.options.JSON || r.options.CSV {
					data := &Result{IP: hostIP, TimeStamp: time.Now().UTC()}
					if r.options.OutputCDN {
						data.IsCDNIP = isCDNIP
						data.CDNName = cdnName
					}
					if host != hostIP {
						data.Host = host
					}
				}
				if r.options.JSON {
					gologger.Silent().Msgf("%s", buffer.String())
				} else if r.options.CSV {
					writer.Flush()
					gologger.Silent().Msgf("%s", buffer.String())
				} else {
					if r.options.OutputCDN && isCDNIP {
						gologger.Silent().Msgf("%s [%s]\n", host, cdnName)
					} else {
						gologger.Silent().Msgf("%s\n", host)
					}
				}
				// file output
				if file != nil {
					if r.options.JSON {
						err = WriteJSONOutput(host, hostIP, nil, r.options.OutputCDN, isCDNIP, cdnName, file)
					} else if r.options.CSV {
						err = WriteCsvOutput(host, hostIP, nil, r.options.OutputCDN, isCDNIP, cdnName, csvFileHeaderEnabled, file)
					} else {
						err = WriteHostOutput(host, nil, r.options.OutputCDN, cdnName, file)
					}
					if err != nil {
						gologger.Error().Msgf("Could not write results to file %s for %s: %s\n", output, host, err)
					}
				}

				if r.options.OnResult != nil {
					r.options.OnResult(&HostResult{Host: host, IP: hostIP})
				}
			}
			csvFileHeaderEnabled = false
		}
	}
}

func (r *Runner) parseExcludedIps(options *Options) ([]string, error) {
	var excludedIps []string
	if options.ExcludeIps != "" {
		for _, host := range strings.Split(options.ExcludeIps, ",") {
			ips, err := r.getExcludeItems(host)
			if err != nil {
				return nil, err
			}
			excludedIps = append(excludedIps, ips...)
		}
	}

	if options.ExcludeIpsFile != "" {
		cdata, err := ReadFile(options.ExcludeIpsFile)
		if err != nil {
			return excludedIps, err
		}
		for host := range cdata {
			ips, err := r.getExcludeItems(host)
			if err != nil {
				return nil, err
			}
			excludedIps = append(excludedIps, ips...)
		}
	}

	return excludedIps, nil
}

func IsIP(str string) bool {
	return net.ParseIP(str) != nil
}

func IsCIDR(str string) bool {
	_, _, err := net.ParseCIDR(str)
	return err == nil
}

func (r *Runner) host2ips(target string) (targetIPsV4 []string, targetIPsV6 []string, err error) {
	// If the host is a Domain, then perform resolution and discover all IP
	// addresses for a given host. Else use that host for port scanning
	if !IsIP(target) {
		dnsData, err := r.dnsclient.QueryMultiple(target)
		if err != nil || dnsData == nil {
			gologger.Warning().Msgf("Could not get IP for host: %s\n", target)
			return nil, nil, err
		}
		if len(r.options.IPVersion) > 0 {
			if Contains(r.options.IPVersion, IPv4) {
				targetIPsV4 = append(targetIPsV4, dnsData.A...)
			}
			if Contains(r.options.IPVersion, IPv6) {
				targetIPsV6 = append(targetIPsV6, dnsData.AAAA...)
			}
		} else {
			targetIPsV4 = append(targetIPsV4, dnsData.A...)
		}
		if len(targetIPsV4) == 0 && len(targetIPsV6) == 0 {
			return targetIPsV4, targetIPsV6, fmt.Errorf("no IP addresses found for host: %s", target)
		}
	} else {
		targetIPsV4 = append(targetIPsV6, target)
		gologger.Debug().Msgf("Found %d addresses for %s\n", len(targetIPsV4), target)
	}

	return
}

func (r *Runner) getExcludeItems(s string) ([]string, error) {
	if isIpOrCidr(s) {
		return []string{s}, nil
	}

	ips4, ips6, err := r.host2ips(s)
	if err != nil {
		return nil, err
	}
	return append(ips4, ips6...), nil
}

func ipMatchesIpVersions(ip string, ipVersions ...string) bool {
	for _, ipVersion := range ipVersions {
		if ipVersion == IPv4 && IsIPv4(ip) {
			return true
		}
		if ipVersion == IPv6 && IsIPv6(ip) {
			return true
		}
	}
	return false
}
func ContainsAny(s string, ss ...string) bool {
	for _, sss := range ss {
		if strings.Contains(s, sss) {
			return true
		}
	}
	return false
}

func IsIPv6(ips ...interface{}) bool {
	for _, ip := range ips {
		switch ipv := ip.(type) {
		case string:
			parsedIP := net.ParseIP(ipv)
			isIP6 := parsedIP != nil && parsedIP.To16() != nil && ContainsAny(ipv, ":")
			if !isIP6 {
				return false
			}
		case net.IP:
			isIP6 := ipv != nil && ipv.To16() != nil && ContainsAny(ipv.String(), ":")
			if !isIP6 {
				return false
			}
		}
	}

	return true
}

func IsIPv4(ips ...interface{}) bool {
	for _, ip := range ips {
		switch ipv := ip.(type) {
		case string:
			parsedIP := net.ParseIP(ipv)
			isIP4 := parsedIP != nil && parsedIP.To4() != nil && strings.Contains(ipv, ".")
			if !isIP4 {
				return false
			}
		case net.IP:
			isIP4 := ipv != nil && ipv.To4() != nil && strings.Contains(ipv.String(), ".")
			if !isIP4 {
				return false
			}
		}
	}

	return true
}

func (s *Scanner) CdnCheck(ip string) (bool, string, error) {
	if s.cdn == nil {
		return false, "", errors.New("cdn client not initialized")
	}
	if !IsIP(ip) {
		return false, "", fmt.Errorf("%s is not a valid ip", ip)
	}

	// the goal is to check if ip is part of cdn/waf to decide if target should be scanned or not
	// since 'cloud' itemtype does not fit logic here , we consider target is not part of cdn/waf
	matched, value, itemType, err := s.cdn.Check(net.ParseIP((ip)))
	if itemType == "cloud" {
		return false, "", err
	}
	return matched, value, err
}

type Result struct {
	Host      string    `json:"host,omitempty" csv:"host"`
	IP        string    `json:"ip,omitempty" csv:"ip"`
	Port      int       `json:"port,omitempty" csv:"port"`
	Protocol  string    `json:"protocol,omitempty" csv:"protocol"`
	TLS       bool      `json:"tls,omitempty" csv:"tls"`
	IsCDNIP   bool      `json:"cdn,omitempty" csv:"cdn"`
	CDNName   string    `json:"cdn-name,omitempty" csv:"cdn-name"`
	TimeStamp time.Time `json:"timestamp,omitempty" csv:"timestamp"`
}

type jsonResult struct {
	Result
	PortNumber int    `json:"port"`
	Protocol   string `json:"protocol"`
	TLS        bool   `json:"tls"`
}

func (r *Result) JSON() ([]byte, error) {
	data := jsonResult{}
	data.TimeStamp = r.TimeStamp
	if r.Host != r.IP {
		data.Host = r.Host
	}
	data.IP = r.IP
	data.IsCDNIP = r.IsCDNIP
	data.CDNName = r.CDNName
	data.PortNumber = r.Port
	data.Protocol = r.Protocol
	data.TLS = r.TLS

	return json.Marshal(data)
}

var (
	NumberOfCsvFieldsErr = errors.New("exported fields don't match csv tags")
	headers              = []string{}
)

func (r *Result) CSVHeaders() ([]string, error) {
	ty := reflect.TypeOf(*r)
	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		csvTag := field.Tag.Get("csv")
		if !slices.Contains(headers, csvTag) {
			headers = append(headers, csvTag)
		}
	}
	return headers, nil
}

func (r *Result) CSVFields() ([]string, error) {
	var fields []string
	vl := reflect.ValueOf(*r)
	ty := reflect.TypeOf(*r)
	for i := 0; i < vl.NumField(); i++ {
		field := vl.Field(i)
		csvTag := ty.Field(i).Tag.Get("csv")
		fieldValue := field.Interface()
		if slices.Contains(headers, csvTag) {
			fields = append(fields, fmt.Sprint(fieldValue))
		}
	}
	return fields, nil
}

// WriteHostOutput writes the output list of host ports to an io.Writer
func WriteHostOutput(host string, ports []*Port, outputCDN bool, cdnName string, writer io.Writer) error {
	bufwriter := bufio.NewWriter(writer)
	sb := &strings.Builder{}

	for _, p := range ports {
		sb.WriteString(host)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(p.Port))
		if outputCDN && cdnName != "" {
			sb.WriteString(" [" + cdnName + "]")
		}
		sb.WriteString("\n")
		_, err := bufwriter.WriteString(sb.String())
		if err != nil {
			bufwriter.Flush()
			return err
		}
		sb.Reset()
	}
	return bufwriter.Flush()
}

// WriteJSONOutput writes the output list of subdomain in JSON to an io.Writer
func WriteJSONOutput(host, ip string, ports []*Port, outputCDN bool, isCdn bool, cdnName string, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	data := jsonResult{}
	data.TimeStamp = time.Now().UTC()
	if host != ip {
		data.Host = host
	}
	data.IP = ip
	if outputCDN {
		data.IsCDNIP = isCdn
		data.CDNName = cdnName
	}
	for _, p := range ports {
		data.PortNumber = p.Port
		data.Protocol = p.Protocol.String()
		data.TLS = p.TLS
		if err := encoder.Encode(&data); err != nil {
			return err
		}
	}
	return nil
}

// WriteCsvOutput writes the output list of subdomain in csv format to an io.Writer
func WriteCsvOutput(host, ip string, ports []*Port, outputCDN bool, isCdn bool, cdnName string, header bool, writer io.Writer) error {
	encoder := csv.NewWriter(writer)
	data := &Result{IP: ip, TimeStamp: time.Now().UTC(), Port: 0, Protocol: TCP.String(), TLS: false}
	if host != ip {
		data.Host = host
	}
	if outputCDN {
		data.IsCDNIP = isCdn
		data.CDNName = cdnName
	}
	if header {
		writeCSVHeaders(data, encoder)
	}

	for _, p := range ports {
		data.Port = p.Port
		data.Protocol = p.Protocol.String()
		data.TLS = p.TLS
		writeCSVRow(data, encoder)
	}
	encoder.Flush()
	return nil
}

func writeCSVHeaders(data *Result, writer *csv.Writer) {
	headers, err := data.CSVHeaders()
	if err != nil {
		gologger.Error().Msgf(err.Error())
		return
	}

	if err := writer.Write(headers); err != nil {
		errMsg := errors.Wrap(err, "Could not write headers")
		gologger.Error().Msgf(errMsg.Error())
	}
}

func writeCSVRow(data *Result, writer *csv.Writer) {
	rowData, err := data.CSVFields()
	if err != nil {
		gologger.Error().Msgf(err.Error())
		return
	}
	if err := writer.Write(rowData); err != nil {
		errMsg := errors.Wrap(err, "Could not write row")
		gologger.Error().Msgf(errMsg.Error())
	}
}

func (r *Runner) onReceive(hostResult *HostResult) {
	if !ipMatchesIpVersions(hostResult.IP, r.options.IPVersion...) {
		return
	}

	dt, err := r.scanner.IPRanger.GetHostsByIP(hostResult.IP)
	if err != nil {
		return
	}

	// receive event has only one port
	for _, p := range hostResult.Ports {
		ipPort := net.JoinHostPort(hostResult.IP, fmt.Sprint(p.Port))
		if r.unique.Has(ipPort) {
			return
		}
	}

	// recover hostnames from ip:port combination
	for _, p := range hostResult.Ports {
		ipPort := net.JoinHostPort(hostResult.IP, fmt.Sprint(p.Port))
		if dtOthers, ok := r.scanner.IPRanger.Hosts.Get(ipPort); ok {
			if otherName, _, err := net.SplitHostPort(string(dtOthers)); err == nil {
				// replace bare ip:port with host
				for idx, ipCandidate := range dt {
					if IsIP(ipCandidate) {
						dt[idx] = otherName
					}
				}
			}
		}
		_ = r.unique.Set(ipPort, struct{}{})
	}

	csvHeaderEnabled := true

	buffer := bytes.Buffer{}
	writer := csv.NewWriter(&buffer)
	for _, host := range dt {
		buffer.Reset()
		if host == "ip" {
			host = hostResult.IP
		}

		isCDNIP, cdnName, _ := r.scanner.CdnCheck(hostResult.IP)
		// console output
		if r.options.JSON || r.options.CSV {
			data := &Result{IP: hostResult.IP, TimeStamp: time.Now().UTC()}
			if r.options.OutputCDN {
				data.IsCDNIP = isCDNIP
				data.CDNName = cdnName
			}
			if host != hostResult.IP {
				data.Host = host
			}
			for _, p := range hostResult.Ports {
				data.Port = p.Port
				data.Protocol = p.Protocol.String()
				data.TLS = p.TLS
				if r.options.JSON {
					b, err := data.JSON()
					if err != nil {
						continue
					}
					buffer.Write([]byte(fmt.Sprintf("%s\n", b)))
				} else if r.options.CSV {
					if csvHeaderEnabled {
						writeCSVHeaders(data, writer)
						csvHeaderEnabled = false
					}
					writeCSVRow(data, writer)
				}
			}
		}
		if r.options.JSON {
			gologger.Silent().Msgf("%s", buffer.String())
		} else if r.options.CSV {
			writer.Flush()
			gologger.Silent().Msgf("%s", buffer.String())
		} else {
			for _, p := range hostResult.Ports {
				if r.options.OutputCDN && isCDNIP {
					gologger.Silent().Msgf("%s:%d [%s]\n", host, p.Port, cdnName)
				} else {
					gologger.Silent().Msgf("%s:%d\n", host, p.Port)
				}
			}
		}
	}
}

func (r *Runner) Close() {
	_ = os.RemoveAll(r.targetsFile)
	_ = r.scanner.IPRanger.Hosts.Close()
	if r.options.EnableProgressBar {
		_ = r.stats.Stop()
	}
	if r.scanner != nil {
		r.scanner.Close()
	}
}

func isIpOrCidr(s string) bool {
	return IsIP(s) || IsCIDR(s)
}

// Contains if a slice contains an element
func Contains[T comparable](inputSlice []T, element T) bool {
	for _, inputValue := range inputSlice {
		if inputValue == element {
			return true
		}
	}

	return false
}

func main() {
	options := ParseOptions()

	ports, err := ParsePorts(options)
	if err != nil {
		fmt.Errorf("could not parse ports: %s", err)
		return
	}

	options.configureHostDiscovery(ports)

	// default to ipv4 if no ipversion was specified
	if len(options.IPVersion) == 0 {
		options.IPVersion = []string{IPv4}
	}

	if options.Retries == 0 {
		options.Retries = DefaultRetriesSynScan
	}
	runner := &Runner{
		options: options,
	}

	dnsOptions := dnsx.DefaultOptions
	dnsOptions.MaxRetries = runner.options.Retries
	dnsOptions.Hostsfile = true
	if Contains(options.IPVersion, "6") {
		dnsOptions.QuestionTypes = append(dnsOptions.QuestionTypes, dns.TypeAAAA)
	}
	if len(runner.options.baseResolvers) > 0 {
		dnsOptions.BaseResolvers = runner.options.baseResolvers
	}
	dnsclient, err := dnsx.New(dnsOptions)
	if err != nil {
		fmt.Errorf("%s", err.Error())
		return
	}
	runner.dnsclient = dnsclient

	excludedIps, err := runner.parseExcludedIps(options)
	if err != nil {
		fmt.Errorf("%s", err.Error())
		return
	}

	runner.streamChannel = make(chan Target)

	uniqueCache := gcache.New[string, struct{}](1500).Build()
	runner.unique = uniqueCache

	scanOpts := &ScanOptions{
		Timeout:       time.Duration(options.Timeout) * time.Millisecond,
		Retries:       options.Retries,
		Rate:          options.Rate,
		PortThreshold: options.PortThreshold,
		ExcludeCdn:    options.ExcludeCDN,
		OutputCdn:     options.OutputCDN,
		ExcludedIps:   excludedIps,
		Proxy:         options.Proxy,
		ProxyAuth:     options.ProxyAuth,
		Stream:        options.Stream,
		OnReceive:     options.OnReceive,
		ScanType:      options.ScanType,
	}

	if scanOpts.OnReceive == nil {
		scanOpts.OnReceive = runner.onReceive
	}

	scanner, err := NewScanner(scanOpts)
	if err != nil {

	}

	runner.scanner = scanner

	runner.scanner.Ports = ports

	// Setup graceful exits
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			runner.ShowScanResultOnExit()
			gologger.Info().Msgf("CTRL+C pressed: Exiting\n")
			runner.Close()
			os.Exit(1)
		}
	}()

	if isPrivileged() {
		fmt.Println("root mode")
	} else {
		fmt.Println("please root")
		return
	}

}
