package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type RoutesTestSuite struct {
	suite.Suite
	router   *gin.Engine
	recorder *httptest.ResponseRecorder
}

func (suite *RoutesTestSuite) SetupTest() {
	suite.router = gin.Default()
	initializeRoutes(suite.router)
	suite.recorder = httptest.NewRecorder()

	// initialize state for testing
	SetStatus("testing")
}

func TestRoutesTestSuite(t *testing.T) {
	suite.Run(t, new(RoutesTestSuite))
}

func (suite *RoutesTestSuite) TestRoutes_GetWorker() {
	req, _ := http.NewRequest("GET", "/api/v1/workers/0", nil)
	suite.router.ServeHTTP(suite.recorder, req)

	// FIXME: this is lame, should deserialize and check fields
	expected := "{\"id\":1,\"name\":\"Foo\",\"service\":\"org.plantd.dev.Foo\",\"description\":\"The first module.\"}"
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(expected, suite.recorder.Body.String())
}

func (suite *RoutesTestSuite) TestRoutes_Status() {
	req, _ := http.NewRequest("GET", "/api/v1/status", nil)
	suite.router.ServeHTTP(suite.recorder, req)

	expected := "{\"status\":\"testing\"}"
	suite.Equal(expected, suite.recorder.Body.String())
}
