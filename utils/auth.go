package utils

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// secretKey is used for signing and verifying JWT tokens.
var secretKey []byte

// region "GenerateToken" creates a new JWT token with the given claims.
func GenerateToken(jwtClaims jwt.MapClaims) (string, error) {
	// Load the JWT secret key from environment variables.
	secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	// Create a new token with the provided claims.
	claims := jwtClaims

	// Create a new JWT token using HS256 signing method.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and return the token string.
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err // Return an error if signing fails.
	}

	return tokenString, nil // Return the signed token.
}

// endregion

// region "VerifyToken" checks if the provided JWT token is valid.
func VerifyToken(tokenString string) error {
	// Parse the token and use the secret key for verification.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil // Return the secret key for signing method verification.
	})

	if err != nil {
		return err // Return an error if parsing fails.
	}

	if !token.Valid {
		return errors.New("invalid token") // Return an error if the token is not valid.
	}

	return nil // Return nil if the token is valid.
}

// endregion

// region "GetClaims" extracts claims from the provided JWT token string.
func GetClaims(tokenString string) (map[string]interface{}, error) {
	// Verify the token before extracting claims.
	err := VerifyToken(tokenString)
	if err != nil {
		return nil, err // Return nil if token verification fails.
	}

	claims := jwt.MapClaims{} // Create a map to hold claims.
	// Parse the token with claims.
	_, parseErr := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil // Return the secret key for signature verification.
	})

	if parseErr != nil {
		return nil, parseErr // Return an error if parsing fails.
	}

	return claims, nil // Return the extracted claims.
}

// endregion

// region UserInfo holds user session information.
type UserInfo struct {
	Email string // User email
	ID    string // User ID
	Name  string // User name
}

// endregion

// region "GetUserSessionInfo" retrieves user information from the session.
func GetUserSessionInfo(ctx *gin.Context) (*UserInfo, error) {
	session := sessions.Default(ctx) // Get the session from the context.

	// Retrieve user information from the session.
	email := session.Get("email")
	id := session.Get("id")
	name := session.Get("name")

	// Check if any user information is missing.
	if email == nil || id == nil || name == nil {
		return nil, errors.New("user information not found in session") // Return an error if any information is missing.
	}

	// Return user information encapsulated in the UserInfo struct.
	return &UserInfo{
		Email: email.(string), // Type assertion to string
		ID:    id.(string),    // Type assertion to string
		Name:  name.(string),  // Type assertion to string
	}, nil
}

// endregion
