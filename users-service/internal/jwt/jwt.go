package jwt

import (
	"time"

	"github.com/cfex/microservices-in-go/users-service/cmd/config"
	"github.com/golang-jwt/jwt/v5"
)

var Jwt *jwtUtils

type Payload struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type JwtClaims struct {
	Payload Payload `json:"payload"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
}

type jwtUtils struct {
	accessKey  []byte
	ttl        time.Duration
	refreshKey []byte
	refreshTtl time.Duration
}

func NewJwt(c *config.Config) {
	accessKey := []byte(c.Jwt.SecretKey)
	refreshKey := []byte(c.Jwt.SecretKey)
	ttl := c.Jwt.TTL
	refreshTtl := c.Jwt.RefreshTTL

	Jwt = &jwtUtils{
		accessKey:  accessKey,
		ttl:        ttl,
		refreshKey: refreshKey,
		refreshTtl: refreshTtl,
	}
}

func (ju *jwtUtils) GenerateTokens(payload Payload) (string, string, error) {
	accessClaims := &JwtClaims{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ju.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims).SignedString(ju.accessKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := &RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ju.refreshTtl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims).SignedString(ju.refreshKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (ju *jwtUtils) ParseToken(token string) (*JwtClaims, error) {
	parse, err := jwt.ParseWithClaims(token, &JwtClaims{}, func(t *jwt.Token) (any, error) {
		return ju.accessKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parse.Claims.(*JwtClaims)
	if !ok || !parse.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
