package web

import (
	"fmt"
	"net/http"
	"shajaro/actor/config"

	"shajaro/actor/infrastructure/web/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "{\"status\": \"failed\", \"message\": \"%v\"}", config.ErrEndpointNotFound)
	})

	router.HandleFunc("/v1/healthz", controller.HealthzController).Methods(http.MethodGet)
	router.HandleFunc("/v1/healthz", controller.HealthzController).Methods(http.MethodHead)

	return router
}
