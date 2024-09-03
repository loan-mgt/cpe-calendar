package ical

import (
	"fmt"
	"time"
)

func GenerateICS(events []Event) string {
	ics := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Your Organization//Your Product//EN\n"

	for _, event := range events {
		start, _ := time.Parse(time.RFC3339, event.Start)
		end, _ := time.Parse(time.RFC3339, event.End)

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
