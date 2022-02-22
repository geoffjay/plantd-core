package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	handler *Handler
}

type testCallback struct {
	name string
}

func (cb *testCallback) Execute(data string) ([]byte, error) {
	return []byte(""), nil
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) SetupSuite() {
	suite.handler = NewHandler()
	suite.handler.AddCallback("test", &testCallback{name: "test"})
}

func (suite *HandlerTestSuite) TestHandler_Get() {
	callback, err := suite.handler.GetCallback("test")
	suite.NotNil(callback, "callback exists for `test`")
	suite.Nil(err, "callback retrieve failed")
	callback, err = suite.handler.GetCallback("foo")
	suite.Nil(callback, "callback doesn't exist for `foo`")
	suite.NotNil(err, "callback retrieval failed for unavailable name")
	suite.Equal(err.Error(), "callback not found for foo")
}
