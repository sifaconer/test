package usecase

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/admin/domain"
	"api-test/src/modules/admin/repository"
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type Auth interface {
	Login(ctx context.Context, model domain.DTOUserDirectory) (*domain.DTOAuth, error)
	Register(ctx context.Context, model domain.DTOUserDirectory) (*domain.DTOAuth, error)
	Logout(ctx context.Context, token *domain.DTOAuth) error
	Refresh(ctx context.Context, token *domain.DTOAuth) (*domain.DTOAuth, error)
}

type auth struct {
	log         common.Logger
	config      *config.Config
	tenant      *common.TenantConnectionManager
	repo        repository.UserDirectoryRepository
	repoTenant  repository.TenantRepository
	argonParams *params
}

// Login implements Auth.
func (a *auth) Login(ctx context.Context, model domain.DTOUserDirectory) (*domain.DTOAuth, error) {
	// Buscar el usuario por email
	userDirectory, err := a.repo.GetUserDirectoryByEmail(ctx, model.Email)
	if err != nil {
		return nil, err
	}

	if userDirectory == nil {
		return nil, common.NotFoundError("user not found")
	}

	// Verificar la contraseña
	match, err := a.comparePasswordAndHash(model.Password, userDirectory.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, common.UnauthorizedError("invalid credentials")
	}

	// Generar el token
	token, err := common.GenerateJWT(ctx, a.config, userDirectory.ToDTO())
	if err != nil {
		return nil, err
	}

	return &domain.DTOAuth{
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

// Logout implements Auth.
func (a *auth) Logout(ctx context.Context, token *domain.DTOAuth) error {
	panic("unimplemented")
}

// Refresh implements Auth.
func (a *auth) Refresh(ctx context.Context, token *domain.DTOAuth) (*domain.DTOAuth, error) {
	// Validate refresh token
	claims, err := common.ValidateJWT(ctx, token.RefreshToken, a.config)
	if err != nil {
		return nil, fmt.Errorf("refresh token inválido: %w", err)
	}
	if claims.Type != domain.TokenTypeRefresh {
		return nil, errors.New("token no es de tipo refresh")
	}

	// Get user
	user, err := a.repo.GetUserDirectoryByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	return common.GenerateJWT(ctx, a.config, user.ToDTO())
}

// Register implements Auth.
func (a *auth) Register(ctx context.Context, model domain.DTOUserDirectory) (*domain.DTOAuth, error) {

	// Generar el hash de la contraseña
	hashedPassword, err := a.generateFromPassword(model.Password, a.argonParams)
	if err != nil {
		return nil, err
	}

	// Crear el usuario
	model.Password = hashedPassword
	model.ID = uuid.New()
	model.IsActive = true
	model.CreationDate = time.Now()
	result, err := a.repo.CreateUserDirectory(ctx, model.ToTable())
	if err != nil {
		return nil, err
	}

	// Consultar tenants asociados al usuario
	tenants, err := a.repo.GetTenantsByUser(ctx, result.ID)
	if err != nil {
		return nil, err
	}
	for _, tenant := range tenants {
		result.UserTenants = append(result.UserTenants, tenant)
	}

	// Generar el token
	token, err := common.GenerateJWT(ctx, a.config, result.ToDTO())
	if err != nil {
		return nil, err
	}

	return &domain.DTOAuth{
		Token:        token.Token,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
	}, nil
}

func (a *auth) generateFromPassword(password string, p *params) (encodedHash string, err error) {
	salt, err := a.generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (a *auth) generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (a *auth) comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := a.decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func (a *auth) decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("invalid hash")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version")
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func NewAuth(log common.Logger, config *config.Config, tenant *common.TenantConnectionManager, repo repository.UserDirectoryRepository) Auth {
	return &auth{
		log:    log,
		config: config,
		tenant: tenant,
		repo:   repo,
		argonParams: &params{
			memory:      64 * 1024, // 64MB
			iterations:  3,
			parallelism: 2,
			saltLength:  16,
			keyLength:   32,
		},
	}
}

var _ Auth = (*auth)(nil)
