package web

import (
	"fmt"
	"net/http"
	"shajaro/actor/config"

	"shajaro/actor/infrastructure/web/controller"

	"shajaro/actor/infrastructure/web/controller/actor"

	"shajaro/actor/infrastructure/web/middleware"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
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

	router.Handle("/v1/countries", negroni.New(
		negroni.HandlerFunc(middleware.ValidateScope(config.ActorRead)),
		negroni.WrapFunc(actor.CountryListController),
	)).Methods(http.MethodGet)

	return router
}
