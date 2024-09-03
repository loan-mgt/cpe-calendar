package ical

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

// Event struct to hold individual event data
type Event struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Start     string `json:"start"`
	End       string `json:"end"`
	AllDay    bool   `json:"allDay"`
	Editable  bool   `json:"editable"`
	ClassName string `json:"className"`
}

// Define a structure to parse the XML partial-response
type PartialResponse struct {
	XMLName xml.Name `xml:"partial-response"`
	Changes []Update `xml:"changes>update"`
}

// Define a structure for each update within the XML
type Update struct {
	ID    string `xml:"id,attr"`
	CDATA string `xml:",innerxml"`
}

// Function to parse events from the provided data
func ParseEvents(data []byte) ([]Event, error) {
	var partialResponse PartialResponse

	// Parse the XML data
	err := xml.Unmarshal(data, &partialResponse)
	if err != nil {
		return nil, err
	}

	// Find the update with the event JSON data
	for _, update := range partialResponse.Changes {
		if strings.Contains(update.CDATA, `"events"`) {
			// Extract the JSON within the CDATA section
			start := strings.Index(update.CDATA, `[`)       // Start of the array
			end := strings.LastIndex(update.CDATA, `]`) + 1 // End of the array
			if start == -1 || end == -1 {
				return nil, errors.New("failed to find JSON in CDATA section")
			}
			jsonData := update.CDATA[start:end]

			// Remove CDATA markers if present
			jsonData = strings.TrimPrefix(jsonData, "[CDATA[")
			jsonData = strings.TrimSuffix(jsonData, "]]")

			// Debugging: Print extracted JSON data
			fmt.Println("Extracted JSON Data:", jsonData)

			// Parse the JSON data
			var events struct {
				Events []Event `json:"events"`
			}
			err = json.Unmarshal([]byte(jsonData), &events)
			if err != nil {
				return nil, err
			}

			// Assuming the array contains a single object with the `events` key
			return events.Events, nil
		}
	}

	return nil, errors.New("no events found")
}
