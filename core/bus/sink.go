package bus

// Sink type to function as a bus subscriber device.
type Sink struct{}

// NewSink constructs an instance of a message bus sink.
func NewSink() *Sink {
	return &Sink{}
}
