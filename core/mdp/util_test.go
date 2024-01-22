package mdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopWorker(t *testing.T) {
	workers := []*brokerWorker{
		{
			idString: "1",
		},
		{
			idString: "2",
		},
		{
			idString: "3",
		},
	}

	worker, otherWorkers := popWorker(workers)
	assert.Equal(t, worker.idString, "1")
	assert.Equal(t, len(otherWorkers), 2)
}

func TestDelWorker(t *testing.T) {
	workers := []*brokerWorker{
		{
			idString: "1",
		},
		{
			idString: "2",
		},
		{
			idString: "3",
		},
	}

	workers = delWorker(workers, workers[1])
	assert.Equal(t, len(workers), 2)
	assert.Equal(t, workers[0].idString, "1")
	assert.Equal(t, workers[1].idString, "3")
}

func TestStringArrayToByte2D(t *testing.T) {
	input := []string{"1", "2", "3", "4"}

	expected := [][]byte{
		[]byte("1"),
		[]byte("2"),
		[]byte("3"),
		[]byte("4"),
	}

	output := stringArrayToByte2D(input)
	assert.Equal(t, expected, output)
}

func TestByte2DToStringArray(t *testing.T) {
	input := [][]byte{
		[]byte("1"),
		[]byte("2"),
		[]byte("3"),
		[]byte("4"),
	}

	expected := []string{"1", "2", "3", "4"}

	output := byte2DToStringArray(input)
	assert.Equal(t, expected, output)
}
