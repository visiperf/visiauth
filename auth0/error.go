package auth0

func isStatusCodeError(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}
