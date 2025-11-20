package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Claims represents the JWT claims structure
type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID primitive.ObjectID, email string, secret string, duration time.Duration) (string, error) {
	expiresAt := time.Now().Add(duration).Unix()

	claims := jwt.MapClaims{
		"user_id": userID.Hex(),
		"email":   email,
		"exp":     expiresAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string, secret string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Extract claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user_id in claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email in claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid exp in claims")
	}

	return &Claims{
		UserID:    userID,
		Email:     email,
		ExpiresAt: int64(exp),
	}, nil
}
