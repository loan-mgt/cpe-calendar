package main

import (
	"cpe/calendar/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
}

func main() {
	r := mux.NewRouter()

	// Handle the calendar.ics route with default filename "calendar.ics"
	r.HandleFunc("/your-cpe-calendar.ics", generate3IRCHandler).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))

	// Use the router in the http server
	log.Fatal(http.ListenAndServe(":8080", r))
}

// generate3IRCHandler is a wrapper around generateICSHandler that uses a specific filename
func generate3IRCHandler(w http.ResponseWriter, r *http.Request) {
	// Call generateICSHandler with the specific filename and calendar name
	handlers.GenerateICSHandler(w, r, "3irc_calendar.ics", "3IRC Calendar")
}
