package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

// Claims JWT 载荷
type Claims struct {
	UserID   uint64 `json:"sub"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwtlib.RegisteredClaims
}

// Manager JWT 签发与解析
type Manager struct {
	secret              []byte
	expireHours         int
	rememberExpireHours int
}

func NewManager(secret string, expireHours, rememberExpireHours int) *Manager {
	return &Manager{
		secret:              []byte(secret),
		expireHours:         expireHours,
		rememberExpireHours: rememberExpireHours,
	}
}

// Generate 签发 Token
func (m *Manager) Generate(userID uint64, username, role string, remember bool) (string, int64, error) {
	hours := m.expireHours
	if remember {
		hours = m.rememberExpireHours
	}
	expiresAt := time.Now().Add(time.Duration(hours) * time.Hour)

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(expiresAt),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", 0, err
	}
	return signed, int64(hours * 3600), nil
}

// Parse 解析 Token
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtlib.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
