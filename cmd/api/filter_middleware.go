package api

import "github.com/gofiber/fiber/v2"

/*
filter	string	Filtros avanzados (ej: AND(name:like:Jugo,brand.name:eq:Natura))
order	string	Orden de resultados (ej: ASC(created_at), DESC(price))

Este middleware se encarga de aplicar filtros y ordenamientos a las consultas SQL, toma los parámetros filter y order de los query params y los parsea a un struct que se enviara por el contexto para que los querys dinámicos puedan usarlos en las consultas sql
*/

func (r *Rest) FilterMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}