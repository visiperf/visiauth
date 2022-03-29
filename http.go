package visiauth

func isStatusCodeError(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}
