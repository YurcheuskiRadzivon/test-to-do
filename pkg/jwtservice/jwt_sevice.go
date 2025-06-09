package jwtservice

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

const (
	ParamID             = "id"
	HeaderAuthorization = "Authorization"
	StatusInvalidToken  = "INVALID_TOKEN"
)

var (
	ErrInvalidToken = errors.New(StatusInvalidToken)
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
		return "", err
	}
	return t, nil
}

func (jwts *JWTService) GetUserID(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwts.secretKey, nil
	})

	if err != nil {
		return -1, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(float64)
		return int(userID), nil
	}

	return -1, ErrInvalidToken
}
