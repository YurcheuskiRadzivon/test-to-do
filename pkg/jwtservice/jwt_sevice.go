package jwtservice

import (
	"errors"
	"log"

	"github.com/golang-jwt/jwt"
)

const (
	ParamID                  = "id"
	HeaderAuthorization      = "Authorization"
	ErrInvalidOrExpiredToken = "INVALID_OR_EXPIRED_TOKEN"
	ErrInvalidCreateToken    = "INVALID_CREATE_TOKEN"
)

type JWTService struct {
	secretKey []byte
}

func New(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
	}
}

func (jwts *JWTService) GetJwtSecretKey() []byte {
	return jwts.secretKey
}

func (jwts *JWTService) CreateToken(payload jwt.MapClaims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwts.secretKey)
	if err != nil {
		log.Printf("Failed to create token: %v", err)
		return "", errors.New(ErrInvalidCreateToken)
	}
	return t, nil
}

func (jwts *JWTService) GetUserID(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwts.secretKey, nil
	})
	if err != nil {
		log.Printf("Failed parse token: %v", err)
		return -1, errors.New(ErrInvalidOrExpiredToken)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(float64)
		return int(userID), nil
	}

	log.Printf("Failed to parse token: %v", err)
	return -1, errors.New(ErrInvalidOrExpiredToken)
}
