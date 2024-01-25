package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/geoffjay/plantd/core"
	"github.com/geoffjay/plantd/core/bus"
	"github.com/geoffjay/plantd/core/http"
	"github.com/geoffjay/plantd/core/mdp"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Service defines the service type.
type Service struct {
	buses    []*bus.Bus
	config   *brokerConfig
	handler  *Handler
	endpoint string
	broker   *mdp.Broker
	running  bool
	worker   *mdp.Worker
}

// NewService creates an instance of the service.
func NewService(config *brokerConfig) *Service {
	service := &Service{
		buses:    initBuses(config),
		config:   config,
		handler:  NewHandler(),
		endpoint: config.Endpoint,
		broker:   nil,
		running:  false,
		worker:   nil,
	}

	service.dumpConfig()
	service.RegisterCallback("service", &serviceCallback{name: "service"})
	service.RegisterCallback("services", &servicesCallback{name: "services"})

	if err := service.initBroker(); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("failed to initialize broker")
		return nil
	}

	if err := service.initWorker(config); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("failed to initialize worker")
		return nil
	}

	return service
}

func (s *Service) dumpConfig() {
	json, err := core.MarshalConfig(s.config)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to marshal config")
	}

	log.WithFields(log.Fields{"context": "service"}).Trace(json)
}

func initBuses(config *brokerConfig) (buses []*bus.Bus) {
	for _, b := range config.Buses {
		log.WithFields(log.Fields{
			"bus":      b.Name,
			"backend":  b.Backend,
			"frontend": b.Frontend,
			"capture":  b.Capture,
		}).Info("initializing message bus")
		buses = append(buses, bus.NewBus(b.Name, b.Name, b.Backend, b.Frontend, b.Capture))
	}

	return
}

func (s *Service) initBroker() error {
	s.broker, _ = mdp.NewBroker(s.endpoint)
	if err := s.broker.Bind(); err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"endpoint": s.endpoint,
		}).Error("failed to bind to endpoint")
		return err
	}

	return nil
}

func (s *Service) initWorker(config *brokerConfig) error {
	var err error
	if s.worker, err = mdp.NewWorker(config.ClientEndpoint, "org.plantd.Broker"); err != nil {
		log.WithFields(log.Fields{
			"err":             err,
			"client-endpoint": config.ClientEndpoint,
		}).Error("failed to setup message queue worker")
		return err
	}
	return nil
}

// Run handles the service execution.
func (s *Service) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer s.worker.Close()
	defer wg.Done()

	log.WithFields(log.Fields{"context": "service.run"}).Debug("starting")

	wg.Add(len(s.buses))
	for _, item := range s.buses {
		elem := item
		go func() {
			if err := elem.Start(ctx, wg); err != nil {
				log.Error(err)
			}
		}()
	}

	wg.Add(3)
	go s.runWorker(ctx, wg)
	go s.runBroker(ctx)
	go s.runApp(ctx, wg)

	<-ctx.Done()

	s.broker.Close()

	log.WithFields(log.Fields{"context": "service.run"}).Debug("exiting")
}

func (s *Service) runBroker(ctx context.Context) {
	done := make(chan bool, 1)
	go s.broker.Run(done)

	SetLastError(errors.New("none"))
	log.Debug("starting broker")

	for {
		var err error
		var event mdp.Event
		log.Debug("waiting for message")
		select {
		case event = <-s.broker.EventChannel:
			log.Debug(event)
		case err = <-s.broker.ErrorChannel:
			SetLastError(err)
			log.WithFields(log.Fields{"error": err}).Error("received error from message queue")
		case <-ctx.Done():
			_ = s.broker.Close()
			done <- true
			log.Debug("exiting broker")
			return
		}
	}
}

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

// App creates a web application to serve static website content.
func (s *Service) runApp(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	staticContents := "./public"
	bindAddress := "0.0.0.0"
	bindPort := 4999

	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()

		r.Use(static.Serve("/", static.LocalFile(staticContents, false)))
		r.Use(gin.Recovery())
		r.Use(http.LoggerMiddleware())
		r.Use(brokerMiddleware(s))

		initializeRoutes(r)

		if err := r.Run(fmt.Sprintf("%s:%d", bindAddress, bindPort)); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	log.Debug("exiting web application")
}

// RegisterCallback is a pointless wrapper around the handler.
func (s *Service) RegisterCallback(name string, callback HandlerCallback) {
	s.handler.AddCallback(name, callback)
}
