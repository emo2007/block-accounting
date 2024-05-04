package jwt

import (
	"fmt"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTInteractor interface {
	NewToken(user models.UserIdentity, duration time.Duration) (string, error)
}

type jwtInteractor struct {
	secret []byte
}

func NewWardenJWT(secret []byte) JWTInteractor {
	return &jwtInteractor{secret}
}

// NewToken creates new JWT token for given user
func (w *jwtInteractor) NewToken(user models.UserIdentity, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id()
	claims["exp"] = time.Now().Add(duration).Unix()

	secret := w.secret

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error sign token. %w", err)
	}

	return tokenString, nil
}
