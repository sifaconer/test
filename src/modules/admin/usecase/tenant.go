package usecase

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/database/postgres"
	"api-test/src/modules/admin/domain"
	"api-test/src/modules/admin/repository"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/unicode/norm"
)

type Tenant interface {
	CreateTenant(ctx context.Context, tenant domain.DTOTenant) (*domain.DTOTenant, error)
	GetTenantByID(ctx context.Context, id uuid.UUID) (*domain.DTOTenant, error)
	GetTenantByName(ctx context.Context, name string) (*domain.DTOTenant, error)
	UpdateTenant(ctx context.Context, tenant domain.TableTenant) (*domain.TableTenant, error)
	DeleteTenant(ctx context.Context, id uuid.UUID) error
	RegisterAllTenants(ctx context.Context) error
	ListTenants(ctx context.Context) ([]domain.DTOTenant, error)
}

type tenant struct {
	log           common.Logger
	repo          repository.TenantRepository
	tenantManager *common.TenantConnectionManager
	migrations    TenantMigrations
	psql          postgres.Database
	config        *config.Config
	crypto        *encryption
}

func (t *tenant) CreateTenant(ctx context.Context, tenant domain.DTOTenant) (*domain.DTOTenant, error) {
	// Extraer el usuario ID del contexto
	userFromCtx := ctx.Value(t.tenantManager.UserIDKey)
	if userFromCtx == nil {
		t.log.Error(ctx, "Error creating user tenant", "error", "user not found in context")
		return nil, errors.New("user not found to associate with tenants")
	}
	userUUID, ok := userFromCtx.(uuid.UUID)
	if !ok {
		t.log.Error(ctx, "Error creating user tenant", "error", "invalid user id")
		return nil, errors.New("invalid user id to associate with tenants")
	}
	// Generar credenciales
	password, err := t.crypto.GenerateRandomPassword()
	if err != nil {
		return nil, err
	}
	ciphertext, iv, err := t.crypto.Encrypt(password)
	if err != nil {
		return nil, err
	}
	id, _ := uuid.NewRandom()
	dbName := t.generateDBName(tenant.Name, id)

	// Create tenant in database
	newTenant, err := t.repo.CreateTenant(ctx, domain.TableTenant{
		ID:                id,
		Name:              tenant.Name,
		DBName:            dbName,
		DBHost:            t.config.DBConfig.DBHost,
		DBPort:            t.config.DBConfig.DBPort,
		DBUser:            t.generateDBUser(id),
		DBPassword:        ciphertext,
		IsActive:          true,
		IV:                iv,
		Version:           "AES-GCM-256",
		PasswordPlaintext: password,
	})
	if err != nil {
		t.log.Error(ctx, "Error creating tenant", "error", err)
		return nil, err
	}

	// Asociar el tenant con el usuario
	_, err = t.repo.CreateUserTenant(ctx, domain.TableUserTenant{
		ID:       uuid.New(),
		TenantID: newTenant.ID,
		UserID:   userUUID,
	})
	if err != nil {
		t.log.Error(ctx, "Error creating user tenant", "error", err)
		return nil, err
	}

	// Register tenant in connection manager
	pwd, err := t.crypto.Decrypt(ciphertext, iv)
	if err != nil {
		t.log.Error(ctx, "Error decrypting password", "error", err)
		return nil, err
	}
	dsn := t.tenantManager.DSN(common.DSNConfig{
		Host:     newTenant.DBHost,
		Port:     newTenant.DBPort,
		User:     newTenant.DBUser,
		Password: string(pwd),
		Database: newTenant.DBName,
		SSLMode:  t.config.SSLMode,
	})
	if err := t.tenantManager.RegisterTenant(&common.TenantConfig{
		TenantID:         newTenant.ID,
		Name:             newTenant.Name,
		ConnectionString: dsn,
	}, t.psql.Connect); err != nil {
		t.log.Error(ctx, "Error registering tenant", "error", err)
		return nil, err
	}

	// Run migrations
	if err := t.migrations.RunAllMigrations(ctx, newTenant.ID); err != nil {
		t.log.Error(ctx, "Error running migrations", "error", err)
		return nil, err
	}

	result := newTenant.ToDTO()
	return &result, nil
}

func (t *tenant) GetTenantByID(ctx context.Context, id uuid.UUID) (*domain.DTOTenant, error) {
	result, err := t.repo.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}
	dto := result.ToDTO()
	return &dto, nil
}

