package main

import (
	"sync"
)

func SetStatus(value string) {
	s.setStatus(value)
}

func GetStatus() string {
	return s.getStatus()
}

func SetLastError(err error) {
	s.setLastError(err)
}

func GetErrorCount() int {
	return s.getErrorCount()
}

func GetLastError() error {
	return s.getLastError()
}

type state struct {
	sync.RWMutex
	status     string
	errorCount int
	lastError  error
}

func (s *state) setStatus(value string) {
	s.Lock()
	s.status = value
	s.Unlock()
}

func (s *state) getStatus() string {
	s.RLock()
	defer s.RUnlock()
	return s.status
}

func (s *state) setLastError(err error) {
	s.Lock()
	s.lastError = err
	s.errorCount++
	s.Unlock()
}

func (s *state) getErrorCount() int {
	s.RLock()
	defer s.RUnlock()
	return s.errorCount
}

func (s *state) getLastError() error {
	s.RLock()
	defer s.RUnlock()
	return s.lastError
}

var s = &state{}
