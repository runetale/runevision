package safesocket

import (
	"context"
	"errors"
	"net"
	"os"
	"runtime"
	"time"
)

type closeable interface {
	CloseRead() error
	CloseWrite() error
}

func ConnCloseRead(c net.Conn) error {
	return c.(closeable).CloseRead()
}

func ConnCloseWrite(c net.Conn) error {
	return c.(closeable).CloseWrite()
}

func ConnectContext(ctx context.Context, path string) (net.Conn, error) {
	for {
		c, err := connect(ctx, path)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			time.Sleep(250 * time.Millisecond)
			continue
		}
		return c, err
	}
}

func Connect(path string) (net.Conn, error) {
	return ConnectContext(context.Background(), path)
}

func Listen(path string) (net.Listener, error) {
	return listen(path)
}

var (
	ErrTokenNotFound = errors.New("no token found")
	ErrNoTokenOnOS   = errors.New("no token on " + runtime.GOOS)
)

func PlatformUsesPeerCreds() bool { return GOOSUsesPeerCreds(runtime.GOOS) }

func GOOSUsesPeerCreds(goos string) bool {
	switch goos {
	case "linux", "darwin", "freebsd":
		return true
	}
	return false
}

func VisionSocketPath() string {
	if runtime.GOOS == "windows" {
		return `\\.\pipe\ProtectedPrefix\Administrators\Vision\visiond`
	}
	if runtime.GOOS == "darwin" {
		return "/var/run/visiond.socket"
	}
	if fi, err := os.Stat("/var/run"); err == nil && fi.IsDir() {
		return "/var/run/thor/visiond.sock"
	}
	return "visiond.sock"
}
