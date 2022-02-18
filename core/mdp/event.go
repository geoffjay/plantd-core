package mdp

const (
	// BrokerEvent indicates that the event originated from the broker.
	BrokerEvent = iota + 1

	// ClientEvent indicates that the event originated from the client.
	ClientEvent

	// WorkerEvent indicates that the event originated from the worker.
	WorkerEvent
)

// Event instances are passed up through a channel.
type Event struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}

// NewBrokerEvent instantiates an event with the type set to broker.
func NewBrokerEvent(message string) Event {
	return Event{
		Type:    BrokerEvent,
		Message: message,
	}
}
