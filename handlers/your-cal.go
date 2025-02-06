package handlers

import (
	"cpe/calendar/decrypt"
	"cpe/calendar/ical"
	"cpe/calendar/logger"
	"cpe/calendar/request"
	"net/http"
	"os"
	"strings"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	logger.Log.Info().
		Msg("Health check endpoint hit, status OK")
}

// GenerateICSHandler generates the ICS file and sends it in the response with a given filename
func GenerateICSHandler(w http.ResponseWriter, r *http.Request) {
	// Get start and end times from environment variables
	start := os.Getenv("START_TIMESTAMP")
	end := os.Getenv("END_TIMESTAMP")
	separator := os.Getenv("SEPARATOR")

	// Log environment variables
	logger.Log.Info().
		Str("start", start).
		Str("end", end).
		Str("separator", separator).
		Msg("Using environment variables for start, end, and separator")

	filename := "cpe-calendar.ics"
	calendarName := "CPE Calendar"

	// Get query param 'creds'
	cryptedCreds := r.URL.Query().Get("creds")

	// Load the RSA private key
	privateKey, err := decrypt.LoadPrivateKey()
	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Error loading private key")
		http.Error(w, "Failed to load private key", http.StatusInternalServerError)
		return
	}

	// Decrypt the message
	decryptedMessage, err := decrypt.DecryptMessage(cryptedCreds, privateKey)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("cryptedCreds", cryptedCreds).
			Msg("Error decrypting message")
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	// Split the decrypted message using the separator
	parts := strings.Split(decryptedMessage, separator)
	if len(parts) < 2 {
		logger.Log.Error().
			Str("decryptedMessage", decryptedMessage).
			Msg("Invalid credentials format")
		http.Error(w, "Invalid credentials format", http.StatusBadRequest)
		return
	}
	username := parts[0]
	pass := parts[1]

	// Log successful decryption of message
	logger.Log.Info().
		Str("username", username).
		Msg("Credentials decrypted successfully")

	// Fetch data from the source
	events, err := request.FetchData(start, end, username, pass)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("username", username).
			Msg("Failed to fetch data")
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	logger.Log.Info().
		Int("eventsCount", len(events)).
		Msg("Fetched events successfully")

	// Generate the iCal file with the calendar name
	icsContent := ical.GenerateICS(events, calendarName)

	// Set headers for the iCal file response with the provided filename
	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	// Write the iCal content to the response
	w.Write([]byte(icsContent))
}

// ValidateHandler validates the credentials and checks if the login is successful
func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	separator := os.Getenv("SEPARATOR")

	// Get query param 'creds'
	cryptedCreds := r.URL.Query().Get("creds")

	// Log incoming credentials request
	logger.Log.Info().
		Str("cryptedCreds", cryptedCreds).
		Msg("Validate credentials request received")

	// Load the RSA private key
	privateKey, err := decrypt.LoadPrivateKey()
	if err != nil {
		logger.Log.Error().
			Err(err).
			Msg("Error loading private key")
		http.Error(w, "Failed to load private key", http.StatusInternalServerError)
		return
	}

	// Decrypt the message
	decryptedMessage, err := decrypt.DecryptMessage(cryptedCreds, privateKey)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("cryptedCreds", cryptedCreds).
			Msg("Error decrypting message")
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	// Split the decrypted message using the separator
	parts := strings.Split(decryptedMessage, separator)
	if len(parts) < 2 {
		logger.Log.Error().
			Str("decryptedMessage", decryptedMessage).
			Msg("Invalid credentials format")
		http.Error(w, "Invalid credentials format", http.StatusBadRequest)
		return
	}
	username := parts[0]
	pass := parts[1]

	// Log successful decryption
	logger.Log.Info().
		Str("username", username).
		Msg("Credentials decrypted successfully")

	// Fetch data to validate credentials
	_, err = request.Login(username, pass)
	if err != nil {
		logger.Log.Error().
			Err(err).
			Str("username", username).
			Msg("Failed to validate credentials")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	logger.Log.Info().
		Str("username", username).
		Msg("User validated successfully")
	w.WriteHeader(http.StatusOK)
}
