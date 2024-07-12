package util

import "time"

func TimeToString(time time.Time) string {
	stringTime := time.Format("15:04:05")
	return stringTime
}

func DateToString(date time.Time) string {
	stringDate := date.Format("02-01-2006")
	return stringDate
}

func StringToTime(timeString string) time.Time {
	formatedTime, _ := time.Parse("15:04:05", timeString)
	return formatedTime
}

func StringToDate(dateString string) time.Time {
	formatedDate, _ := time.Parse("02-01-2006", dateString)
	return formatedDate
}

func InTimeSpanStart(start, end, checkStart time.Time) bool {
	if checkStart.After(start) || checkStart.Equal(start) {
		if checkStart.Equal(end) || checkStart.After(end) {
			return false
		}

		if checkStart.Before(end) {
			return true
		}

		return true
	}

	return false
}

func InTimeSpanEnd(start, end, checkEnd time.Time) bool {
	if checkEnd.After(start) && (checkEnd.Before(end) || checkEnd.Equal(end)) {
		return true
	}

	return false
}
