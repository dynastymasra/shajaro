package middleware

import (
	"shajaro/actor/config"

	"net/http"
	"shajaro/actor/helper"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/random"
)

func RequestKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(config.HeaderRequestID)
		if len(requestID) > 0 {
			c.Set(config.TraceKey, requestID)
			c.Next()
		}
		c.Set(config.TraceKey, random.String(11))
		c.Next()
	}
}

func RequestType() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentType := c.Request.Header.Get("Content-Type")
		httpMethod := c.Request.Method

		if httpMethod == http.MethodPost || httpMethod == http.MethodPut {
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, helper.FailResponse("content type is not supported"))
				c.Abort()
			}
		}
		c.Next()
	}
}
