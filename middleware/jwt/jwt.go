package jwtHelper

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go-douyin-demo/config"
)

type UserClaims struct {
	jwt.StandardClaims
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

func GenToken(id uint, username string) (string, error) {
	userClaim := &UserClaims{
		UserID:         id,
		Username:       username,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenStr, err := token.SignedString([]byte(config.Cfg.JwtConfig.JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ParseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JwtConfig.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("parse token error:%v", err)
	}
	return userClaim, nil
}
