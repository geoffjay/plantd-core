package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/geoffjay/plantd/core/util"

	log "github.com/sirupsen/logrus"
)

func main() {
	initLogging()

	port, _ := strconv.Atoi(util.Getenv("PLANTD_MODULE_ECHO_PORT", "5001"))
	bind := util.Getenv("PLANTD_MODULE_ECHO_ADDRESS", "0.0.0.0")
	service := NewService(port, bind)

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go service.run(ctx, wg)

	log.WithFields(log.Fields{"module": "echo"}).Debug("started")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	log.WithFields(log.Fields{"module": "echo"}).Debug("terminated")

	cancelFunc()
	wg.Wait()

	log.WithFields(log.Fields{"module": "echo"}).Debug("exiting")
}

func initLogging() {
	level := util.Getenv("PLANTD_MODULE_ECHO_LOG_LEVEL", "info")
	if logLevel, err := log.ParseLevel(level); err == nil {
		log.SetLevel(logLevel)
	}
}
