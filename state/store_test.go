package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	scope string
	store *Store
}

func TestStoreLoad(t *testing.T) {
	store := NewStore()
	if err := os.Mkdir("./tmp", 0664); err != nil {
		panic(err)
	}
	err := store.Load("./tmp")
	assert.NotNil(t, err)
	os.Remove("./tmp")
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
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

// nolint: typecheck
func (suite *StoreTestSuite) TestStore_GetMissingKey() {
	value, err := suite.store.Get("org.plantd.State.Test", "missing")
	suite.Nil(err)
	suite.Equal(value, "")
}

// nolint: typecheck
func (suite *StoreTestSuite) TestStore_SetGet() {
	err := suite.store.Set("org.plantd.State.Test", "foo", "bar")
	suite.Nil(err)
	value, err := suite.store.Get("org.plantd.State.Test", "foo")
	suite.Nil(err)
	suite.Equal(value, "bar")
}

// nolint: typecheck
func (suite *StoreTestSuite) TestStore_Scope() {
	var err error
	err = suite.store.CreateScope("test")
	suite.Nil(err, err)
	err = suite.store.DeleteScope("fake")
	suite.NotNil(err, err)
	err = suite.store.DeleteScope("test")
	suite.Nil(err, err)
}

func (suite *StoreTestSuite) TestStore_HasScope() {
	err := suite.store.CreateScope("test")
	suite.Nil(err, err)
	ok := suite.store.HasScope("test")
	suite.Equal(ok, true)
	err = suite.store.DeleteScope("test")
	suite.Nil(err, err)
}

func (suite *StoreTestSuite) TestStore_ListAllScope() {
	var err error
	err = suite.store.CreateScope("test1")
	suite.Nil(err, err)
	err = suite.store.CreateScope("test2")
	suite.Nil(err, err)
	scopes := suite.store.ListAllScope()
	suite.Equal(len(scopes), 3)
	suite.Equal(scopes[1], "test1")
	suite.Equal(scopes[2], "test2")
	err = suite.store.DeleteScope("test1")
	suite.Nil(err, err)
	err = suite.store.DeleteScope("test2")
	suite.Nil(err, err)
}
