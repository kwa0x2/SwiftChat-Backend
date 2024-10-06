package utils

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var secretKey []byte

func GenerateToken(jwtClaims jwt.MapClaims) (string, error) {
	secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	claims := jwtClaims

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func GetClaims(tokenString string) (map[string]interface{}, error) {
	err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	_, parseErr := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if parseErr != nil {
		return nil, parseErr
	}

	return claims, nil
}

type UserInfo struct {
	Email string
	ID    string
	Name  string
}

func GetUserSessionInfo(ctx *gin.Context) (*UserInfo, error) {
	session := sessions.Default(ctx)

	email := session.Get("email")
	id := session.Get("id")
	name := session.Get("name")

	if email == nil || id == nil || name == nil {
		return nil, errors.New("user information not found in session")
	}

	return &UserInfo{
		Email: email.(string),
		ID:    id.(string),
		Name:  name.(string),
	}, nil
}
