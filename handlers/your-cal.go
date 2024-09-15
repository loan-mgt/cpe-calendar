package handlers

import (
	"cpe/calendar/decrypt"
	"cpe/calendar/ical"
	"cpe/calendar/request"
	"log"
	"net/http"
	"os"
	"strings"
)

// generateICSHandler generates the ICS file and sends it in the response with a given filename
func GenerateICSHandler(w http.ResponseWriter, r *http.Request) {
	// Get start and end times from environment variables
	start := os.Getenv("START_TIMESTAMP")
	end := os.Getenv("END_TIMESTAMP")

	// Get separator from environment variable
	separator := os.Getenv("SEPARATOR")

	filename := "cpe-calendar" + ".ics"

	calendarName := "CPE Calendar"

	// Get query param 'creds'
	cryptedCreds := r.URL.Query().Get("creds")

	// Load the RSA private key
	privateKey, err := decrypt.LoadPrivateKey()
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}

	// Decrypt the message
	decryptedMessage, err := decrypt.DecryptMessage(cryptedCreds, privateKey)
	if err != nil {
		log.Printf("Error decrypting message: %v", err)
	}

	// Split the decrypted message using the separator
	parts := strings.Split(decryptedMessage, separator)
	if len(parts) < 2 {
		http.Error(w, "Invalid credentials format", http.StatusBadRequest)
		return
	}
	username := parts[0]
	pass := parts[1]

	// Fetch data from the source
	data, err := request.FetchData(start, end, username, pass)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	// Parse the fetched data
	events, err := ical.ParseEvents(data)
	if err != nil {
		log.Printf("Failed to parse events: %v", err)
		http.Error(w, "Failed to parse events", http.StatusInternalServerError)
		return
	}

	// Generate the iCal file with the calendar name
	icsContent := ical.GenerateICS(events, calendarName)

	// Set headers for the iCal file response with the provided filename
	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	// Write the iCal content to the response
	w.Write([]byte(icsContent))
}

func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	// Get start and end times from environment variables
	start := os.Getenv("START_TIMESTAMP")
	end := os.Getenv("END_TIMESTAMP")

	// Get separator from environment variable
	separator := os.Getenv("SEPARATOR")

	// Get query param 'creds'
	cryptedCreds := r.URL.Query().Get("creds")

	// Load the RSA private key
	privateKey, err := decrypt.LoadPrivateKey()
	if err != nil {
		log.Printf("Error loading private key: %v", err)
		http.Error(w, "Failed to load private key", http.StatusInternalServerError)
		return
	}

	// Decrypt the message
	decryptedMessage, err := decrypt.DecryptMessage(cryptedCreds, privateKey)
	if err != nil {
		log.Printf("Error decrypting message: %v", err)
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	// Split the decrypted message using the separator
	parts := strings.Split(decryptedMessage, separator)
	if len(parts) < 2 {
		http.Error(w, "Invalid credentials format", http.StatusBadRequest)
		return
	}
	username := parts[0]
	pass := parts[1]

	// Fetch data from the source
	data, err := request.FetchData(start, end, username, pass)
	if err != nil {
		log.Printf("Failed to fetch data: %v", err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	// Parse the fetched data
	events, err := ical.ParseEvents(data)
	if err != nil {
		log.Printf("Failed to parse events: %v", err)
		http.Error(w, "Failed to parse events", http.StatusInternalServerError)
		return
	}

	if len(events) > 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
