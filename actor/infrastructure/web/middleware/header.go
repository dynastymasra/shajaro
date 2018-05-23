package middleware

import (
	"sirius/actor/config"

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
