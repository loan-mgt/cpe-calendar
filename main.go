package main

import (
	"cpe/calendar/ical"
	"cpe/calendar/request"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	r := mux.NewRouter()

	// Handle the calendar.ics route with default filename "calendar.ics"
	r.HandleFunc("/3irc.ics", generate3IRCHandler).Methods("GET")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static/")))

	// Use the router in the http server
	log.Fatal(http.ListenAndServe(":8080", r))
}

// generateICSHandler generates the ICS file and sends it in the response with a given filename
func generateICSHandler(w http.ResponseWriter, _ *http.Request, filename, calendarName string) {
	// Get start and end times from environment variables
	start := os.Getenv("START_TIMESTAMP")
	end := os.Getenv("END_TIMESTAMP")

	// Step 1: Fetch data from the source using the updated FetchData function
	data, err := request.FetchData(start, end)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	log.Printf("Fetched data: %s\n======================", data)

	// Step 2: Parse the fetched data
	events, err := ical.ParseEvents(data)
	if err != nil {
		log.Printf("Failed to parse events: %v", err)
		http.Error(w, "Failed to parse events", http.StatusInternalServerError)
		return
	}

	// Step 3: Generate the iCal file with the calendar name
	icsContent := ical.GenerateICS(events, calendarName)

	// Step 4: Set headers for the iCal file response with the provided filename
	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	// Step 5: Write the iCal content to the response
	w.Write([]byte(icsContent))
}

// generate3IRCHandler is a wrapper around generateICSHandler that uses a specific filename
func generate3IRCHandler(w http.ResponseWriter, r *http.Request) {
	// Call generateICSHandler with the specific filename and calendar name
	generateICSHandler(w, r, "3irc_calendar.ics", "3IRC Calendar")
}
