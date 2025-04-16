package handlers

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/admin/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MigrationsHandler struct {
	log common.Logger
	uc  usecase.TenantMigrations
	config config.Config
}

func (m *MigrationsHandler) RunAdminMigrations(c *fiber.Ctx) error {
	// Use case
	err := m.uc.RunAdminMigrations(common.Context(c))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}
	return c.Status(fiber.StatusOK).JSON(common.Response[any]{
		Status:  "success",
		Code:    fiber.StatusOK,
		Message: "Admin migrations run successfully",
	})
}

func (m *MigrationsHandler) RunTenantMigrations(c *fiber.Ctx) error {
	// Use case
	err := m.uc.RunAllMigrations(common.Context(c), uuid.Nil)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusUnauthorized,
			Message: "Unauthorized",
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}
	return c.Status(fiber.StatusOK).JSON(common.Response[any]{
		Status:  "success",
		Code:    fiber.StatusOK,
		Message: "Tenant migrations run successfully",
	})
}

func NewMigrationsHandler(log common.Logger, uc usecase.TenantMigrations, config config.Config) *MigrationsHandler {
	return &MigrationsHandler{
		log: log,
		uc:  uc,
		config: config,
	}
}


