package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/web"
)

type JWTInterface interface {
	GenerateToken(request web.ResponseLogin) (*string, error)
	GetEmail(token string) (string, error)
	GetUserID(token string) (uint, error)
}

type jwtclaim struct {
	KeyID  string `json:"id"`
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type jsonwebtoken struct {
	config *config.App
}

func NewJWTMiddleware() JWTInterface {
	return &jsonwebtoken{
		config: config.GetConfig(),
	}
}

func (j *jsonwebtoken) GenerateToken(request web.ResponseLogin) (*string, error) {
	claim := &jwtclaim{
		KeyID:  request.KeyID,
		Email:  request.Email,
		UserID: request.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := accessToken.SignedString([]byte(j.config.Server.Secret))
	if err != nil {
		return nil, web.InternalServerError("Cannot generate JWT Token")
	}

	return &token, nil
}

func (j *jsonwebtoken) captureToken(token string) (*jwtclaim, error) {
	tkn, err := jwt.ParseWithClaims(token, &jwtclaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.config.Server.Secret), nil
	})

	if err != nil {
		return nil, web.InternalServerError("Error while parsing Token")
	}

	claim, valid := tkn.Claims.(*jwtclaim)
	if !valid {
		return nil, web.Forbidden("Sorry your token is invalid")
	}

	return claim, nil
}

func (j *jsonwebtoken) GetEmail(token string) (string, error) {
	claim, err := j.captureToken(token)
	if err != nil {
		return "", web.InternalServerError(err.Error())
	}

	return claim.Email, nil
}

func (j *jsonwebtoken) GetUserID(token string) (uint, error) {
	claim, err := j.captureToken(token)
	if err != nil {
		return 0, web.InternalServerError(err.Error())
	}

	return claim.UserID, nil
}
