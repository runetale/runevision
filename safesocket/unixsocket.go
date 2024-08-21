//go:build !windows && !js && !plan9

package safesocket

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
)

func connect(ctx context.Context, path string) (net.Conn, error) {
	if runtime.GOOS == "js" {
		return nil, errors.New("safesocket.Connect not yet implemented on js/wasm")
	}
	var std net.Dialer
	return std.DialContext(ctx, "unix", path)
}

func listen(path string) (net.Listener, error) {
	c, err := net.Dial("unix", path)
	if err == nil {
		c.Close()
		return nil, fmt.Errorf("%v: address already in use", path)
	}
	_ = os.Remove(path)

	perm := socketPermissions()

	sockDir := filepath.Dir(path)
	if _, err := os.Stat(sockDir); os.IsNotExist(err) {
		os.MkdirAll(sockDir, 0755)
		if perm == 0666 {
			if fi, err := os.Stat(sockDir); err == nil && fi.Mode()&0077 == 0 {
				if err := os.Chmod(sockDir, 0755); err != nil {
					log.Print(err)
				}
			}
		}
	}
	pipe, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}
	os.Chmod(path, perm)
	return pipe, err
}

func socketPermissions() os.FileMode {
	if PlatformUsesPeerCreds() {
		return 0666
	}
	return 0600
}
