package handlers

import (
	"api-test/src/common"
	"api-test/src/modules/admin/domain"
	"api-test/src/modules/admin/usecase"

	"github.com/gofiber/fiber/v2"
)

type tenantHandler struct {
	log common.Logger
	uc  usecase.Tenant
}

// Create implements TenantHandler.
func (t *tenantHandler) Create(c *fiber.Ctx) error {
	// Decode
	dto := domain.DTOTenant{}
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}
	// Validate
	if validationErrors := common.Validate(dto); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Validation error",
			Errors:  validationErrors,
		})
	}

	// Use case
	tenant, err := t.uc.CreateTenant(common.Context(c), dto)
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
	return c.Status(fiber.StatusCreated).JSON(common.Response[any]{
		Status:  "success",
		Code:    fiber.StatusCreated,
		Message: "Tenant created successfully",
		Data:    tenant,
	})
}

// Delete implements TenantHandler.
func (t *tenantHandler) Delete(c *fiber.Ctx) error {
	panic("unimplemented")
}

// Get implements TenantHandler.
func (t *tenantHandler) Get(c *fiber.Ctx) error {
	panic("unimplemented")
}

// List implements TenantHandler.
func (t *tenantHandler) List(c *fiber.Ctx) error {
	// Use case
	tenants, err := t.uc.ListTenants(common.Context(c))
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
		Message: "Success",
		Data:    tenants,
	})
}

// Update implements TenantHandler.
func (t *tenantHandler) Update(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewTenantHandler(log common.Logger, uc usecase.Tenant) *tenantHandler {
	return &tenantHandler{
		log: log,
		uc:  uc,
	}
}

var _ TenantHandler = (*tenantHandler)(nil)
