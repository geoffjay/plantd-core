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

func (suite *StoreTestSuite) SetupSuite() {
	suite.scope = "org.plantd.State.Test"
	suite.store = NewStore()
	if err := suite.store.Load("/tmp/test.db"); err != nil {
		panic(err)
	}
	if err := suite.store.CreateScope(suite.scope); err != nil {
		panic(err)
	}
}

func (suite *StoreTestSuite) TearDownSuite() {
	if err := suite.store.DeleteScope(suite.scope); err != nil {
		panic(err)
	}
	suite.store.Unload()
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (suite *StoreTestSuite) TestStore_GetMissingKey() {
	value, err := suite.store.Get("org.plantd.State.Test", "missing")
	suite.Nil(err)
	suite.Equal(value, "")
}

func (suite *StoreTestSuite) TestStore_SetGet() {
	err := suite.store.Set("org.plantd.State.Test", "foo", "bar")
	suite.Nil(err)
	value, err := suite.store.Get("org.plantd.State.Test", "foo")
	suite.Nil(err)
	suite.Equal(value, "bar")
}

func (suite *StoreTestSuite) TestStore_Scope() {
	var err error
	err = suite.store.CreateScope("test")
	suite.Nil(err, err)
	err = suite.store.DeleteScope("fake")
	suite.NotNil(err, err)
	err = suite.store.DeleteScope("test")
	suite.Nil(err, err)
}
