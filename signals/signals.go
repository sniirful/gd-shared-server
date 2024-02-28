package signals

import (
	"os"
	"os/signal"
)

func CaptureSIGINT(handlerFunction func()) func() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		<-sigCh
		handlerFunction()
	}()
	return func() {
		close(sigCh)
		signal.Stop(sigCh)
	}
}
