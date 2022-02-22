package mdp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var broker *Broker
var done chan bool

// TestMain will execute each test.
func TestMain(m *testing.M) {
	setUp("tcp://127.0.0.1:5555")
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}

// setUp initializes the broker.
func setUp(endpoint string) {
	var err error
	done = make(chan bool, 1)
	if broker, err = NewBroker(endpoint); err != nil || broker == nil {
		panic("failed to initialize a new broker instance")
	}
	if err = broker.Bind(); err != nil {
		panic("failed to bind to endpoint")
	}
	go broker.Run(done)
}

// tearDown closes the broker.
func tearDown() {
	<-done
	if err := broker.Close(); err != nil {
		panic("failed to close broker endpoint connection")
	}
}

func TestBrokerClient(t *testing.T) {
	worker, err := NewWorker("tcp://localhost:5555", "test")
	assert.Nil(t, err)

	client, err := NewClient("tcp://localhost:5555")
	assert.Nil(t, err)

	err = client.Send("test", "foo")
	assert.Nil(t, err)

	var reply []string
	request, err := worker.Recv(reply)
	reply = request

	assert.NotNil(t, reply)
	assert.Nil(t, err)
	expected := []string{"foo"}
	assert.Equal(t, expected, request)

	response, err := client.Recv()
	assert.Nil(t, err)
	assert.Equal(t, expected, response)
}
