package utils

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(userID uint, email, role string) (string, error) {
	secret := GetJWTSecret()
	expiry := time.Now().Add(GetJWTExpiry()).Unix()

	claims := JwtClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseTokenString(tokenStr string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(GetJWTSecret()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
