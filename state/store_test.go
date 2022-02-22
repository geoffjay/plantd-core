package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	scope string
	store *Store
}

func (s *StoreTestSuite) SetupSuite() {
	s.scope = "org.plantd.State.Test"
	s.store = NewStore()
	if err := s.store.Load("/tmp/test.db"); err != nil {
		panic(err)
	}
	if err := s.store.CreateScope(s.scope); err != nil {
		panic(err)
	}
}

func (s *StoreTestSuite) TearDownSuite() {
	if err := s.store.DeleteScope(s.scope); err != nil {
		panic(err)
	}
	s.store.Unload()
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (s *StoreTestSuite) TestStore_GetMissingKey() {
	value, err := s.store.Get("org.plantd.State.Test", "missing")
	suite.Nil(err)
	suite.Equal(value, "")
}

func (s *StoreTestSuite) TestStore_SetGet() {
	err := s.store.Set("org.plantd.State.Test", "foo", "bar")
	suite.Nil(err)
	value, err := s.store.Get("org.plantd.State.Test", "foo")
	suite.Nil(err)
	suite.Equal(value, "bar")
}

func (s *StoreTestSuite) TestStore_Scope() {
	var err error
	err = s.store.CreateScope("test")
	suite.Nil(err, err)
	err = s.store.DeleteScope("fake")
	suite.NotNil(err, err)
	err = s.store.DeleteScope("test")
	suite.Nil(err, err)
}
