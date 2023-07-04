package common

import (
	"time"

	"github.com/dgrijalva/jwt-go" // https://www.youtube.com/watch?v=-Eei8eik1Io (from NerdCademyDev channel)
)

var JWTSecret = []byte("SuperSecret")

func GenerateJWT(id uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) //si
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}
