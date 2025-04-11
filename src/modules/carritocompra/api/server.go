package api

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/carritocompra/domain"
	"api-test/src/modules/carritocompra/usecase"

	"github.com/gofiber/fiber/v2"
)

type CarritoCompraAPI struct {
	log    common.Logger
	config *config.Config
	app    *fiber.App
	uc     usecase.CarritoCompra
	routes CarritoCompraRoutes
	tenant *common.TenantConnectionManager
}

func NewCarritoCompraAPI(log common.Logger, app *fiber.App, config *config.Config, tenant *common.TenantConnectionManager) *CarritoCompraAPI {
	repo := common.NewRepository[domain.TableCarritoCompra, int64](config, log, tenant)
	uc := usecase.NewCarritoCompra(config, log, tenant, repo)
	routes := NewCarritoCompraRoutes(log, app, uc)

	return &CarritoCompraAPI{
		log:    log,
		config: config,
		uc:     uc,
		app:    app,
		routes: routes,
		tenant: tenant,
	}
}

func (api *CarritoCompraAPI) Register() {
	api.routes.RegisterRoutes()
}
