package rest

import (
	"errors"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/interface/rest/controllers"
)

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func buildApiError(code int, message string) apiError {
	return apiError{
		Code:    code,
		Message: message,
	}
}

func mapError(err error) apiError {
	switch {
	case errors.Is(err, controllers.ErrorAuthInvalidMnemonic):
		return buildApiError(http.StatusBadRequest, "Invalid Mnemonic")
	default:
		return buildApiError(http.StatusInternalServerError, "Internal Server Error")
	}
}
