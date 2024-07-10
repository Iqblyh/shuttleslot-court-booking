package util

func IsValidFilter(filter string) bool {
	return filter == "daily" || filter == "monthly" || filter == "yearly"
}
