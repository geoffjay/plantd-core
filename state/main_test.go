package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func TestMainFunc(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-V")
	err := cmd.Run()

	var e *exec.ExitError
	if ok := errors.As(err, &e); ok && !e.Success() {
		return
	}

	t.Fatalf("process ran with err %v, want exit status 1", err)
}
