// このパッケージはapi serverがLocalBackendに接続するための構造体です
package localclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/runetale/runevision/domain/entity"
	"github.com/runetale/runevision/domain/requests"
	"github.com/runetale/runevision/safesocket"
	"github.com/runetale/runevision/utility"
)

const LocalBackendAPIHost = "visiond.sock"
const localBackendFullEndpoint = "http://" + LocalBackendAPIHost + "/localapi" + "/v0"

const (
	scanAPI = "/scan"
)

type LocalClient struct {
	endpoint     string
	vsClient     *http.Client
	tsClientOnce sync.Once

	logger *utility.Logger
}

func NewLocalClient(logger *utility.Logger) *LocalClient {
	return &LocalClient{
		endpoint: localBackendFullEndpoint,
		logger:   logger,
	}
}

func (lc *LocalClient) GetRequestEndPoint() string {
	return lc.endpoint
}

func (lc *LocalClient) socket() string {
	return safesocket.VisionSocketPath()
}

func (lc *LocalClient) dialer() func(ctx context.Context, network, addr string) (net.Conn, error) {
	return lc.defaultDialer
}

func (lc *LocalClient) defaultDialer(ctx context.Context, network, addr string) (net.Conn, error) {
	if addr != "visiond.sock:80" {
		return nil, fmt.Errorf("unexpected URL address %q", addr)
	}
	return safesocket.Connect(lc.socket())
}

// DoLocalRequestはapi serverがvision daemonと会話するときに必要
func (lc *LocalClient) doLocalRequest(req *http.Request) (*http.Response, error) {
	lc.tsClientOnce.Do(func() {
		lc.vsClient = &http.Client{
			Transport: &http.Transport{
				DialContext: lc.dialer(),
			},
		}
	})
	return lc.vsClient.Do(req)
}

func (lc *LocalClient) DoScan(seqID string, request *requests.HackDoScanRequest, c echo.Context) (*entity.HackHistory, error) {
	apiURL := lc.GetRequestEndPoint() + scanAPI
	localAPIURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	body := io.NopCloser(bytes.NewBuffer(jsonBytes))

	c.SetRequest(&http.Request{
		Method: c.Request().Method,
		URL:    localAPIURL,
		Body:   body,
		Header: http.Header{
			"SequentialID": []string{seqID},
		},
	})

	res, err := lc.doLocalRequest(c.Request())
	if err != nil {
		return nil, err
	}

	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var hh entity.HackHistory
	err = json.Unmarshal(resbody, &hh)
	if err != nil {
		return nil, err
	}

	return &hh, nil
}
