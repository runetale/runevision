package backend

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/runetale/runevision/domain/entity"
	"github.com/runetale/runevision/domain/requests"
	"github.com/runetale/runevision/localclient"
	"github.com/runetale/runevision/types"
)

type backendAPIHandler func(*Handler, http.ResponseWriter, *http.Request)

// local backend api pathの種類
// わかりやすい様に一旦ここに列挙
var backendApiPath = map[string]backendAPIHandler{
	"ping": (*Handler).ping,
	"scan": (*Handler).scan,
}

type Handler struct {
	lb *LocalBackend
}

func NewHandler(b *LocalBackend) *Handler {
	return &Handler{lb: b}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.lb == nil {
		http.Error(w, "server has no local backend", http.StatusInternalServerError)
		return
	}
	if !h.validHost(r.Host) {
		http.Error(w, "invalid localapi request", http.StatusForbidden)
		return
	}

	if fn, ok := handlerForPath(r.URL.Path); ok {
		fn(h, w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (h *Handler) validHost(hostname string) bool {
	switch hostname {
	case localclient.LocalBackendAPIHost:
		return true
	}
	host, _, err := net.SplitHostPort(hostname)
	if err != nil {
		return false
	}
	if host == "localhost" {
		return true
	}
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return false
	}
	return addr.IsLoopback()
}

func (*Handler) serveLocalAPIRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "vision daemon\n")
}

func writeErrorJSON(w http.ResponseWriter, err error) {
	if err == nil {
		err = errors.New("unexpected nil error")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	type E struct {
		Error string `json:"error"`
	}
	json.NewEncoder(w).Encode(E{err.Error()})
}

func handlerForPath(urlPath string) (h backendAPIHandler, ok bool) {
	if urlPath == "/" {
		return (*Handler).serveLocalAPIRoot, true
	}

	suff, ok := strings.CutPrefix(urlPath, "/localapi/v0/")
	if !ok {
		return nil, false
	}

	if fn, ok := backendApiPath[suff]; ok {
		return fn, true
	}

	return nil, false
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "required POST method", http.StatusBadRequest)
		return
	}

	ping, err := h.lb.Ping()
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	type res struct {
		Ping string `json:"ping"`
	}

	json.NewEncoder(w).Encode(res{Ping: ping})
}

// hack/scan apiから呼ばれる。serverからのリクエストに加え、sequential idが追加された状態で
// リクエストが送られてくる
func (h *Handler) scan(w http.ResponseWriter, r *http.Request) {
	var req requests.HackDoScanRequest
	sid := r.Header.Get("SequentialID")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	err = json.Unmarshal(b, &req)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	err = h.lb.Scan(types.SequenceID(sid), &req)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	status := h.lb.GetStatus(types.SequenceID(sid))

	hh := entity.NewHackHistory(req.Name, sid, string(status))

	json.NewEncoder(w).Encode(hh)
}
