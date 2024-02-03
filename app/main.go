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

// @title Plantd Web Application
// @version 1.0
// @description The Plantd web application.
// @contact.name Geoff Johnson
// @contact.email geoff.jay@gmail.com
// @license.name MIT
// @license.url https://opensource.org/license/mit/
func main() {
	processArgs()
	initLogging()

	service := service{}
	service.init()

	fields := log.Fields{"service": "app", "context": "main"}

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go service.run(ctx, wg)

	log.WithFields(fields).Debug("starting")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	log.WithFields(fields).Debug("terminated")

	cancelFunc()
	wg.Wait()

	log.WithFields(fields).Debug("exiting")
}

func initLogging() {
	level := util.Getenv("PLANTD_APP_LOG_LEVEL", "info")
	if logLevel, err := log.ParseLevel(level); err == nil {
		log.SetLevel(logLevel)
	}

	format := util.Getenv("PLANTD_APP_LOG_FORMAT", "text")
	if format == "json" {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	opts := loki.NewLokiHookOptions().WithLevelMap(
		loki.LevelMap{log.PanicLevel: "critical"},
	).WithFormatter(
		&log.JSONFormatter{},
	).WithStaticLabels(
		loki.Labels{
			"app":         "app",
			"environment": "development",
		},
	)

	lokiAddress := util.Getenv("PLANTD_APP_LOG_LOKI_ADDRESS", "http://localhost:3100")
	hook := loki.NewLokiHookWithOpts(
		lokiAddress,
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
