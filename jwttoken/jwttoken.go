package jwttoken

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		fmt.Println("failed to create token")
		return "", err
	}
	return tokenString, nil
}
func VerifyToken(tokenString string) error {
	privateKeyStr := os.Getenv("SECRET")
	privateKey := []byte(privateKeyStr)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
