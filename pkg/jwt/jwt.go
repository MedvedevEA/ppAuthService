package jwt

import (
	"crypto/rsa"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	Jti        *uuid.UUID `json:"jti"`
	Sub        *uuid.UUID `json:"sub"`
	DeviceCode string     `json:"device"`
	TokenType  string     `json:"type"`
	jwt.RegisteredClaims
}

func ParseToken(tokenString string, publicKey *rsa.PublicKey) (*TokenClaims, error) {
	tokenJwt, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	tokenClaims, ok := tokenJwt.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("get token claims error")
	}
	return tokenClaims, nil
}
