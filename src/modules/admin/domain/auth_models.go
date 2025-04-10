package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	UserID  uuid.UUID   `json:"user_id"`
	Tenants []uuid.UUID `json:"tenants,omitempty"`
	Type    TokenType   `json:"type"`
	jwt.RegisteredClaims
}
