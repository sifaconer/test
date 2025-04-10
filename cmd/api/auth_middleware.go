package api

import (
	"api-test/src/common"
	"slices"

	"github.com/gofiber/fiber/v2"
)

// JWT JSON Web Token
func (r *Rest) AuthenticationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Validar las rutas excluidas
		if slices.Contains(r.EXCLUDE_PATHS, c.Path()) {
			return c.Next()
		}

		// Extraer el token
		token := c.GetReqHeaders()["Authorization"]
		if len(token) == 0 {
			r.log.Error(c.Context(), "Auth Middleware", "path", c.Path(), "status", 401, "error", "No token found")
			return fiber.NewError(401, "No token found")
		}

		// Validar el token
		claims, err := common.ValidateJWT(c.Context(), token[0], r.conf)
		if err != nil {
			r.log.Error(c.Context(), "Auth Middleware", "path", c.Path(), "status", 401, "error", err.Error())
			return fiber.NewError(401, err.Error())
		}

		// Almacenar el user ID en el contexto
		c.Locals(r.tenant.UserIDKey, claims.UserID)

		return c.Next()
	}
}

// RBAC Role Based Access Control
func (r *Rest) AuthorizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			status := c.Response().StatusCode()
			r.log.Error(c.Context(), "Auth Middleware", "path", c.Path(), "status", status, "error", err.Error())
			return err
		}

		role := c.Get("role")
		if role == "" {
			r.log.Error(c.Context(), "Auth Middleware", "path", c.Path(), "status", 401, "error", "No role found")
			return fiber.NewError(401, "No role found")
		}

		return err
	}
}
