package infra_interface

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTManager interface {
	Name() string
	Start() error
	Stop() error

	Sign(isRefresh bool, userID string, roles ...string) (string, error)
	Verify(token string) (*JWTClaims, error)
}
