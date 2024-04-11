package server

import (
	"github.com/golang-jwt/jwt/v5"
)

type BearerToken struct {
	UserName string `json:"userName"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
