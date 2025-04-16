package api

import (
	"api-test/src/common"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Tenant Middleware para capturar el tenant ID
func (r *Rest) TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Validar las rutas excluidas
		if slices.Contains(r.EXCLUDE_PATHS, c.Path()) {
			return c.Next()
		}
		// // Excluir la ruta del POST para /tenants
		if c.Method() == "POST" && c.Path() == "/api/v1/tenants" {
			return c.Next()
		}
		// Excluir la ruta del GET para /tenants
		if c.Method() == "GET" && c.Path() == "/api/v1/tenants" {
			return c.Next()
		}

		// Capturar el tenant ID del header "X-Tenant-Id"
		tenantID := c.GetReqHeaders()["X-Tenant-Id"]
		if len(tenantID) == 0 {
			r.log.Error(c.Context(), "Tenant Middleware", "path", c.Path(), "status", 401, "error", "No tenant found")
			return fiber.NewError(401, "No tenant found assigned")
		}
		tenantUUID, err := uuid.Parse(tenantID[0])
		if err != nil {
			r.log.Error(c.Context(), "Tenant Middleware", "path", c.Path(), "status", 401, "error", err.Error())
			return fiber.NewError(401, "Invalid tenant ID")
		}

		// Validar el JWT
		token := c.GetReqHeaders()["Authorization"]
		if len(token) == 0 {
			r.log.Error(c.Context(), "Tenant Middleware", "path", c.Path(), "status", 401, "error", "No token found")
			return fiber.NewError(401, "No token found")
		}

		// Validar el token
		err = common.ValidateTenant(c.Context(), r.conf, tenantUUID, token[0])
		if err != nil {
			r.log.Error(c.Context(), "Tenant Middleware", "path", c.Path(), "status", 401, "error", err.Error())
			return fiber.NewError(401, err.Error())
		}

		// Inyectar el tenant ID en el contexto
		c.Locals(r.tenant.TenantKey, tenantUUID)

		return c.Next()
	}
}