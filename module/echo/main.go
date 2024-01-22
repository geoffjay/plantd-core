package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/geoffjay/plantd/core/mdp"
	"github.com/geoffjay/plantd/core/util"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	worker *mdp.Worker
}

func (s *Service) setup() {
	var err error

	endpoint := util.Getenv("PLANTD_BROKER_ENDPOINT", "tcp://127.0.0.1:9797")
	if s.worker, err = mdp.NewWorker(endpoint, "org.plantd.module.Echo"); err != nil {
		log.WithFields(
			log.Fields{"err": err},
		).Panic("failed to setup message queue worker")
	}
}

func (s *Service) run(ctx context.Context, wg *sync.WaitGroup) {
	s.setup()

	defer s.worker.Close()
	defer wg.Done()
	log.Debug("starting worker")

	wg.Add(1)
	go s.runWorker(ctx, wg)

	<-ctx.Done()
	log.Debug("exiting worker")
}

func (s *Service) runWorker(ctx context.Context, wg *sync.WaitGroup) {
	var err error

	defer wg.Done()

	go func() {
		var request, reply []string
		for !s.worker.Terminated() {
			log.Debug("waiting for request")

			if request, err = s.worker.Recv(reply); err != nil {
				log.Error("failed while receiving request")
			}

			log.Debug("received request")

			if len(request) == 0 {
				log.Debug("received request is empty")
				continue
			}

			msgType := request[0]

			// reset reply
			reply = []string{}

			for _, part := range request[1:] {
				log.WithFields(log.Fields{
					"context": "worker",
					"part":    part,
				}).Debug("processing message")
				var data []byte
				switch msgType {
				case "echo":
					log.Tracef("part: %s", part)
					// pong
					data = []byte(part)
					log.Tracef("data: %s", data)
				default:
					log.Error("invalid message type provided")
				}

				reply = append(reply, string(data))
			}

			log.Tracef("reply: %+v", reply)
		}
	}()

	<-ctx.Done()
	s.worker.Shutdown()

	log.WithFields(log.Fields{"context": "worker"}).Debug("exiting")
}

func main() {
	service := &Service{}

	log.SetLevel(log.DebugLevel)

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go service.run(ctx, wg)

	log.Debug("service started")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	log.Debug("service terminated")

	cancelFunc()
	wg.Wait()

	log.Debug("exiting")
}
