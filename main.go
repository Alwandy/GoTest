package main

import (
	"log"
	"net/http"
	"time"

	"./api/fleet"
	"./pkg/database"
	"github.com/gorilla/mux"
)

/*
*  Main method to start the Deathstar
 */
func main() {
	r := mux.NewRouter()                                                // Create the router, r = router
	r.HandleFunc("/", HomeHandler)                                      // Set root path handler
	r.HandleFunc("/api/fleet", fleet.GetAllShipsHandler).Methods("GET") // To retrieve all ships in the fleet
	r.HandleFunc("/api/fleet", fleet.CreateShipHandler).Methods("POST") // To add a ship into the fleet
	r.HandleFunc("/api/fleet/id?={id}", TestHandler).Methods("GET")     // Get single ship based on ID
	r.HandleFunc("/api/fleet/class?={id}", TestHandler).Methods("GET")  // Filter ships based on class
	r.HandleFunc("/api/fleet/status?={id}", TestHandler).Methods("GET") // Filter ships based on status
	r.HandleFunc("/api/fleet/name?={id}", TestHandler).Methods("GET")   // Filter ships based on name
	r.HandleFunc("/api/fleet/id?={id}", TestHandler).Methods("PUT")     // Update existing ship
	r.HandleFunc("/api/fleet/id?={id}", TestHandler).Methods("DELETE")  // Delete ship
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")          // Make sure the server runs fine and returns 200
	http.Handle("/", r)                                                 // Lets us pass the Handlefunc for route.

	srv := &http.Server{
		Handler:      r,                // What router to use (middleware)
		Addr:         "127.0.0.1:8000", // Define Address
		WriteTimeout: 15 * time.Second, // Timeout to write, great practice to set to avoid calls getting stuck
		ReadTimeout:  15 * time.Second, // Time to read, to avoid calls getting stuck
	}
	database.Init()                 // Initialize the database and create the tables if it doesn't exist.
	log.Fatal(srv.ListenAndServe()) // Start the gateway to listen and serve on 127.0.0.1:8000, log if any issues occurs to start.
}

// Returning a welcome page for root path
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Deathstar!"))
}

//Empty handler to return WIP...
func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WIP..."))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok!"))
}
