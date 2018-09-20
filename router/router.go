package router

import (
	"holy-war-web/api/controllers"
	"holy-war-web/api/middleware"
	"github.com/gorilla/mux"
	"net/http"
	)

func Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/registration", controllers.Registration).Methods("POST")

	router.HandleFunc("/login", controllers.Login).Methods("POST")

	router.HandleFunc("/graphql",  middleware.Auth(controllers.GraphQl)).Methods("GET")

	//Read does not need protection
	router.HandleFunc("/read", controllers.Read).Methods("GET")

	router.Handle("/create", middleware.Auth(controllers.Create)).Methods("POST")

	router.Handle("/update", middleware.Auth(controllers.Update)).Methods("PUT")

	router.Handle("/delete", middleware.Auth(controllers.Del)).Methods("DELETE")

	router.HandleFunc("/graphql/ws", controllers.GraphQlWs)

	return router
}
