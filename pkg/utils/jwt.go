package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joobers/backend/pkg/models"
	"github.com/joobers/backend/pkg/models/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey = os.Getenv("JWT_KEY")

func GenerateJwt(id primitive.ObjectID) (string, error) {
	expirationTime := time.Now().Add(336 * time.Hour) // 2 weeks
	claims := &models.JwtClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJwt(token string) bool {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return false
	}

	return true
}

func ParseJwt(token string) (*models.JwtClaims, error) {
	claims := &models.JwtClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return nil, errors.InvalidToken
	}

	return claims, nil
}

func RefreshJwt(token string) (string, error) {
	claims, err := ParseJwt(token)
	if err != nil {
		return "", err
	}

	newToken, err := GenerateJwt(claims.ID)
	if err != nil {
		return "", errors.InternalServerError
	}

	return newToken, nil
}
