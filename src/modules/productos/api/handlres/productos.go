package handlres

import (
	"api-test/src/common"
	"api-test/src/modules/productos/domain"
	"api-test/src/modules/productos/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProductosHandler struct {
	log common.Logger
	uc  usecase.ProductosUseCase
}

func (h *ProductosHandler) Create(c *fiber.Ctx) error {
	// Decode
	dto := domain.ProductosDTO{}
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
	_, err := h.uc.Create(common.Context(c), *dto.ToTable())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Internal server error",
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusCreated).JSON(common.Response[any]{
		Status:  "success",
		Code:    fiber.StatusCreated,
		Message: "Success",
	})
}

func (h *ProductosHandler) GetById(c *fiber.Ctx) error {
	// Decode
	dto := domain.ProductosDTO{}
	if err := c.ParamsParser(&dto); err != nil {
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
	_, err := h.uc.GetById(common.Context(c), dto.Id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Internal server error",
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusCreated).JSON(common.Response[any]{
		Status:  "success",
		Code:    fiber.StatusCreated,
		Message: "Success",
	})
}

func (h *ProductosHandler) GetAll(c *fiber.Ctx) error {
	// Use case
	result, err := h.uc.Search(common.Context(c), nil)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Internal server error",
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	var dtos []domain.ProductosDTO
	for _, result := range result {
		dtos = append(dtos, *result.ToDTO())
	}

	// Response
	return c.Status(fiber.StatusOK).JSON(common.Response[any]{
		Status:  "success",
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    dtos,
	})
}

func (h *ProductosHandler) Update(c *fiber.Ctx) error {
	// Decode
	dto := domain.ProductosDTO{}
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
	_, err := h.uc.Update(common.Context(c), dto.Id, *dto.ToTable())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Internal server error",
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusCreated).JSON(common.Response[any]{
		Status:  "success",
		Code:    fiber.StatusCreated,
		Message: "Success",
	})
}

func (h *ProductosHandler) Delete(c *fiber.Ctx) error {
	// Decode
	dto := domain.ProductosDTO{}
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
	err := h.uc.Delete(common.Context(c), dto.Id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Internal server error",
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusOK).JSON(common.Response[any]{
		Status:  "success",
		Code:    fiber.StatusOK,
		Message: "Success",
	})
}

func NewProductosHandler(log common.Logger, uc usecase.ProductosUseCase) *ProductosHandler {
	return &ProductosHandler{
		log: log,
		uc:  uc,
	}
}