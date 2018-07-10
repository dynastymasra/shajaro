package middleware

import (
	"github.com/dynastymasra/shajaro/actor/config"

	"strings"

	"net/http"
	"reflect"
	"runtime"

	"github.com/dynastymasra/shajaro/actor/helper"

	"fmt"

	log "github.com/dynastymasra/gochill"
	"github.com/urfave/negroni"
)

func ValidateScope(scope string) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		scopeHeader := r.Header.Get(config.ScopeHeader)
		scopes := strings.Split(scopeHeader, " ")

		if !isAllowedAccess(scope, scopes) {
			log.Warn(log.Msg("User doesn't have right scope access"), log.O("project", config.ProjectName),
				log.O("package", runtime.FuncForPC(reflect.ValueOf(ValidateScope).Pointer()).Name()),
				log.O("version", config.Version), log.O("required", scope), log.O("actual", scopes))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, helper.FailResponse("user don't have right access").Stringify())
			return
		}
		next(w, r)
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
