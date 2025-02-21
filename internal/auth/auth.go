package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v4"
)

var SecretKey = "secret-yy-xz"

func GenerateJWT(userID uint) (string, error) {
	claims:= jwt.MapClaims{
		"user_id": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SecretKey))
}


