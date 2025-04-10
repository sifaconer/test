package api

import (
	"api-test/src/common"
	"api-test/src/modules/carritocompra/repository/implements"
	"api-test/src/modules/carritocompra/usecase"

	"github.com/gofiber/fiber/v2"
)

type CarritoCompraAPI struct {
	log    common.Logger
	app    *fiber.App
	uc     usecase.CarritoCompra
	routes CarritoCompraRoutes
	tenant *common.TenantConnectionManager
}

func NewCarritoCompraAPI(log common.Logger, app *fiber.App, tenant *common.TenantConnectionManager) *CarritoCompraAPI {
	repo := implements.NewCarritoCompraRepository(log, tenant)
	uc := usecase.NewCarritoCompra(log, repo)
	routes := NewCarritoCompraRoutes(log, app, uc)

	return &CarritoCompraAPI{
		log:    log,
		uc:     uc,
		app:    app,
		routes: routes,
		tenant: tenant,
	}
}

func (api *CarritoCompraAPI) Register() {
	api.routes.RegisterRoutes()
}
