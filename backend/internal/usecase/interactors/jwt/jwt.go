package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrorInvalidTokenClaims = errors.New("invalid token claims")
	ErrorTokenExpired       = errors.New("token expired")
)

type JWTInteractor interface {
	NewToken(user models.UserIdentity, duration time.Duration) (string, error)
	User(token string) (*models.User, error)
}

type jwtInteractor struct {
	secret          []byte
	usersInteractor users.UsersInteractor
}

func NewWardenJWT(secret []byte, usersInteractor users.UsersInteractor) JWTInteractor {
	return &jwtInteractor{
		secret:          secret,
		usersInteractor: usersInteractor,
	}
}

// NewToken creates new JWT token for given user
func (w *jwtInteractor) NewToken(user models.UserIdentity, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id().String()
	claims["exp"] = time.Now().Add(duration).UnixMilli()

	secret := w.secret

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error sign token. %w", err)
	}

	return tokenString, nil
}

func (w *jwtInteractor) User(tokenStr string) (*models.User, error) {
	claims := make(jwt.MapClaims)

	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return w.secret, nil
	})
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error parse jwt token. %w", err), ErrorInvalidTokenClaims)
	}

	if expDate, ok := claims["exp"].(float64); ok {
		if time.UnixMilli(int64(expDate)).Before(time.Now()) {
			return nil, fmt.Errorf("error token expired. %w", ErrorTokenExpired)
		}
	} else {
		return nil, errors.Join(fmt.Errorf("error parse exp date. %w", err), ErrorInvalidTokenClaims)
	}

	var userIdString string
	var ok bool

	if userIdString, ok = claims["uid"].(string); !ok {
		return nil, ErrorInvalidTokenClaims
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error parse user id. %w", err), ErrorInvalidTokenClaims)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	users, err := w.usersInteractor.Get(ctx, users.GetParams{
		Ids: uuid.UUIDs{userId},
	})
	if err != nil || len(users) == 0 {
		return nil, fmt.Errorf("error fetch user from repository. %w", err)
	}

	return users[0], nil
}
