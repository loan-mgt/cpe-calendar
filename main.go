package main

import (
	"cpe/calendar/ical"
	"cpe/calendar/request"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", generateICSHandler)
	http.ListenAndServe(":8080", nil)
}

func generateICSHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Fetch data from the source
	data, err := request.FetchData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	// Step 2: Parse the fetched data
	events, err := ical.ParseEvents([]byte(data))
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
