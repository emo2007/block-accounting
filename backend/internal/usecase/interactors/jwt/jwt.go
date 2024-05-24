package jwt

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/emochka2007/block-accounting/internal/pkg/models"
	"github.com/emochka2007/block-accounting/internal/usecase/interactors/users"
	"github.com/emochka2007/block-accounting/internal/usecase/repository/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrorInvalidTokenClaims = errors.New("invalid token claims")
	ErrorTokenExpired       = errors.New("token expired")
)

type JWTInteractor interface {
	NewToken(user models.UserIdentity, duration time.Duration, remoteAddr string) (AccessToken, error)
	User(token string) (*models.User, error)
	RefreshToken(ctx context.Context, token string, rToken string) (AccessToken, error)
}

type jwtInteractor struct {
	secret          []byte
	usersInteractor users.UsersInteractor
	authRepository  auth.Repository
}

func NewJWT(
	secret []byte,
	usersInteractor users.UsersInteractor,
	authRepository auth.Repository,
) JWTInteractor {
	return &jwtInteractor{
		secret:          secret,
		usersInteractor: usersInteractor,
		authRepository:  authRepository,
	}
}

type AccessToken struct {
	Token     string
	ExpiredAt time.Time

	RefreshToken string
	RTExpiredAt  time.Time
}

