package api

import (
	"api-test/src/common"
	"api-test/src/modules/carritocompra/api/handlers"
	"api-test/src/modules/carritocompra/usecase"

	"github.com/gofiber/fiber/v2"
)

type CarritoCompraRoutes interface {
	RegisterRoutes()
}

type carritoCompraRoutes struct {
	log      common.Logger
	uc       usecase.CarritoCompra
	handlers *handlers.CarritoCompraHandler
	app      fiber.Router
}

func (c *carritoCompraRoutes) RegisterRoutes() {
	c.app.Get("/carrito-compra", c.handlers.Search)
	c.app.Post("/carrito-compra", c.handlers.Create)
	c.app.Get("/carrito-compra/:id", c.handlers.Get)
	c.app.Put("/carrito-compra/:id", c.handlers.Update)
	c.app.Delete("/carrito-compra/:id", c.handlers.Delete)
}

func NewCarritoCompraRoutes(log common.Logger, app fiber.Router, uc usecase.CarritoCompra) CarritoCompraRoutes {
	return &carritoCompraRoutes{
		log:      log,
		uc:       uc,
		app:      app,
		handlers: handlers.NewCarritoCompraHandler(log, uc),
	}
}
