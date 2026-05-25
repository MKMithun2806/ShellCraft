//go:build windows

package listener

// setupWinchHandler is a no-op on Windows — there is no SIGWINCH equivalent.
func setupWinchHandler() func() {
	return func() {} // no-op
}
