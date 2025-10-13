package jwtutils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

func CreateJwtToken(claims jwt.Claims, secretKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, err
}

func ValidateJwtToken(tokenStr string, secretKey string) (jwt.MapClaims, error) {
	logrus.Info(tokenStr)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("bad jwt token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logrus.Error("in validate error")
		return nil, fmt.Errorf("bad jwt payload")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, fmt.Errorf("expire time out")
	}
	return claims, nil
}
