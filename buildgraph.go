package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/vitalyisaev2/buildgraph/config"
	"github.com/vitalyisaev2/buildgraph/service"
	"github.com/vitalyisaev2/buildgraph/webserver"
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

	logger := makeLogger()
	run(logger, cfg)
}

func makeLogger() *logrus.Logger {
	var log = logrus.New()
	//log.Formatter = new(logrus.JSONFormatter)
	log.Level = logrus.DebugLevel
	log.Out = os.Stdout
	return log
}

func run(logger *logrus.Logger, cfg *config.Config) {
	var (
		stopChan = make(chan os.Signal)
		errChan  = make(chan error)
	)

	logger.Info("starting subsystems")
	services, err := service.NewCollection(logger, cfg)
	if err != nil {
		logger.Fatal("service initialization error", zap.Error(err))
	}

	logger.Info("starting webserver")
	ws := webserver.NewWebServer(services, cfg.Webserver, errChan)

	signal.Notify(stopChan)

	select {
	case <-stopChan:
		logger.Info("interruption signal has been received")
	case err := <-errChan:
		logger.Error("fatal error", zap.Error(err))
	}

	ws.Stop()
	services.Stop()
}
