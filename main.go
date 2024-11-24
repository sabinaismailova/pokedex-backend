package main

import (
	"log"
	"net/http"
	"pokedex/pokemon/handler"
)

func main() {
	// Serve locally using the Handler function from the handler package
	http.HandleFunc("/", handler.Handler)

	// Start the server locally on port 8080
	log.Println("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
