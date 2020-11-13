package signal

import (
	"os"
	"os/signal"
)

var onlyOneSignalHandler = make(chan struct{})
var shutdownHandler chan os.Signal

func SetupSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler)

	stop := make(chan struct{})
	shutdownHandler = make(chan os.Signal, 2)
	signal.Notify(shutdownHandler, shutdownSignals...)
	go func() {
		<-shutdownHandler
		close(stop)
		<-shutdownHandler
		os.Exit(1)
	}()

	return stop
}
