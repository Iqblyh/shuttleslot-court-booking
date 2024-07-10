package util

func IsValidFilter(filter string) bool {
	return filter == "daily" || filter == "monthly" || filter == "yearly"
}

func IsValidPaymentMethod(method string) bool {
	return method == "mid" || method == "cash"
}
