package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/geoffjay/plantd/core/http"
	"github.com/geoffjay/plantd/core/mdp"
	"github.com/geoffjay/plantd/core/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Service type for REST API.
type Service struct {
	port   int
	bind   string
	worker *mdp.Worker
}

// NewService constructs and instance of a service type.
func NewService(port int, bind string) *Service {
	return &Service{
		port:   port,
		bind:   bind,
		worker: nil,
	}
}

func (s *Service) setup() {
	var err error

	endpoint := util.Getenv("PLANTD_MODULE_ECHO_BROKER_ENDPOINT", "tcp://127.0.0.1:9797")
	if s.worker, err = mdp.NewWorker(endpoint, "org.plantd.module.Echo"); err != nil {
		log.WithFields(log.Fields{
			"module": "echo",
			"err":    err,
		}).Panic("failed to setup message queue worker")
	}
}

func (s *Service) run(ctx context.Context, wg *sync.WaitGroup) {
	s.setup()

	defer s.worker.Close()
	defer wg.Done()

	wg.Add(1)
	go s.runWorker(ctx, wg)

	log.WithFields(log.Fields{"module": "echo", "context": "web"}).Debug("starting")

	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()

		r.Use(gin.Recovery())
		r.Use(http.LoggerMiddleware())

		initializeRoutes(r)

		if err := r.Run(fmt.Sprintf("%s:%d", s.bind, s.port)); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()

	log.WithFields(log.Fields{"module": "echo", "context": "web"}).Debug("exiting")
}

func (s *Service) runWorker(ctx context.Context, wg *sync.WaitGroup) {
	var err error

	defer wg.Done()

	log.WithFields(log.Fields{"module": "echo", "context": "worker"}).Debug("starting")

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

	log.WithFields(log.Fields{"module": "echo", "context": "worker"}).Debug("exiting")
}

func initializeRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthCheckHandler)
	}
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
