package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"assetManager/internal/models"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents JWT claims
type Claims struct {
	UserID   int64  `json:"UserID"`
	Username string `json:"Username"`
	jwt.RegisteredClaims
}

// JWTService handles JWT operations
type JWTService struct {
	secret      []byte
	expiryHours int
}

// NewJWTService creates a new JWT service
func NewJWTService(secret string, expiryHours int) *JWTService {
	return &JWTService{
		secret:      []byte(secret),
		expiryHours: expiryHours,
	}
}

// GenerateToken creates a new JWT token for a user
func (s *JWTService) GenerateToken(user *models.User, remember bool) (string, int64, error) {
	expiry := time.Duration(s.expiryHours) * time.Hour
	if remember {
		expiry = 30 * 24 * time.Hour // 30 days if remember is checked
	}

	expiresAt := time.Now().Add(expiry)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
