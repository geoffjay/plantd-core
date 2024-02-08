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
	plog "github.com/geoffjay/plantd/core/log"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := GetConfig()

	processArgs()
	plog.Initialize(config.Log)
	// initLogging(config.Log)

	service := NewService()
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

// func initLogging(logConfig cfg.LogConfig) {
// 	if logLevel, err := log.ParseLevel(logConfig.Level); err == nil {
// 		log.SetLevel(logLevel)
// 	}
//
// 	if logConfig.Formatter == "json" {
// 		log.SetFormatter(&log.JSONFormatter{
// 			TimestampFormat: "2006-01-02 15:04:05",
// 		})
// 	} else {
// 		log.SetFormatter(&log.TextFormatter{
// 			FullTimestamp:   true,
// 			TimestampFormat: "2006-01-02 15:04:05",
// 		})
// 	}
//
// 	opts := loki.NewLokiHookOptions().WithLevelMap(
// 		loki.LevelMap{log.PanicLevel: "critical"},
// 	).WithFormatter(
// 		&log.JSONFormatter{},
// 	).WithStaticLabels(
// 		logConfig.Loki.Labels,
// 	)
//
// 	hook := loki.NewLokiHookWithOpts(
// 		logConfig.Loki.Address,
// 		opts,
// 		log.InfoLevel,
// 		log.WarnLevel,
// 		log.ErrorLevel,
// 		log.FatalLevel,
// 	)
//
// 	log.AddHook(hook)
// }

func processArgs() {
	if len(os.Args) > 1 {
		r := regexp.MustCompile("^-V$|(-{2})?version$")
		if r.Match([]byte(os.Args[1])) {
			fmt.Println(core.VERSION)
		}
		os.Exit(0)
	}
}
