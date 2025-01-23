package ical

import (
	"cpe/calendar/types"
	"fmt"
	"time"
)

// GenerateICS generates an ICS string from a list of events
func GenerateICS(events []types.Event, calendarName string) string {
	ics := "BEGIN:VCALENDAR\n"
	ics += "VERSION:2.0\n"
	ics += "PRODID:-//github.com/qypol342 //CPE Calendar//EN\n"
	ics += fmt.Sprintf("NAME:%s\n", calendarName)
	ics += fmt.Sprintf("X-WR-CALNAME:%s\n", calendarName)
	ics += fmt.Sprintf("Description:%s: %s\n", "CPE Calendar", calendarName)
	ics += fmt.Sprintf("X-WR-CALDESC:%s: %s\n", "CPE Calendar", calendarName)
	ics += "REFRESH-INTERVAL;VALUE=DURATION:PT1H\n"

	// Define the layout for parsing the datetime with a timezone offset
	const layout = "2006-01-02T15:04:05.000"

	for _, event := range events {

		if event.Favori == nil {
			continue
		}

		// Extract components from regex matches
		location := event.Favori.F2
		summary := event.Favori.F5 + event.Favori.F3
		description := event.Favori.F4

		// Parse the start and end times in the given time zone
		start, err := time.Parse(layout, event.DateDebut)
		if err != nil {
			// Handle parsing error
			fmt.Println("Error parsing start time:", err)
			continue
		}
		end, err := time.Parse(layout, event.DateFin)
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
		ics += fmt.Sprintf("UID:%d\n", event.ID)
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
