//go:build !windows

package signals

import "syscall"

func CaptureSIGWINCH(handlerFunction func()) func() {
	return captureSignal(handlerFunction, syscall.SIGWINCH, true)
}
