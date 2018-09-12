package router

import (
	"awesomeProject/api/controllers"
	"awesomeProject/api/middleware"
	"github.com/gorilla/mux"
	"net/http"
	)

func Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/registration", controllers.Registration).Methods("POST")

	router.HandleFunc("/login", controllers.Login).Methods("POST")

	router.HandleFunc("/graphql",  middleware.MiddlewareAuth(controllers.GraphQl)).Methods("GET")

	//Read does not need protection
	router.HandleFunc("/read", controllers.Read).Methods("GET")

	router.Handle("/create", middleware.MiddlewareAuth(controllers.Create)).Methods("POST")

	router.Handle("/update", middleware.MiddlewareAuth(controllers.Update)).Methods("PUT")

	router.Handle("/delete", middleware.MiddlewareAuth(controllers.Del)).Methods("DELETE")

	router.HandleFunc("/graphql/ws", controllers.GraphQlWs)

	return router
}
