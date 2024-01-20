package bus

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	czmq "github.com/zeromq/goczmq/v4"
)

// Sink type to function as a bus subscriber device.
type Sink struct {
	endpoint string
	filter   string
	running  bool
	handler  *SinkHandler
}

// SinkHandler defines the type of a callback.
type SinkHandler struct {
	Callback SinkCallback
}

// SinkCallback defines an interface for a callback handler.
type SinkCallback interface {
	Handle(data []byte) error
}

// NewSink constructs an instance of a message bus sink.
func NewSink(endpoint, filter string) *Sink {
	return &Sink{
		endpoint: endpoint,
		filter:   filter,
		running:  false,
	}
}

func (s *Sink) defaultFields(err error) log.Fields {
	fields := log.Fields{
		"endpoint": s.endpoint,
		"filter":   s.filter,
	}
	if err != nil {
		fields["err"] = err
	}
	return fields
}

// SetHandler sets the message handler for the sink to use.
func (s *Sink) SetHandler(handler *SinkHandler) {
	s.handler = handler
}

// Run the sink routine.
func (s *Sink) Run(ctx context.Context, wg *sync.WaitGroup) {
	var subscriber *czmq.Sock
	var poller *czmq.Poller
	var err error

	defer wg.Done()

	if subscriber, err = czmq.NewSub(s.endpoint, s.filter); err != nil {
		log.WithFields(s.defaultFields(err)).Panic("subscriber create")
	}
	log.WithFields(s.defaultFields(nil)).Debug("created message queue sink socket")
	defer subscriber.Destroy()

	if poller, err = czmq.NewPoller(); err != nil {
		log.WithFields(log.Fields{"err": err}).Panic("poller create")
	}
	log.Debug("created message queue sink poller")
	defer poller.Destroy()

	if err = poller.Add(subscriber); err != nil {
		log.WithFields(s.defaultFields(err)).Panic("poller setup")
	}
	log.Debug("successfully added subscriber to poller")

	s.running = true

	go func() {
		for s.running {
			log.WithFields(log.Fields{"socket": s.endpoint}).Trace("waiting for data...")
			socket, err := poller.Wait(1000)
			if err != nil {
				break
			}

			if socket == nil {
				log.Trace("poller timeout reached")
				continue
			}

			log.Trace("poller received data")
			data, _, rerr := socket.RecvFrame()
			// TODO while more flag (_) is set
			if rerr != nil {
				log.Error(rerr)
				continue
			}
			log.Trace("handling received data")
			if err = s.handler.Callback.Handle(data); err != nil {
				log.WithFields(log.Fields{"error": err}).Error("failed to handle message")
			}
		}
	}()

	<-ctx.Done()
	log.Debug("sink received shutdown")
	s.Stop()
}

// Stop sets the flag to shutdown the loop handling messages.
func (s *Sink) Stop() {
	s.running = false
}

// Running is used to check if the message handler should be running.
func (s *Sink) Running() bool {
	return s.running
}
