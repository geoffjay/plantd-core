package main

import "github.com/gin-gonic/gin"

func brokerMiddleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("mq", service.broker)
	}
}
