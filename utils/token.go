package utils

import (
	"IMChat/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"time"
)

type claims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

var secretKey = []byte(viper.GetString("JWT.secretKey"))

func CreateToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(time.Hour)
	claim := &claims{
		Id:       user.ID,
		Username: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		panic(err)
	}
	return tokenString, err
}

func ParseToken(token string) (*claims, error) {
	claims := &claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	return claims, err
}
