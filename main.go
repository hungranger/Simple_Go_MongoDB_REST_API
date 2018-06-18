package main

import (
	"Simple_Go_MongoDB_REST_API/store"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	router := store.NewRouter() // create routes

	cssHandler := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", cssHandler))

	// These two lines are important if you're designing a front-end to utilise this API methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	// Launch server with CORS validations
	http.Handle("/", handlers.CORS(allowedOrigins, allowedMethods)(router))
	log.Fatal(http.ListenAndServe(getPort(), nil))

}

func getPort() string {
	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":80"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return port
}
