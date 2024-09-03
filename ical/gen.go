package ical

import (
	"fmt"
	"time"
)

func GenerateICS(events []Event) string {
	ics := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Your Organization//Your Product//EN\n"

	// Define the layout for parsing the datetime with a timezone offset
	const layout = "2006-01-02T15:04:05-0700"

	for _, event := range events {
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
		ics += fmt.Sprintf("SUMMARY:%s\n", event.Title)
		ics += "END:VEVENT\n"
	}

	ics += "END:VCALENDAR\n"

	return ics
}
