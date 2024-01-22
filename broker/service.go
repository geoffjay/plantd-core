package main

import (
	"context"
	"sync"

	"github.com/geoffjay/plantd/core"
	"github.com/geoffjay/plantd/core/mdp"

	log "github.com/sirupsen/logrus"
)

// Service defines the service type.
type Service struct {
	config  *brokerConfig
	handler *Handler
	// nolint: unused
	broker *mdp.Broker
	worker *mdp.Worker
}

// NewService creates an instance of the service.
func NewService(config *brokerConfig) *Service {
	json, err := core.MarshalConfig(config)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to marshal config")
	}

	log.WithFields(
		log.Fields{
			"context": "service",
			"config":  json,
		},
	).Trace("creating service")

	return &Service{
		config: config,
	}
}

// nolint: unused
func (s *Service) setupBroker() {
	//  s.broker := service.NewBroker(config)
	// if broker == nil {
	// 	panic("an unrecoverable error occurred while creating the broker")
	// }
}

func (s *Service) setupHandler() {
	s.handler = NewHandler()
	s.RegisterCallback("service", &serviceCallback{name: "service"})
	s.RegisterCallback("services", &servicesCallback{name: "services"})
}

func (s *Service) setupWorker() {
	var err error
	if s.worker, err = mdp.NewWorker("tcp://127.0.0.1:7200", "org.plantd.Broker"); err != nil {
		log.WithFields(log.Fields{"err": err}).Panic("failed to setup message queue worker")
	}
}

// Run handles the service execution.
func (s *Service) Run(ctx context.Context, wg *sync.WaitGroup) {
	s.setupHandler()
	s.setupWorker()

	defer s.worker.Close()

	defer wg.Done()
	log.WithFields(log.Fields{"context": "service.run"}).Debug("starting")

	wg.Add(1)
	go s.runWorker(ctx, wg)

	<-ctx.Done()

	log.WithFields(log.Fields{"context": "service.run"}).Debug("exiting")
}

// func (s *Service) runBroker(ctx context.Context, wg *sync.WaitGroup) {
//   var err error
//
//   defer wg.Done()
//
//   wg.Add(3)
// 	go broker.App(ctx, wg)
// 	go broker.Serve(ctx, wg)
// 	go broker.Start(ctx, wg)
//
//   <-ctx.Done()
//   s.broker.Shutdown()
//
//   log.WithFields(log.Fields{"context": "broker"}).Debug("exiting")
// }

// nolint: funlen
func (s *Service) runWorker(ctx context.Context, wg *sync.WaitGroup) {
	var err error

	defer wg.Done()

	go func() {
		var request, reply []string
		for !s.worker.Terminated() {
			log.WithFields(log.Fields{"context": "service.worker"}).Debug("waiting for request")

			if request, err = s.worker.Recv(reply); err != nil {
				log.WithFields(log.Fields{"error": err}).Error("failed while receiving request")
			}

			log.WithFields(
				log.Fields{
					"context": "service.worker",
					"request": request,
				},
			).Debug("received request")

			if len(request) == 0 {
				log.WithFields(
					log.Fields{"context": "service.worker"},
				).Debug("received request is empty")
				continue
			}

			msgType := request[0]

			// Reset reply
			reply = []string{}

			for _, part := range request[1:] {
				log.WithFields(log.Fields{
					"context": "worker",
					"part":    part,
				}).Debug("processing message")
				var data []byte
				switch msgType {
				case "services", "service":
					log.Tracef("part: %s", part)
					if data, err = s.handler.callbacks[msgType].Execute(part); err != nil {
						log.WithFields(log.Fields{
							"context": "worker",
							"type":    msgType,
							"error":   err,
						}).Warn("message failed")
						break
					}
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

// RegisterCallback is a pointless wrapper around the handler.
func (s *Service) RegisterCallback(name string, callback HandlerCallback) {
	s.handler.AddCallback(name, callback)
}
