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
	ucAuth         usecase.Auth
	config         *config.Config
	app            fiber.Router
	tenantHandlers handlers.TenantHandler
	authHandlers   handlers.AuthHandler
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
	t.app.Get("/tenants/:id", t.tenantHandlers.Get)
	t.app.Put("/tenants/:id", t.tenantHandlers.Update)
	t.app.Delete("/tenants/:id", t.tenantHandlers.Delete)
}

func NewAdminRoutes(log common.Logger, app fiber.Router, ucTenant usecase.Tenant, ucAuth usecase.Auth, config *config.Config) AdminRoutes {
	return &adminRoutes{
		log:            log,
		ucTenant:       ucTenant,
		ucAuth:         ucAuth,
		config:         config,
		app:            app,
		tenantHandlers: handlers.NewTenantHandler(log, ucTenant),
		authHandlers:   handlers.NewAuthHandler(log, *config, ucAuth),
	}
}

var _ AdminRoutes = (*adminRoutes)(nil)
