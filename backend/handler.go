package backend

import (
	"io"
	"net"
	"net/http"
	"net/netip"
	"strings"
)

type backendAPIHandler func(*Handler, http.ResponseWriter, *http.Request)

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
	if r.Referer() != "" || r.Header.Get("Origin") != "" || !h.validHost(r.Host) {
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
}

func (h *Handler) scan(w http.ResponseWriter, r *http.Request) {
}
