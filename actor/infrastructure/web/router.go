package web

import (
	"fmt"
	"net/http"

	"github.com/dynastymasra/shajaro/actor/config"

	"github.com/dynastymasra/shajaro/actor/infrastructure/web/controller"

	"github.com/dynastymasra/shajaro/actor/infrastructure/web/controller/actor"

	"github.com/dynastymasra/shajaro/actor/infrastructure/web/middleware"

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

	router.HandleFunc("/v1/countries", actor.CountryListController).Methods(http.MethodGet)

	router.HandleFunc("/v1/register", actor.RegisterController).Methods(http.MethodPost)

	router.HandleFunc("/v1/login", actor.LoginController).Methods(http.MethodPost)

	router.Handle("/v1/actor", negroni.New(
		negroni.HandlerFunc(middleware.ValidateScope(config.ActorRead)),
		negroni.WrapFunc(actor.GetUserByIDController),
	)).Methods(http.MethodGet)

	router.Handle("/v1/actor", negroni.New(
		negroni.HandlerFunc(middleware.ValidateScope(config.ActorUpdate)),
		negroni.WrapFunc(actor.UpdateUserController),
	)).Methods(http.MethodPut)

	router.Handle("/v1/password", negroni.New(
		negroni.HandlerFunc(middleware.ValidateScope(config.ActorUpdate)),
		negroni.WrapFunc(actor.UpdatePasswordController),
	)).Methods(http.MethodPut)

	router.Handle("/v1/actor", negroni.New(
		negroni.HandlerFunc(middleware.ValidateScope(config.ActorDelete)),
		negroni.WrapFunc(actor.DeleteUserController),
	)).Methods(http.MethodDelete)

	return router
}
