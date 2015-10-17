package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Login struct {
	TokenSecretKey string
	TokenExpiredIn int64
}

func NewLogin(secret string, sec int64) *Login {
	return &Login{TokenSecretKey: secret, TokenExpiredIn: sec}
}

func (l *Login) GetToken(user *User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Header["alg"] = "HS256"
	token.Header["typ"] = "JWT"

	token.Claims["id"] = user.Id
	token.Claims["username"] = user.Nickname
	token.Claims["exp"] = time.Now().Add(time.Second * time.Duration(l.TokenExpiredIn)).Unix()

	return token.SignedString([]byte(l.TokenSecretKey))
}
