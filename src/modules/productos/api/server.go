package api

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/productos/domain"
	"api-test/src/modules/productos/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProductosAPI struct {
	log    common.Logger
	config *config.Config
	app    fiber.Router
	uc     usecase.ProductosUseCase
	routes ProductosRoutes
	tenant *common.TenantConnectionManager
}

func NewProductosAPI(log common.Logger, app fiber.Router, config *config.Config, tenant *common.TenantConnectionManager) *ProductosAPI {
	repo := common.NewRepository[domain.ProductosTable, int64](config, log, tenant)
	uc := usecase.NewProductosUseCase(config, log, tenant, repo)
	routes := NewProductosRoutes(log, uc, app)
	return &ProductosAPI{
		log:    log,
		config: config,
		app:    app,
		uc:     uc,
		routes: routes,
		tenant: tenant,
	}
}

func (api *ProductosAPI) Register() {
	api.routes.RegisterRoutes()
}
