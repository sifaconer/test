package usecase

import (
	"api-test/src/common"
	"api-test/src/config"
	"context"
	"database/sql"
	"embed"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/uptrace/bun"
)

type TenantMigrations interface {
	RunAdminMigrations(ctx context.Context) error
	RunAllMigrations(ctx context.Context, tenantID uuid.UUID) error
	RunMigration(ctx context.Context, tenantID uuid.UUID, migrationID int64) error
	RollbackAllMigrations(ctx context.Context, tenantID uuid.UUID) error
	RollbackMigration(ctx context.Context, tenantID uuid.UUID, migrationID int64) error
}

type tenantMigrations struct {
	log                common.Logger
	config             *config.Config
	tenantManager      *common.TenantConnectionManager
	adminMigrationsFS  embed.FS
	tenantMigrationsFS embed.FS
	commonMigrationsFS embed.FS
}

// RollbackAllMigrations implements TenantMigrations.
func (t *tenantMigrations) RollbackAllMigrations(ctx context.Context, tenantID uuid.UUID) error {
db, err := t.getDB(ctx, tenantID)
	if err != nil {
		return err
	}

	// run tenant migrations
	if err := t.down(ctx, db.DB, t.tenantMigrationsFS, "migrations/tenants"); err != nil {
		return err
	}

	return nil
}

// RollbackMigration implements TenantMigrations.
func (t *tenantMigrations) RollbackMigration(ctx context.Context, tenantID uuid.UUID, migrationID int64) error {
	db, err := t.getDB(ctx, tenantID)
	if err != nil {
		return err
	}

	if err := t.downTo(ctx, db.DB, t.tenantMigrationsFS, "migrations/tenants", migrationID); err != nil {
		return err
	}

	return nil
}

// RunAllMigrations implements TenantMigrations.
func (t *tenantMigrations) RunAllMigrations(ctx context.Context, tenantID uuid.UUID) error {
	db, err := t.getDB(ctx, tenantID)
	if err != nil {
		return err
	}

	// run common migrations
	if err := t.up(ctx, db.DB, t.commonMigrationsFS, "migrations/common"); err != nil {
		return err
	}

	// run tenant migrations
	if err := t.up(ctx, db.DB, t.tenantMigrationsFS, "migrations/tenants"); err != nil {
		return err
	}

	return nil
}

// RunMigration implements TenantMigrations.
func (t *tenantMigrations) RunMigration(ctx context.Context, tenantID uuid.UUID, migrationID int64) error {
	db, err := t.getDB(ctx, tenantID)
	if err != nil {
		return err
	}

	if err := t.upTo(ctx, db.DB, t.tenantMigrationsFS, "migrations/tenants", migrationID); err != nil {
		return err
	}

	return nil
}

// RunAdminMigrations implements TenantMigrations.
func (t *tenantMigrations) RunAdminMigrations(ctx context.Context) error {
	db, err := t.tenantManager.GetKosviTenantDB()
	if err != nil {
		return err
	}

	// run common migrations
	t.log.Info(ctx, "Running Common Migrations")
	if err := t.up(ctx, db.DB, t.commonMigrationsFS, "migrations/common"); err != nil {
		return err
	}

	// run admin migrations
	t.log.Info(ctx, "Running Admin Migrations")
	if err := t.up(ctx, db.DB, t.adminMigrationsFS, "migrations/admin"); err != nil {
		return err
	}

	return nil
}


func (t *tenantMigrations) getDB(ctx context.Context, tenantID uuid.UUID) (*bun.DB, error) {
	var db *bun.DB
	if tenantID != uuid.Nil {
		dbTenant, err := t.tenantManager.GetDB(tenantID)
		if err != nil {
			return nil, err
		}
		db = dbTenant
	} else {
		dbCtx, err := t.tenantManager.GetDBContext(ctx)
		if err != nil {
			return nil, err
		}
		db = dbCtx
	}
	return db, nil
}

func (t *tenantMigrations) up(ctx context.Context, db *sql.DB, embedMigrations embed.FS, folder string) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		t.log.Error(ctx, "Error setting dialect", "error", err)
		return err
	}

	if err := goose.Up(db, folder); err != nil {
		t.log.Error(ctx, "Error running migrations", "error", err)
		return err
	}

	return nil
}

func (t *tenantMigrations) upTo(ctx context.Context, db *sql.DB, embedMigrations embed.FS, folder string, version int64) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		t.log.Error(ctx, "Error setting dialect", "error", err)
		return err
	}

	if err := goose.UpTo(db, folder, version); err != nil {
		t.log.Error(ctx, "Error running migrations", "error", err)
		return err
	}

	return nil
}

func (t *tenantMigrations) down(ctx context.Context, db *sql.DB, embedMigrations embed.FS, folder string) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		t.log.Error(ctx, "Error setting dialect", "error", err)
		return err
	}

	if err := goose.Down(db, folder); err != nil {
		t.log.Error(ctx, "Error running migrations", "error", err)
		return err
	}

	return nil
}

func (t *tenantMigrations) downTo(ctx context.Context, db *sql.DB, embedMigrations embed.FS, folder string, version int64) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		t.log.Error(ctx, "Error setting dialect", "error", err)
		return err
	}

	if err := goose.DownTo(db, folder, version); err != nil {
		t.log.Error(ctx, "Error running migrations", "error", err)
		return err
	}

	return nil
}

func NewTenantMigrations(
	log common.Logger,
	config *config.Config,
	tenantManager *common.TenantConnectionManager,
	adminMigrationsFS embed.FS,
	tenantMigrationsFS embed.FS,
	commonMigrationsFS embed.FS) TenantMigrations {
	return &tenantMigrations{
		log:                log,
		config:             config,
		tenantManager:      tenantManager,
		adminMigrationsFS:  adminMigrationsFS,
		tenantMigrationsFS: tenantMigrationsFS,
		commonMigrationsFS: commonMigrationsFS,
	}
}

var _ TenantMigrations = (*tenantMigrations)(nil)
