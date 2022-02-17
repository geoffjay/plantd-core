package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/geoffjay/plantd-core/core"

	log "github.com/sirupsen/logrus"
)

func main() {
	processArgs()
	initLogging()

	store := NewStore()
	path := core.Getenv("PLANTD_STATE_DB", "plantd-state.db")
	if err := store.Load(path); err != nil {
		log.Error(err)
	}
	defer store.Unload()

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go run(ctx, wg)

	log.WithFields(log.Fields{"scope": "main"}).Debug("starting")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	log.WithFields(log.Fields{"scope": "main"}).Debug("terminated")

	cancelFunc()
	wg.Wait()

	log.WithFields(log.Fields{"scope": "main"}).Debug("exiting")
}

func run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.WithFields(log.Fields{"scope": "run"}).Debug("starting")

	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.WithFields(log.Fields{"scope": "run"}).Debug("processing")
		}
	}()

	<-ctx.Done()
	log.WithFields(log.Fields{"scope": "run"}).Debug("exiting")
}

func initLogging() {
	level := core.Getenv("PLANTD_STATE_LOG_LEVEL", "info")
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
