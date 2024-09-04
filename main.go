package main

import (
	"cpe/calendar/ical"
	"cpe/calendar/request"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/calendar.ics", generateICSHandler).Methods("GET")
	mux := http.NewServeMux()

	mux.Handle("/ics/", r)

	mux.Handle("/", http.FileServer(http.Dir("static")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateICSHandler(w http.ResponseWriter, r *http.Request) {
	// Set start and end times (these could be retrieved from request parameters if needed)
	start := "1725228000000" // Example start timestamp
	end := "1728684000000"   // Example end timestamp

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

	// Step 3: Generate the iCal file
	icsContent := ical.GenerateICS(events)

	// Step 4: Write the iCal file to the response
	w.Header().Set("Content-Type", "text/calendar")
	w.Write([]byte(icsContent))
}
