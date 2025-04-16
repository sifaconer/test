package api

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/admin/api/handlers"
	"api-test/src/modules/admin/usecase"

	"github.com/gofiber/fiber/v2"
)

type AdminRoutes interface {
	RegisterRoutes()
}

type adminRoutes struct {
	log            common.Logger
	ucTenant       usecase.Tenant
	ucTenantMigration usecase.TenantMigrations
	ucAuth         usecase.Auth
	config         *config.Config
	app            fiber.Router
	tenantHandlers *handlers.TenantHandler
	authHandlers   *handlers.AuthHandler
	migrationsHandlers *handlers.MigrationsHandler
}

func (t *adminRoutes) RegisterRoutes() {

	// Auth
	t.app.Post("/login", t.authHandlers.Login)
	t.app.Post("/register", t.authHandlers.Register)
	t.app.Post("/logout", t.authHandlers.Logout)
	t.app.Post("/refresh", t.authHandlers.Refresh)

	// Tenant
	t.app.Get("/tenants", t.tenantHandlers.List)
	t.app.Post("/tenants", t.tenantHandlers.Create)

	// Migrations
	t.app.Post("/migrations/admin", t.migrationsHandlers.RunAdminMigrations)
	t.app.Post("/migrations/tenant", t.migrationsHandlers.RunTenantMigrations)
}

func NewAdminRoutes(log common.Logger, app fiber.Router, ucTenant usecase.Tenant, ucTenantMigration usecase.TenantMigrations, ucAuth usecase.Auth, config *config.Config) AdminRoutes {
	
	return &adminRoutes{
		log:            log,
		ucTenant:       ucTenant,
		ucTenantMigration: ucTenantMigration,
		ucAuth:         ucAuth,
		config:         config,
		app:            app,
		tenantHandlers: handlers.NewTenantHandler(log, ucTenant),
		authHandlers:   handlers.NewAuthHandler(log, *config, ucAuth),
		migrationsHandlers: handlers.NewMigrationsHandler(log, ucTenantMigration, *config),
	}
}

var _ AdminRoutes = (*adminRoutes)(nil)
