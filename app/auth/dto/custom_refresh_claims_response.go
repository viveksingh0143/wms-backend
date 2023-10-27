package dto

import "github.com/golang-jwt/jwt/v5"

type CustomRefreshClaims struct {
	jwt.RegisteredClaims
	ExpireLong bool `json:"expireLong"`
}
