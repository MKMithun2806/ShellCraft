//go:build !windows

package listener

import (
	"os"
	"os/signal"
	"syscall"
)

// setupWinchHandler registers a SIGWINCH handler for terminal resize events.
// Returns a stop function to clean up the signal registration.
func setupWinchHandler() func() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			// Terminal resized — the kernel handles the actual resize;
			// this goroutine keeps the channel drained to avoid blocking.
			// Future enhancement: forward TIOCSWINSZ over the TCP connection.
		}
	}()
	return func() { signal.Stop(ch) }
}
