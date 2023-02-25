package app

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	UserId    uint   `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	jwt.StandardClaims
}
