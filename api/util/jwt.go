package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateToken(jwtkey []byte, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.StandardClaims{
		Subject:   email,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(jwtKey []byte, tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}
	return claims, nil
}

func RefreshJWTToken(jwtKey []byte, oldTokenString string) (string, error) {
	oldClaims, err := ParseToken(jwtKey, oldTokenString)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	newClaims := &jwt.StandardClaims{
		Subject:   oldClaims.Subject,
		ExpiresAt: expirationTime.Unix(),
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err := newToken.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return newTokenString, nil
}

func GetUserIDFromJWT(c *gin.Context, jwtkey []byte) (string, error) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString == "" {
		return "", errors.New("Authorization header is empty")
	}

	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return jwtkey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	return "", errors.New("Invalid token")
}
