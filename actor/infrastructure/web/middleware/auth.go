package middleware

import (
	"shajaro/actor/config"

	"strings"

	"net/http"
	"reflect"
	"runtime"
	"shajaro/actor/helper"

	log "github.com/dynastymasra/gochill"
	"github.com/gin-gonic/gin"
)

func ValidateScope(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopeHeader := c.Request.Header.Get(config.ScopeHeader)
		scopes := strings.Split(scopeHeader, " ")

		if !isAllowedAccess(scope, scopes) {
			c.Header("Content-Type", "application/json")

			log.Warn(log.Msg("User doesn't have right scope access"), log.O("project", config.ProjectName),
				log.O("package", runtime.FuncForPC(reflect.ValueOf(ValidateScope).Pointer()).Name()),
				log.O("version", config.Version), log.O("required", scope), log.O("actual", scopes))

			c.JSON(http.StatusForbidden, helper.FailResponse("user don't have right access"))
			c.Abort()
		}
		c.Next()
	}
}

func isAllowedAccess(scope string, scopeList []string) bool {
	for _, v := range scopeList {
		if v == scope {
			return true
		}
	}
	return false
}
