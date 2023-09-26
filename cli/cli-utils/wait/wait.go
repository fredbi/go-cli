package wait

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

// ErrTimeout tells a wait for an event has timed out
var ErrTimeout = errors.New("the waiting for mounts has been too long: bailed")

// FileIsPresent waits for a file to be present on the local file system.
//
// This does nothing if the pth argument is empty.
func FileIsPresent(pth string, opts ...Option) (bool, error) {
	if pth == "" {
		return true, nil
	}

	o := applyOptions(opts)
	ctx, cancel := context.WithTimeoutCause(context.Background(), o.timeout, ErrTimeout)
	defer cancel()

	ticker := time.NewTicker(o.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()

		case <-ticker.C:
			if _, err := os.Stat(pth); err == nil {
				return true, nil
			}
		}
	}
}

// PortIsOpen waits for a network port to be open.
//
// This does nothing if the hostport argument is empty.
func PortIsOpen(hostport string, opts ...Option) (bool, error) {
	if hostport == "" {
		return true, nil
	}

	o := applyOptions(opts)

	ctx, cancel := context.WithTimeoutCause(context.Background(), o.timeout, ErrTimeout)
	defer cancel()

	var d net.Dialer
	_, err := d.DialContext(ctx, o.dialProtocol, hostport)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Done creates a file on the local file system.
//
// A typical use is to signal to other containers that we are done (e.g. proxy container),
// using a touch file to communicate between the containers on the same pod.
//
// This does nothing if the touch file argument is empty.
func Done(touch string, _ ...Option) error {
	if touch == "" {
		return nil
	}

	if dir := filepath.Dir(touch); dir != "." {
		if err := os.MkdirAll(filepath.Dir(touch), 0755); err != nil {
			return fmt.Errorf("warn: could not create dir for signal file: %q: %v", touch, err)
		}
	}

	f, err := os.Create(touch)
	if err != nil {
		return fmt.Errorf("warn: could not create signal file: %q: %v", touch, err)
	}

	return f.Close()
}