// NewToken creates new JWT token for given user
func (w *jwtInteractor) NewToken(user models.UserIdentity, duration time.Duration, remoteAddr string) (AccessToken, error) {
	tokens, err := w.newTokens(user.Id(), duration)
	if err != nil {
		return AccessToken{}, fmt.Errorf("error create new tokens. %w", err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	if err := w.authRepository.AddToken(ctx, auth.AddTokenParams{
		UserId:                user.Id(),
		Token:                 tokens.Token,
		TokenExpiredAt:        tokens.ExpiredAt,
		RefreshToken:          tokens.RefreshToken,
		RefreshTokenExpiredAt: tokens.RTExpiredAt,
		CreatedAt:             time.Now(),
	}); err != nil {
		return AccessToken{}, fmt.Errorf("error save tokens into repository. %w", err)
	}

	return tokens, nil
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

	tokens, err := w.authRepository.GetTokens(ctx, auth.GetTokenParams{
		UserId: userId,
		Token:  tokenStr,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetch token from repository. %w", err)
	}

	if tokens.TokenExpiredAt.Before(time.Now()) {
		return nil, fmt.Errorf("error token expired. %w", ErrorTokenExpired)
	}

	if tokens.UserId != userId {
		return nil, errors.Join(fmt.Errorf("error invalid user id. %w", err), ErrorInvalidTokenClaims)
	}

	users, err := w.usersInteractor.Get(ctx, users.GetParams{
		Ids: uuid.UUIDs{tokens.UserId},
	})
	if err != nil || len(users) == 0 {
		return nil, fmt.Errorf("error fetch user from repository. %w", err)
	}

	return users[0], nil
}

func (w *jwtInteractor) RefreshToken(ctx context.Context, token string, rToken string) (AccessToken, error) {
	claims := make(jwt.MapClaims)

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return w.secret, nil
	})
	if err != nil {
		return AccessToken{}, errors.Join(fmt.Errorf("error parse jwt token. %w", err), ErrorInvalidTokenClaims)
	}

	var userIdString string
	var ok bool

	if userIdString, ok = claims["uid"].(string); !ok {
		return AccessToken{}, ErrorInvalidTokenClaims
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return AccessToken{}, errors.Join(fmt.Errorf("error parse user id. %w", err), ErrorInvalidTokenClaims)
	}

	_, err = jwt.ParseWithClaims(rToken, claims, func(t *jwt.Token) (interface{}, error) {
		return w.secret, nil
	})
	if err != nil {
		return AccessToken{}, errors.Join(fmt.Errorf("error parse refresh jwt token. %w", err), ErrorInvalidTokenClaims)
	}

	if expDate, ok := claims["exp"].(float64); ok {
		if time.UnixMilli(int64(expDate)).Before(time.Now()) {
			return AccessToken{}, fmt.Errorf("error refresh token expired. %w", ErrorTokenExpired)
		}
	} else {
		return AccessToken{}, errors.Join(fmt.Errorf("error parse exp date. %w", err), ErrorInvalidTokenClaims)
	}

	if userIdString, ok = claims["uid"].(string); !ok {
		return AccessToken{}, ErrorInvalidTokenClaims
	}

	rTokenUserId, err := uuid.Parse(userIdString)
	if err != nil {
		return AccessToken{}, errors.Join(
			fmt.Errorf("error parse user id from refresh token. %w", err),
			ErrorInvalidTokenClaims,
		)
	}

	if userId != rTokenUserId {
		return AccessToken{}, fmt.Errorf("error user ids corrupted. %w", ErrorInvalidTokenClaims)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()

	tokens, err := w.authRepository.GetTokens(ctx, auth.GetTokenParams{
		UserId:       userId,
		Token:        token,
		RefreshToken: rToken,
	})
	if err != nil {
		return AccessToken{}, fmt.Errorf("error fetch token from repository. %w", err)
	}

	if tokens.RefreshTokenExpiredAt.Before(time.Now()) {
		return AccessToken{}, fmt.Errorf("error token expired. %w", ErrorTokenExpired)
	}

	rtHash := sha512.New()
	rtHash.Write([]byte(tokens.Token))

	rtHashStringValid := base64.StdEncoding.EncodeToString(rtHash.Sum(nil))

	rtHashRaw, ok := claims["rt_hash"]
	if !ok {
		return AccessToken{}, fmt.Errorf("error refresh token claims corrupted. %w", ErrorInvalidTokenClaims)
	}

	rtHashString, ok := rtHashRaw.(string)
	if !ok {
		return AccessToken{}, fmt.Errorf("error refresh token claims corrupted. %w", ErrorInvalidTokenClaims)
	}

	if rtHashString != rtHashStringValid {
		return AccessToken{}, fmt.Errorf("error refresh token hash corrupted. %w", ErrorInvalidTokenClaims)
	}

	newTokens, err := w.newTokens(userId, 24*time.Hour)
	if err != nil {
		return AccessToken{}, fmt.Errorf("error create new tokens. %w", err)
	}

	if err = w.authRepository.RefreshToken(ctx, auth.RefreshTokenParams{
		UserId:                userId,
		OldToken:              token,
		Token:                 newTokens.Token,
		TokenExpiredAt:        newTokens.ExpiredAt,
		OldRefreshToken:       rToken,
		RefreshToken:          newTokens.RefreshToken,
		RefreshTokenExpiredAt: newTokens.RTExpiredAt,
	}); err != nil {
		return AccessToken{}, fmt.Errorf("error update tokens. %w", err)
	}

	return newTokens, nil
}

func (w *jwtInteractor) newTokens(userId uuid.UUID, duration time.Duration) (AccessToken, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	expAt := time.Now().Add(duration)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = userId.String()
	claims["exp"] = expAt.UnixMilli()

	secret := w.secret

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return AccessToken{}, fmt.Errorf("error sign token. %w", err)
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)

	rtHash := sha512.New()

	rtHash.Write([]byte(tokenString))

	rtExpAt := expAt.Add(time.Hour * 24 * 5)

	claims = refreshToken.Claims.(jwt.MapClaims)
	claims["uid"] = userId.String()
	claims["exp"] = rtExpAt.UnixMilli()
	claims["rt_hash"] = base64.StdEncoding.EncodeToString(rtHash.Sum(nil))

	rtokenString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return AccessToken{}, fmt.Errorf("error sign refresh token. %w", err)
	}

	return AccessToken{
		Token:        tokenString,
		ExpiredAt:    expAt,
		RefreshToken: rtokenString,
		RTExpiredAt:  rtExpAt,
	}, nil
}
