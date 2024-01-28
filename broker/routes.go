package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/geoffjay/plantd/core/mdp"

	"github.com/gin-gonic/gin"
)

type serviceErrors struct {
	ErrorCount int    `json:"count"`
	LastError  string `json:"last"`
}

type mockWorker struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Service     string `json:"service"`
	Description string `json:"description"`
}

var mockWorkers = []mockWorker{
	{
		ID:          1,
		Name:        "Foo",
		Service:     "org.plantd.dev.Foo",
		Description: "The first module.",
	},
	{
		ID:          2,
		Name:        "Bar",
		Service:     "org.plantd.dev.Bar",
		Description: "The second module.",
	},
	{
		ID:          3,
		Name:        "Baz",
		Service:     "org.plantd.dev.Baz",
		Description: "The third module.",
	},
}

func initializeRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Title": "Main website",
		})
	})

	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", statusHandler)
		v1.GET("/hello", helloHandler)
		v1.GET("/errors", errorsHandler)
		v1.GET("/workers", getWorkersHandler)
		v1.GET("/workers/:id", getWorkerHandler)
		v1.GET("/info", getWorkerInfoHandler)
	}
}

func statusHandler(c *gin.Context) {
	status := fmt.Sprintf(`{"status":"%s"}`, GetStatus())
	c.String(200, status)
}

func helloHandler(c *gin.Context) {
	c.String(200, `{"message": "hello, hello, hello"}`)
}

func errorsHandler(c *gin.Context) {
	errorInfo := serviceErrors{
		ErrorCount: GetErrorCount(),
		LastError:  fmt.Sprintf("%s", GetLastError()),
	}
	c.JSON(200, errorInfo)
}

func getWorkerHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
	if id < 0 || id > len(mockWorkers) {
		c.JSON(http.StatusNotFound, nil)
	}
	c.JSON(http.StatusOK, mockWorkers[id])
}

func getWorkersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, mockWorkers)
}

func getWorkerInfoHandler(c *gin.Context) {
	_mq, _ := c.Get("mq")
	mq := _mq.(*mdp.Broker)
	c.JSON(http.StatusOK, mq.GetWorkerInfo())
}
