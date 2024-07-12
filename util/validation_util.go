package util

import (
	"time"
)

func IsValidFilter(filter string) bool {
	return filter == "daily" || filter == "monthly" || filter == "yearly"
}

func IsValidPaymentMethod(method string) bool {
	return method == "mid" || method == "cash"
}

func IsValidDate(dateString string) bool {
	const layout = "02-01-2006"
	_, err := time.Parse(layout, dateString)
	return err == nil
}

func IsValidTime(timeString string) bool {
	const layout = "15:04:05"
	_, err := time.Parse(layout, timeString)
	return err == nil
}
