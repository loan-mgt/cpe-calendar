package request

import (
	"bytes"
	"compress/gzip"
	"cpe/calendar/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
	"cpe/calendar/logger"
)

func FetchData(start, end, username, password string) ([]types.Event, error) {
	// Log the operation with context about the start and end times
	logger.Log.Info().
		Str("username", username).
		Str("start", start).
		Str("end", end).
		Msg("Fetching data from CPE calendar")

	token, err := Login(username, password)
	if err != nil {
		logger.Log.Error().
			Str("username", username).
			Err(err).
			Msg("Failed to login")
		return nil, err
	}

	body, err := getCalendar(token, start, end)
	if err != nil {
		logger.Log.Error().
			Str("username", username).
			Err(err).
			Msg("Failed to fetch calendar data")
		return nil, err
	}

	logger.Log.Info().
		Str("username", username).
		Msg("Data fetched successfully")
	return body, nil
}

func Login(username, password string) (types.TokenResponse, error) {
	// Log the login request with username context
	logger.Log.Info().
		Str("username", username).
		Msg("Initiating login request")

	// Prepare the login request
	urlStr := "https://mycpe.cpe.fr/mobile/login"
	loginData := map[string]string{
		"login":    username,
		"password": password,
	}

	// Marshal login data to JSON
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		logger.Log.Error().
			Str("username", username).
			Err(err).
			Msg("Failed to marshal login data")
		return types.TokenResponse{}, fmt.Errorf("failed to marshal login data: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log.Error().
			Str("username", username).
			Err(err).
			Msg("Failed to create login request")
		return types.TokenResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 15; sdk_gphone64_x86_64 Build/AE3A.240806.005)")
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error().
			Str("username", username).
			Err(err).
			Msg("Login request failed")
		return types.TokenResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		logger.Log.Error().
			Str("username", username).
			Int("statusCode", resp.StatusCode).
			Msg("Received non-200 response for login")
		return types.TokenResponse{}, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Read and unmarshal the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error().
			Str("username", username).
			Err(err).
			Msg("Failed to read login response body")
		return types.TokenResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var formattedResp types.TokenResponse
	if err := json.Unmarshal(body, &formattedResp); err != nil {
		logger.Log.Error().
			Str("username", username).
			Err(err).
			Msg("Failed to parse login response JSON")
		return types.TokenResponse{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	logger.Log.Info().
		Str("username", username).
		Msg("Login successful")
	return formattedResp, nil
}

func getCalendar(token types.TokenResponse, startUnix, endUnix string) ([]types.Event, error) {
	// Log the request to fetch calendar data with the token and time context
	logger.Log.Info().
		Str("token", token.Normal).
		Str("startUnix", startUnix).
		Str("endUnix", endUnix).
		Msg("Fetching calendar data")

	// Define the base URL and query parameters
	baseURL := "https://mycpe.cpe.fr/mobile/mon_planning"

	startTime, err := unixToDateTime(startUnix)
	if err != nil {
		logger.Log.Error().
			Str("startUnix", startUnix).
			Err(err).
			Msg("Failed to parse start time")
		return nil, fmt.Errorf("failed to parse start time: %w", err)
	}

	endTime, err := unixToDateTime(endUnix)
	if err != nil {
		logger.Log.Error().
			Str("endUnix", endUnix).
			Err(err).
			Msg("Failed to parse end time")
		return nil, fmt.Errorf("failed to parse end time: %w", err)
	}

	query := fmt.Sprintf("?date_debut=%s&date_fin=%s", startTime, endTime)
	logger.Log.Debug().
		Str("finalURL", baseURL+query).
		Msg("Generated final URL")

	// Create the GET request
	req, err := http.NewRequest("GET", baseURL+query, nil)
	if err != nil {
		logger.Log.Error().
			Str("finalURL", baseURL+query).
			Err(err).
			Msg("Failed to create calendar request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers to match the curl request
	req.Header.Add("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 15; sdk_gphone64_x86_64 Build/AE3A.240806.005)")
	req.Header.Add("Authorization", "Bearer "+token.Normal)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Host", "mycpe.cpe.fr")

	// Send the GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error().
			Str("finalURL", baseURL+query).
			Err(err).
			Msg("Request failed to get calendar data")
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle gzip encoding if necessary
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			logger.Log.Error().
				Str("finalURL", baseURL+query).
				Err(err).
				Msg("Failed to create gzip reader")
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.(*gzip.Reader).Close()
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		logger.Log.Error().
			Str("finalURL", baseURL+query).
			Int("statusCode", resp.StatusCode).
			Msg("Received non-200 response while fetching calendar")
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(reader)
	if err != nil {
		logger.Log.Error().
			Str("finalURL", baseURL+query).
			Err(err).
			Msg("Failed to read calendar response body")
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response into the events slice
	var events []types.Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		logger.Log.Error().
			Str("finalURL", baseURL+query).
			Err(err).
			Msg("Failed to parse calendar JSON response")
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	logger.Log.Info().
		Str("token", token.Normal).
		Str("startUnix", startUnix).
		Str("endUnix", endUnix).
		Msg("Calendar data fetched successfully")
	return events, nil
}

func unixToDateTime(rawTime string) (string, error) {
	start, err := strconv.ParseInt(rawTime, 10, 64)
	if err != nil {
		logger.Log.Error().
			Str("rawTime", rawTime).
			Err(err).
			Msg("Invalid Unix time format")
		return "", fmt.Errorf("invalid unix time: %v", err)
	}

	start = start / 1000
	dateTime := time.Unix(start, 0)

	// Format the time to match the query format ("YYYY-MM-DD")
	return dateTime.Format("2006-01-02"), nil
}
