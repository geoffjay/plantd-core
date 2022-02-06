package main

import "github.com/gin-gonic/gin"

func initializeRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/jobs", submitJob)
	}
}
