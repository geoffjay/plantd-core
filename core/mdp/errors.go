package mdp

import (
	"errors"
)

var (
	errPermanent = errors.New("permanent error, abandoning request")
)