func (t *tenant) GetTenantByName(ctx context.Context, name string) (*domain.DTOTenant, error) {
	result, err := t.repo.GetTenantByName(ctx, name)
	if err != nil {
		return nil, err
	}
	dto := result.ToDTO()
	return &dto, nil
}

func (t *tenant) UpdateTenant(ctx context.Context, tenant domain.TableTenant) (*domain.TableTenant, error) {
	return t.repo.UpdateTenant(ctx, tenant)
}

func (t *tenant) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	return t.repo.DeleteTenant(ctx, id)
}

func (t *tenant) RegisterAllTenants(ctx context.Context) error {
	tenants, err := t.repo.GetAllTenants(ctx)
	if err != nil {
		t.log.Error(ctx, "Error getting all tenants", "error", err)
		return err
	}

	for _, tenant := range tenants {
		pwd, err := t.crypto.Decrypt(tenant.DBPassword, tenant.IV)
		if err != nil {
			t.log.Error(ctx, "Error decrypting password", "error", err)
			return err
		}
		dsn := t.tenantManager.DSN(common.DSNConfig{
			Host:     tenant.DBHost,
			Port:     tenant.DBPort,
			User:     tenant.DBUser,
			Password: string(pwd),
			Database: tenant.DBName,
			SSLMode:  t.config.SSLMode,
		})
		t.log.Info(ctx, "Registering tenant", "tenant", tenant.Name, "tenant_id", tenant.ID)
		if err := t.tenantManager.RegisterTenant(&common.TenantConfig{
			TenantID:         tenant.ID,
			Name:             tenant.Name,
			ConnectionString: dsn,
		}, t.psql.Connect); err != nil {
			t.log.Error(ctx, "Error registering tenant", "error", err)
			return err
		}
		t.log.Info(ctx, "Tenant registered", "tenant", tenant.Name, "tenant_id", tenant.ID)
	}
	return nil
}

func (t *tenant) ListTenants(ctx context.Context) ([]domain.DTOTenant, error) {
	// Get user id
	userFromCtx := ctx.Value(t.tenantManager.UserIDKey)
	if userFromCtx == nil {
		t.log.Error(ctx, "Error getting user id", "error", "user not found in context")
		return nil, errors.New("user not found in context")
	}
	userID, ok := userFromCtx.(uuid.UUID)
	if !ok {
		t.log.Error(ctx, "Error getting user id", "error", "invalid user id")
		return nil, errors.New("invalid user id")
	}
	// Get user tenants
	tenants, err := t.repo.GetTenantsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var dtos []domain.DTOTenant
	for _, tenant := range tenants {
		dtos = append(dtos, domain.DTOTenant{
			ID:   tenant.Tenant.ID,
			Name: tenant.Tenant.Name,
			UserID: tenant.UserID,
		})
	}
	return dtos, nil
}

func (t *tenant) removeAccentsAndSpecialChars(input string) string {
	normStr := norm.NFD.String(input)
	var sb strings.Builder
	for _, r := range normStr {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			if unicode.Is(unicode.Mn, r) {
				continue
			}
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func (t *tenant) generateDBName(tenantName string, id uuid.UUID) string {
	normalized := t.normalizeName(tenantName)
	shortID := id.String()[:8]
	if len(normalized) > 63 { // PostgreSQL tiene un límite de 63 caracteres
		normalized = normalized[:50]
	}
	dbName := fmt.Sprintf("db_%s_%s", normalized, shortID)
	return dbName
}

func (t *tenant) normalizeName(name string) string {
	norm := t.removeAccentsAndSpecialChars(name)
	name = strings.ToLower(norm)

	reg := regexp.MustCompile(`[^a-z0-9]`)
	name = reg.ReplaceAllString(name, "_")

	reg = regexp.MustCompile(`_+`)
	name = reg.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")

	if len(name) > 0 && unicode.IsDigit(rune(name[0])) {
		name = "t_" + name
	}

	return name
}

func (t *tenant) generateDBUser(id uuid.UUID) string {
	shortID := id.String()[:8]
	dbUser := fmt.Sprintf("kosvi_admin_%s", shortID)
	return dbUser
}

func NewTenant(log common.Logger, repo repository.TenantRepository, migrations TenantMigrations, config *config.Config, tenantManager *common.TenantConnectionManager, psql postgres.Database) Tenant {
	return &tenant{
		log:           log,
		repo:          repo,
		config:        config,
		tenantManager: tenantManager,
		migrations:    migrations,
		psql:          psql,
		crypto:        NewEncryption(log, config),
	}
}

var _ Tenant = (*tenant)(nil)
