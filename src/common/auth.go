package common

import (
	"api-test/src/config"
	"api-test/src/modules/admin/domain"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Genera Access y Refresh JWT con expiraciones separadas
func GenerateJWT(ctx context.Context, config *config.Config, user domain.DTOUserDirectory) (*domain.DTOAuth, error) {
	privKey, err := loadECPrivateKey(config)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	// Access Token
	tenants := []uuid.UUID{}
	for _, tenant := range user.UserTenants {
		tenants = append(tenants, tenant.TenantID)
	}

	expired := now.Add(time.Duration(config.JWT.TTL) * time.Second)
	accessClaims := domain.Claims{
		UserID:  user.ID,
		Tenants: tenants,
		Type:    domain.TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expired),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
			Issuer:    "KOSVI",
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, accessClaims).SignedString(privKey)
	if err != nil {
		return nil, fmt.Errorf("error al firmar access token: %w", err)
	}

	// Refresh Token
	expiredRefresh := now.Add(time.Duration(config.JWT.TTL) * time.Second * 2) // 2x TTL
	refreshClaims := domain.Claims{
		UserID: user.ID,
		Type:   domain.TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredRefresh),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
			Issuer:    "KOSVI",
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).SignedString(privKey)
	if err != nil {
		return nil, fmt.Errorf("error al firmar refresh token: %w", err)
	}

	return &domain.DTOAuth{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expired,
	}, nil
}

func RefreshJWT(ctx context.Context, config *config.Config, auth *domain.DTOAuth, user domain.DTOUserDirectory) (*domain.DTOAuth, error) {
	if auth == nil {
		return nil, errors.New("token inválido")
	}
	// Validate access token
	// claims, err := ValidateJWT(ctx, auth.Token, config)
	// if err != nil {
	// 	return nil, fmt.Errorf("access token inválido: %w", err)
	// }
	// if claims.Type != domain.TokenTypeAccess {
	// 	return nil, errors.New("token no es de tipo access")
	// }
	// Validate refresh token
	claims, err := ValidateJWT(ctx, auth.RefreshToken, config)
	if err != nil {
		return nil, fmt.Errorf("refresh token inválido: %w", err)
	}
	if claims.Type != domain.TokenTypeRefresh {
		return nil, errors.New("token no es de tipo refresh")
	}

	// Generate new tokens
	return GenerateJWT(ctx, config, user)
}

func ValidateJWT(ctx context.Context, tokenString string, config *config.Config) (*domain.Claims, error) {
	publicKey, err := loadECPublicKey(config)
	if err != nil {
		return nil, fmt.Errorf("clave pública inválida: %w", err)
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodES256.Alg()}))
	token, err := parser.ParseWithClaims(tokenString, &domain.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("token mal formado")
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expirado")
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token no válido aún")
		}
		if errors.Is(err, jwt.ErrTokenInvalidAudience) {
			return nil, errors.New("token inválido por audiencia")
		}
		if errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			return nil, errors.New("token inválido por emisor")
		}
		if errors.Is(err, jwt.ErrTokenInvalidSubject) {
			return nil, errors.New("token inválido por sujeto")
		}
		if errors.Is(err, jwt.ErrTokenInvalidId) {
			return nil, errors.New("token inválido por id")
		}
		if errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return nil, errors.New("token inválido por claims")
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("token inválido por firma")
		}
		return nil, fmt.Errorf("error al parsear el token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok {
		return nil, errors.New("los claims no son del tipo esperado")
	}

	return claims, nil
}

func ValidateTenant(ctx context.Context, config *config.Config, tenantID uuid.UUID, tokenString string) error {
	claims, err := ValidateJWT(ctx, tokenString, config)
	if err != nil {
		return err
	}
	if !slices.Contains(claims.Tenants, tenantID) {
		return errors.New("tenant no autorizado")
	}
	return nil
}

// Carga la clave privada EC desde env en base64
func loadECPrivateKey(config *config.Config) (*ecdsa.PrivateKey, error) {
	b64 := config.JWT.ECPrivateKeyBase64
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("error decodificando clave privada base64: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("error decodificando PEM de clave privada")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("clave privada inválida: %w", err)
	}
	return key, nil
}

// Carga la clave pública EC desde env en base64
func loadECPublicKey(config *config.Config) (*ecdsa.PublicKey, error) {
	b64 := config.JWT.ECPublicKeyBase64
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("error decodificando clave pública base64: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("error decodificando PEM de clave pública")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("clave pública inválida: %w", err)
	}
	ecPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("clave pública no es EC")
	}
	return ecPub, nil
}
