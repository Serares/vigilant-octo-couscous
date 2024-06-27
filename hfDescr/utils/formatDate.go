package utils

import (
	"fmt"
	"time"
)

// Function to add the suffix to the day of the month
func dayWithSuffix(day int) string {
	if day >= 11 && day <= 13 {
		return fmt.Sprintf("%dth", day)
	}
	switch day % 10 {
	case 1:
		return fmt.Sprintf("%dst", day)
	case 2:
		return fmt.Sprintf("%dnd", day)
	case 3:
		return fmt.Sprintf("%drd", day)
	default:
		return fmt.Sprintf("%dth", day)
	}
}

var FormatDate = func(date time.Time) string {
	// Format date
	return fmt.Sprintf("%s %s, %d at %02d:%02d %s",
		date.Month().String(),
		dayWithSuffix(date.Day()),
		date.Year(),
		date.Hour(),
		date.Minute(),
		date.Format("MST"))
}
