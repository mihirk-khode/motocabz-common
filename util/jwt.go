package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSecretProvider is an interface for providing JWT secret
type JWTSecretProvider interface {
	GetJWTSecret() string
}

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"sub"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secret []byte
}

// NewJWTManager creates a new JWT manager with a secret provider
func NewJWTManager(secretProvider JWTSecretProvider) *JWTManager {
	secret := secretProvider.GetJWTSecret()
	if secret == "" {
		secret = "tp54XJqd7sb7vw8dQXgRZcHdv3k3+YI7fUgaPdZStY8=" // Default fallback
	}
	return &JWTManager{
		secret: []byte(secret),
	}
}

// NewJWTManagerWithSecret creates a new JWT manager with a direct secret
func NewJWTManagerWithSecret(secret string) *JWTManager {
	if secret == "" {
		secret = "tp54XJqd7sb7vw8dQXgRZcHdv3k3+YI7fUgaPdZStY8=" // Default fallback
	}
	return &JWTManager{
		secret: []byte(secret),
	}
}

// GenerateToken generates a new JWT token
func (j *JWTManager) GenerateToken(userID string, role string, ttl time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ParseToken parses a JWT token string
func (j *JWTManager) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// ValidateToken validates a JWT token and checks expiration
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
