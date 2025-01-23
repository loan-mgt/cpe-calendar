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
)

func FetchData(start, end, username, password string) ([]types.Event, error) {

	token, err := Login(username, password)
	if err != nil {
		return nil, err
	}

	body, err := getCalendar(token, start, end)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Login(username, password string) (types.TokenResponse, error) {
	// Prepare the login request
	urlStr := "https://mycpe.cpe.fr/mobile/login"
	loginData := map[string]string{
		"login":    username,
		"password": password,
	}

	// Marshal login data to JSON
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return types.TokenResponse{}, fmt.Errorf("failed to marshal login data: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonData))
	if err != nil {
		return types.TokenResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 15; sdk_gphone64_x86_64 Build/AE3A.240806.005)")
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return types.TokenResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return types.TokenResponse{}, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Read and unmarshal the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.TokenResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var formattedResp types.TokenResponse
	if err := json.Unmarshal(body, &formattedResp); err != nil {
		return types.TokenResponse{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return formattedResp, nil
}

func getCalendar(token types.TokenResponse, startUnix, endUnix string) ([]types.Event, error) {
	// Define the base URL and query parameters
	baseURL := "https://mycpe.cpe.fr/mobile/mon_planning"

	startTime, err := unixToDateTime(startUnix)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start time: %w", err)
	}

	endTime, err := unixToDateTime(endUnix)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end time: %w", err)
	}

	query := fmt.Sprintf("?date_debut=%s&date_fin=%s", startTime, endTime)

	fmt.Println("final url ", baseURL+query)

	// Create the GET request
	req, err := http.NewRequest("GET", baseURL+query, nil)
	if err != nil {
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
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle gzip encoding if necessary
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer reader.(*gzip.Reader).Close()
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response into the events slice
	var events []types.Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return events, nil
}

func unixToDateTime(rawTime string) (string, error) {
	start, err := strconv.ParseInt(rawTime, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid unix time: %v", err)
	}

	start = start / 1000

	dateTime := time.Unix(start, 0)

	// Format the time to match the query format ("YYYY-MM-DD")
	return dateTime.Format("2006-01-02"), nil

}
