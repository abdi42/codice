package main

import (
	"fmt"
	"github.com/abdullahi/codice/controllers"
	"github.com/abdullahi/codice/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},                                     // All origins
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}, // Allowing only get, just an example
	})

	router.Use(services.JwtAuthentication)

	//Workspaces
	router.HandleFunc("/workspace/{id}/files", controllers.GetWorkspaceFiles).Methods("GET")
	router.HandleFunc("/workspace/{id}", controllers.GetWorkSpace).Methods("GET")
	router.HandleFunc("/workspace/{id}/run", controllers.RunProgram).Methods("PUT")
	router.HandleFunc("/workspace/create", controllers.CreateWorkspace).Methods("POST")

	//Users & Auth
	router.HandleFunc("/users/new", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/users/login", controllers.Authenticate).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println("Running service at port: " + port)

	err := http.ListenAndServe(":"+port, c.Handler(router))

	if err != nil {
		fmt.Print(err)
	}
}
