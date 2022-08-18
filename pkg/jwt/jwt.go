package jwt

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Edbeer/Project/internal/entity"
	"github.com/golang-jwt/jwt"
)

// Manager
type Manager struct {
	signingKey string
}

// JWT Manager constructor
func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

// JWT Claims struct
type Claims struct {
	Email string `json:"email"`
	ID string `json:"id"`
	jwt.StandardClaims
}

// Generate JWT token
func (m *Manager) GenerateJWTToken(user *entity.User) (string, error) {
	claims := &Claims{
		Email: user.Email,
		ID: user.ID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Register the JWT string
	tokenString, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

// Parse access token
func (m *Manager) Parse(accessToken string) (string, error) {
	if accessToken == "" {
		log.Fatal("invalid jwt token")
	}

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", t.Header["alg"])
		}
		secret := []byte(m.signingKey)
		return secret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["id"].(string), nil
}

func (m *Manager) NewRefreshToken() string {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", b)
}