package models

import (
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
