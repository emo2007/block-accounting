package rest

import (
	"errors"
	"net/http"

	"github.com/emochka2007/block-accounting/internal/interface/rest/controllers"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/jwt"
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
	// auth controller errors
	case errors.Is(err, controllers.ErrorAuthInvalidMnemonic):
		return buildApiError(http.StatusBadRequest, "Invalid Mnemonic")
	case errors.Is(err, controllers.ErrorTokenRequired):
		return buildApiError(http.StatusUnauthorized, "Token Required")

	// jwt-related errors
	case errors.Is(err, jwt.ErrorTokenExpired):
		return buildApiError(http.StatusUnauthorized, "Token Expired")
	case errors.Is(err, jwt.ErrorInvalidTokenClaims):
		return buildApiError(http.StatusUnauthorized, "Invalid Token")
	default:
		return buildApiError(http.StatusInternalServerError, "Internal Server Error")
	}
}
