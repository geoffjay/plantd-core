package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/geoffjay/plantd/core/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Service type for REST API.
type Service struct {
	port int
	bind string
}

// NewService constructs and instance of a service type.
func NewService(port int, bind string) *Service {
	return &Service{
		port,
		bind,
	}
}

func (s *Service) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

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
	log.Debug("exiting service")
}
