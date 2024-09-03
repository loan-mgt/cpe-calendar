package ical

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Regular expression for splitting title into components
const regexPattern = `(?P<location>.*?)(?P<promo>[1-9][A-Z]{3,})(?P<summary>.*?)(?P<description>(( |n)[A-Z]{3,} .*)|$)`

// GenerateICS generates an ICS string from a list of events
func GenerateICS(events []Event) string {
	ics := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Your Organization//Your Product//EN\n"

	// Define the layout for parsing the datetime with a timezone offset
	const layout = "2006-01-02T15:04:05-0700"

	// Compile the regular expression
	re := regexp.MustCompile(regexPattern)

	for _, event := range events {
		// Remove newline characters from title
		cleanedTitle := strings.ReplaceAll(event.Title, "\n", " ")

		// Apply regex to split title
		matches := re.FindStringSubmatch(cleanedTitle)
		if matches == nil {
			// Handle case where regex does not match
			fmt.Println("Error parsing title:", event.Title)
			continue
		}

		// Extract components from regex matches
		location := matches[1]
		summary := matches[3]
		description := matches[4]

		// Parse the start and end times in the given time zone
		start, err := time.Parse(layout, event.Start)
		if err != nil {
			// Handle parsing error
			fmt.Println("Error parsing start time:", err)
			continue
		}
		end, err := time.Parse(layout, event.End)
		if err != nil {
			// Handle parsing error
			fmt.Println("Error parsing end time:", err)
			continue
		}

		// Convert times to UTC
		start = start.UTC()
		end = end.UTC()

		// Format times for ICS
		ics += "BEGIN:VEVENT\n"
		ics += fmt.Sprintf("UID:%s\n", event.ID)
		ics += fmt.Sprintf("DTSTART:%s\n", start.Format("20060102T150405Z"))
		ics += fmt.Sprintf("DTEND:%s\n", end.Format("20060102T150405Z"))
		ics += fmt.Sprintf("LOCATION:%s\n", location)
		ics += fmt.Sprintf("SUMMARY:%s\n", summary)
		ics += fmt.Sprintf("DESCRIPTION:%s\n", description)
		ics += "END:VEVENT\n"
	}

	ics += "END:VCALENDAR\n"

	return ics
}
