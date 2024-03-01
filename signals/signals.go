package signals

import (
	"os"
	"os/signal"
)

func CaptureInterrupt(handlerFunction func()) func() {
	return captureSignal(handlerFunction, os.Interrupt, false)
}

func captureSignal(handlerFunction func(), signalToSend os.Signal, sendInitialSignal bool) func() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, signalToSend)

	go func() {
		for range sigCh {
			handlerFunction()
		}
	}()
	if sendInitialSignal {
		sigCh <- signalToSend
	}
	return func() {
		close(sigCh)
		signal.Stop(sigCh)
	}
}
