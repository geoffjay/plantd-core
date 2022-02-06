package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"syscall"

	"github.com/geoffjay/plantd-core/core"

	log "github.com/sirupsen/logrus"
)

func main() {
	processArgs()
	initLogging()

	port, _ := strconv.Atoi(core.Getenv("PORT", "5000"))
	bind := core.Getenv("ADDRESS", "0.0.0.0")
	app := NewService(port, bind)

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

	log.Debug("proxy exiting")
}

func initLogging() {
	level := core.Getenv("LOG_LEVEL", "info")
	if logLevel, err := log.ParseLevel(level); err == nil {
		log.SetLevel(logLevel)
	}
}

func processArgs() {
	if len(os.Args) > 1 {
		r := regexp.MustCompile("^-v$|(-{2})?version$")
		if r.Match([]byte(os.Args[1])) {
			fmt.Println(core.VERSION)
		}
		os.Exit(0)
	}
}
