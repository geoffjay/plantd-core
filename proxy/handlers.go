package main

import (
	"github.com/gin-gonic/gin"
)

func submitJob(c *gin.Context) {
	response := `{"msg":"unimplemented"}`
	c.String(200, response)
}
