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
)

func main() {
	processArgs()
	initLogging()

	app := NewService()

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go app.run(ctx, wg)

	log.Debug("service started")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	log.Debug("service terminated")

	cancelFunc()
	wg.Wait()

	log.Debug("logger exiting")
}

func initLogging() {
	level := util.Getenv("PLANTD_LOGGER_LOG_LEVEL", "info")
	if logLevel, err := log.ParseLevel(level); err == nil {
		log.SetLevel(logLevel)
	}
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
