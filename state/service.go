package main

import (
	"context"
	"sync"
	"time"

	"github.com/geoffjay/plantd/core/bus"
	"github.com/geoffjay/plantd/core/mdp"
	"github.com/geoffjay/plantd/core/util"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	store  *Store
	sink   *bus.Sink
	worker *mdp.Worker
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) setupStore() {
	s.store = NewStore()
	path := util.Getenv("PLANTD_STATE_DB", "plantd-state.db")
	if err := s.store.Load(path); err != nil {
		log.WithFields(log.Fields{"err": err}).Panic("failed to setup KV store")
	}
	defer s.store.Unload()
}

func (s *Service) setupWorker() {
	var err error
	if s.worker, err = mdp.NewWorker("tcp://127.0.0.1:7200", "org.plantd.State"); err != nil {
		log.WithFields(log.Fields{"err": err}).Panic("failed to setup message queue worker")
	}
	defer s.worker.Close()
}

func (s *Service) setupSink() {
	s.sink = bus.NewSink()
}

func (s *Service) Run(ctx context.Context, wg *sync.WaitGroup) {
	s.setupStore()
	s.setupSink()
	s.setupWorker()

	defer wg.Done()
	log.WithFields(log.Fields{"context": "run"}).Debug("starting")

	wg.Add(2)
	go s.runSink(ctx, wg)
	go s.runWorker(ctx, wg)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.WithFields(log.Fields{"context": "run"}).Debug("processing")
		}
	}()

	<-ctx.Done()

	log.WithFields(log.Fields{"context": "run"}).Debug("exiting")
}

func (s *Service) runSink(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		for {
			log.WithFields(log.Fields{"context": "sink"}).Debug("listening")
			time.Sleep(10 * time.Second)
		}
	}()

	<-ctx.Done()

	log.WithFields(log.Fields{"context": "worker"}).Debug("exiting")
}

func (s *Service) runWorker(ctx context.Context, wg *sync.WaitGroup) {
	var err error

	defer wg.Done()

	go func() {
		var request, reply []string
		for {
			if s.worker.Terminated() {
				break
			}

			log.WithFields(log.Fields{"context": "worker"}).Debug("waiting for request")
			if request, err = s.worker.Recv(reply); err != nil {
				log.WithFields(log.Fields{"error": err}).Error("failed while receiving request")
			}
			log.WithFields(log.Fields{"context": "worker", "request": request}).Debug("received request")

			if len(request) == 0 {
				log.WithFields(log.Fields{"context": "worker"}).Debug("received request is empty")
				continue
			}

			// msgType := request[0]

			// Reset reply
			reply = []string{}

			for _, part := range request[1:] {
				log.WithFields(log.Fields{"context": "worker", "part": part}).Debug("processing message")
				// Reply with an empty response for now
				reply = append(reply, "{}")
			}
		}
	}()

	<-ctx.Done()
	s.worker.Shutdown()

	log.WithFields(log.Fields{"context": "worker"}).Debug("exiting")
}
