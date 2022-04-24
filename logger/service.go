package main

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	manager *Manager
}

// NewService constructs and instance of a service type.
func NewService() *Service {
	return &Service{
		manager: NewManager(),
	}
}

func (s *Service) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	wg.Add(1)
	go s.manager.Run(ctx, wg)

	<-ctx.Done()
	log.Debug("exiting service")
}
