package database

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/database/postgres"
	"api-test/src/modules/admin/usecase"
	"context"
	"embed"
)

type Database struct {
	conf            *config.Config
	log             common.Logger
	tenant          *common.TenantConnectionManager
	adminMigrations usecase.TenantMigrations
	psql            postgres.Database
}

func NewDatabase(conf *config.Config, log common.Logger, tenant *common.TenantConnectionManager,
	adminMigrationsFS,
	tenantMigrationsFS,
	commonMigrationsFS embed.FS) *Database {
	adminMigrations := usecase.NewTenantMigrations(
		log,
		conf,
		tenant,
		adminMigrationsFS,
		tenantMigrationsFS,
		commonMigrationsFS)
	return &Database{
		conf:            conf,
		log:             log,
		tenant:          tenant,
		adminMigrations: adminMigrations,
	}
}

func (d *Database) PSQL() postgres.Database {
	return d.psql
}

func (d *Database) Run() error {
	d.log.Info(context.Background(), "Starting KOSVI Database")
	psql := postgres.NewPostgres(d.log, d.conf)
	err := d.tenant.RegisterTenant(&common.TenantConfig{
		TenantID:         d.conf.TenantID,
		Name:             "KOSVI",
		ConnectionString: d.conf.DSN,
	}, psql.Connect)
	if err != nil {
		return err
	}
	d.log.Info(context.Background(), "Database KOSVI Started")
	if err := d.AdminMigrations(); err != nil {
		return err
	}
	d.psql = psql
	return nil
}

func (d *Database) AdminMigrations() error {
	d.log.Info(context.Background(), "Running Admin Migrations")
	err := d.adminMigrations.RunAdminMigrations(context.Background())
	if err != nil {
		return err
	}
	d.log.Info(context.Background(), "Admin Migrations Completed")

	return nil
}

func (d *Database) Stop() error {
	d.log.Info(context.Background(), "Stopping Database")
	return d.tenant.CloseAll()
}
