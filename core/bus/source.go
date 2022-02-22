package bus

// Source type to function as a bus publisher device.
type Source struct{}

// NewSource constructs an instance of a message bus sink.
func NewSource() *Source {
	return &Source{}
}
