package common

import (
	"api-test/src/config"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TenantConfig struct {
	TenantID         uuid.UUID
	Name             string
	ConnectionString string
}

type DSNConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	Database string
	SSLMode  string
}

type TenantConnectionManager struct {
	mu        sync.RWMutex
	config    *config.Config
	bunDBs    map[uuid.UUID]*bun.DB
	configs   map[uuid.UUID]*TenantConfig
	TenantKey string
	UserIDKey string
}

func NewTenantConnectionManager(config *config.Config) *TenantConnectionManager {
	return &TenantConnectionManager{
		bunDBs:    make(map[uuid.UUID]*bun.DB),
		configs:   make(map[uuid.UUID]*TenantConfig),
		config:    config,
		TenantKey: TenantKey,
		UserIDKey: UserIDKey,
	}
}

func (m *TenantConnectionManager) RegisterTenant(config *TenantConfig, connectFunc func(tenantID uuid.UUID, dsn string) (*bun.DB, error)) error {
	if _, err := m.GetDB(config.TenantID); err == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.configs[config.TenantID] = config
	_, exists := m.bunDBs[config.TenantID]
	if exists {
		return nil
	}

	db, err := connectFunc(config.TenantID, config.ConnectionString)
	if err != nil {
		return err
	}
	m.bunDBs[config.TenantID] = db
	return nil
}

// dsn format
func (m *TenantConnectionManager) DSN(params DSNConfig) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", params.User, params.Password, params.Host, params.Port, params.Database, params.SSLMode)
}

func (m *TenantConnectionManager) GetTenantConfig(tenantID uuid.UUID) (*TenantConfig, error) {
	m.mu.RLock()
	config, exists := m.configs[tenantID]
	m.mu.RUnlock()

	if exists {
		return config, nil
	}

	return nil, fmt.Errorf("no configuration found for tenant: %s", tenantID)
}

func (m *TenantConnectionManager) GetDB(tenantID uuid.UUID) (*bun.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	db, exists := m.bunDBs[tenantID]
	if exists {
		return db, nil
	}

	return nil, fmt.Errorf("no db found for tenant: %s", tenantID)
}

func (m *TenantConnectionManager) GetKosviTenantDB() (*bun.DB, error) {
	return m.GetDB(m.config.TenantID)
}

func (m *TenantConnectionManager) GetDBContext(ctx context.Context) (*bun.DB, error) {
	tenantID, ok := ctx.Value(m.TenantKey).(uuid.UUID)
	if !ok {
		return nil, fmt.Errorf("no tenant found in context")
	}
	return m.GetDB(tenantID)
}

func (m *TenantConnectionManager) RemoveTenant(tenantID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	db, exists := m.bunDBs[tenantID]
	if exists {
		if err := db.Close(); err != nil {
			return err
		}
		delete(m.bunDBs, tenantID)
	}

	delete(m.configs, tenantID)
	return nil
}

func (m *TenantConnectionManager) CloseAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, db := range m.bunDBs {
		db.Close()
	}

	m.bunDBs = make(map[uuid.UUID]*bun.DB)
	return nil
}
