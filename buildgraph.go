package main

import (
	"os"
	"os/signal"

	"github.com/vitalyisaev2/buildgraph/common"
	"github.com/vitalyisaev2/buildgraph/config"
)

func main() {
}

func launch(logger common.Logger, cfg *config.Config) error {
	var (
		err      error
		stopChan = make(chan os.Signal)
		errChan  = make(chan error)
	)

	signal.Notify(stopChan)

	select {
	case <-stopChan:
		logger.Warn("Interruption signal has been received")
	case err = <-errChan:
	}

	return err
}
