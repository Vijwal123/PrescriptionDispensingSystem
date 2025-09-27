package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var Jwtsecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(id int, role string) (string, error) {

	claims := jwt.MapClaims{
		"user_id": id,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Jwtsecret)

}

func ExtractSecertKey(token *jwt.Token) (interface{}, error) {
	return Jwtsecret, nil
}
