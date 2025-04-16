package api

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/database/postgres"
	"api-test/src/modules/admin/repository/implements"
	"api-test/src/modules/admin/usecase"
	"context"

	"github.com/gofiber/fiber/v2"
)

type AdminAPI struct {
	log    common.Logger
	app    fiber.Router
	routes AdminRoutes
	config *config.Config
	tenant *common.TenantConnectionManager
	ucTenant usecase.Tenant
}

func (t *AdminAPI) Register() {
	t.routes.RegisterRoutes()
}

func (t *AdminAPI) RegisterAllTenants(ctx context.Context) error {
	return t.ucTenant.RegisterAllTenants(ctx)
}

func NewAdminAPI(
	log common.Logger,
	app fiber.Router,
	config *config.Config,
	tenant *common.TenantConnectionManager,
	migrations usecase.TenantMigrations,
	psql postgres.Database) *AdminAPI {

	repoTenant := implements.NewTenantRepository(log, tenant)
	ucTenant := usecase.NewTenant(log, repoTenant, migrations, config, tenant, psql)
	repoUserDirectory := implements.NewUserRepository(log, tenant)
	ucAuth := usecase.NewAuth(log, config, tenant, repoUserDirectory)

	return &AdminAPI{
		log:    log,
		app:    app,
		ucTenant: ucTenant,
		routes: NewAdminRoutes(log, app, ucTenant, ucAuth, config),
		tenant: tenant,
	}
}
