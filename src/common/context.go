package common

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Global Context to store tenant and user ID
var (
	TenantKey = "TenantID"
	UserIDKey = "UserID"
)

func Context(c *fiber.Ctx) context.Context {
	tenantID, ok := c.Locals(TenantKey).(uuid.UUID)
	if !ok {
		NewLogger().Warn(c.Context(), "Tenant not found")
		return c.Context()
	}
	ctx := context.WithValue(c.Context(), TenantKey, tenantID)
	userID, ok := c.Locals(UserIDKey).(uuid.UUID)
	if !ok {
		NewLogger().Warn(c.Context(), "User not found")
		return ctx
	}
	ctx = context.WithValue(ctx, UserIDKey, userID)
	return ctx
}