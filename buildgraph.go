package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/services"
	"go.uber.org/zap"
)

func main() {
	var path string
	flag.StringVar(&path, "c", "", "path to config")
	flag.Parse()

	cfg, err := config.NewConfig(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	run(logger, cfg)
}

func run(logger *zap.Logger, cfg *config.Config) {
	var (
		stopChan = make(chan os.Signal)
		errChan  = make(chan error)
	)

	services, err := services.NewCollection(logger, cfg, errChan)
	if err != nil {
		logger.Fatal("Service initialization error", zap.Error(err))
	}

	signal.Notify(stopChan)

	select {
	case <-stopChan:
		logger.Info("Interruption signal has been received")
	case err := <-errChan:
		logger.Error("Fatal error", zap.Error(err))
	}

	services.Stop()
}
