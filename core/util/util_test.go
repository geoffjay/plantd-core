package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopMsg(t *testing.T) {
	msg := [][]string{
		{"1a", "2a", "3a"},
		{"1b", "2b", "3b"},
		{"1c", "2c", "3c"},
	}
	head, tail := PopMsg(msg)

	assert.Equal(t, []string{"1a", "2a", "3a"}, head)
	assert.Equal(t, 2, len(tail))

	head, tail = PopMsg([][]string{{"foo", "bar", "baz"}})

	assert.Equal(t, []string{"foo", "bar", "baz"}, head)
	assert.NotNil(t, tail)
	assert.Equal(t, [][]string{}, tail)
}

func TestPopStr(t *testing.T) {
	head, tail := PopStr([]string{"foo", "bar", "baz"})

	assert.Equal(t, "foo", head)
	assert.Equal(t, 2, len(tail))
	assert.Equal(t, "bar", tail[0])
	assert.Equal(t, "baz", tail[1])

	head, tail = PopStr([]string{"foo"})
	assert.Equal(t, "foo", head)
	assert.Equal(t, 0, len(tail))
	assert.NotNil(t, tail)
	assert.Equal(t, []string{}, tail)
}

func TestUnwrap(t *testing.T) {
	head, tail := Unwrap([]string{"foo", "bar", "baz"})

	assert.Equal(t, "foo", head)
	assert.Equal(t, 2, len(tail))
	assert.Equal(t, "bar", tail[0])
	assert.Equal(t, "baz", tail[1])

	head, tail = Unwrap([]string{"foo", "", "bar", "baz"})

	assert.Equal(t, "foo", head)
	assert.Equal(t, 2, len(tail))
	assert.Equal(t, "bar", tail[0])
	assert.Equal(t, "baz", tail[1])
}

func TestKeys(t *testing.T) {
	list := make(map[string]string)
	list["foo"] = "a"
	list["bar"] = "b"
	list["baz"] = "c"

	keys := Keys(list)

	assert.ElementsMatch(t, []string{"foo", "bar", "baz"}, keys)
}
