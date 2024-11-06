package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your_secret_key") // Ensure this is stored securely

// GenerateJWT generates a JWT token with the user's email
func GenerateJWT(email string) (string, error) {
	// Create a new token object with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
		"iat":   time.Now().Unix(),                     // Issued at time
	})

	// Sign the token with the secret
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
