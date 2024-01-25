package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/geoffjay/plantd/core"
	"github.com/geoffjay/plantd/core/util"

	log "github.com/sirupsen/logrus"
	loki "github.com/yukitsune/lokirus"
)

var config brokerConfig

func main() {
	processArgs()
	initConfig()
	initLogging()

	service := NewService(&config)
	fields := log.Fields{"service": "broker", "context": "main"}

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go service.Run(ctx, wg)

	log.WithFields(fields).Debug("starting")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	log.WithFields(fields).Debug("terminated")

	cancelFunc()
	wg.Wait()

	log.WithFields(fields).Debug("exiting")
}

func initConfig() {
	if err := core.LoadConfig("broker", &config); err != nil {
		log.Fatalf("error reading config file: %s\n", err)
	}
}

func initLogging() {
	level := util.Getenv("PLANTD_BROKER_LOG_LEVEL", "info")
	if logLevel, err := log.ParseLevel(level); err == nil {
		log.SetLevel(logLevel)
	}

	format := util.Getenv("PLANTD_BROKER_LOG_FORMAT", "text")
	if format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	opts := loki.NewLokiHookOptions().WithLevelMap(
		loki.LevelMap{log.PanicLevel: "critical"},
	).WithFormatter(
		&log.JSONFormatter{},
	).WithStaticLabels(
		loki.Labels{
			"app":         "broker",
			"environment": "development",
		},
	)

	hook := loki.NewLokiHookWithOpts(
		"http://localhost:3100",
		opts,
		log.InfoLevel,
		log.WarnLevel,
		log.ErrorLevel,
		log.FatalLevel,
	)

	log.AddHook(hook)
}

func processArgs() {
	if len(os.Args) > 1 {
		r := regexp.MustCompile("^-V$|(-{2})?version$")
		if r.Match([]byte(os.Args[1])) {
			fmt.Println(core.VERSION)
		}
		os.Exit(0)
	}
}
