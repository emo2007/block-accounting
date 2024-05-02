package rest

import "net/http"

type apiError struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func buildApiError(code int, message string) apiError {
	return apiError{
		Code:  code,
		Error: message,
	}
}

func mapError(_ error) apiError {
	// todo map typed errors
	switch {
	default:
		return buildApiError(http.StatusInternalServerError, "Internal Server Error")
	}
}
