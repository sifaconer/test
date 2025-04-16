package api

import (
	"api-test/src/common"
	"api-test/src/modules/productos/api/handlres"
	"api-test/src/modules/productos/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProductosRoutes interface {
	RegisterRoutes()
}

type productosRoutes struct {
	log      common.Logger
	uc       usecase.ProductosUseCase
	handlers *handlres.ProductosHandler
	app      fiber.Router
}

func (r *productosRoutes) RegisterRoutes() {
	r.app.Get("/productos", r.handlers.Search)
	r.app.Get("/productos/:id", r.handlers.Get)
	r.app.Post("/productos", r.handlers.Create)
	r.app.Put("/productos/:id", r.handlers.Update)
	r.app.Delete("/productos/:id", r.handlers.Delete)
}

func NewProductosRoutes(log common.Logger, uc usecase.ProductosUseCase, app fiber.Router) ProductosRoutes {
	return &productosRoutes{
		log:      log,
		uc:       uc,
		handlers: handlres.NewProductosHandler(log, uc),
		app:      app,
	}
}