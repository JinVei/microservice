package entity

import "github.com/dgrijalva/jwt-go"

type Jwt struct {
	Uid   string `json:"uid"`
	Sid   string `json:"sid"`
	Uname string `json:"uname"`
	jwt.StandardClaims
}
