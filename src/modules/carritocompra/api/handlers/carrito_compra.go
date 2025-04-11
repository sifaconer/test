package handlers

import (
	"api-test/src/common"
	"api-test/src/modules/carritocompra/domain"
	"api-test/src/modules/carritocompra/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CarritoCompraHandler struct {
	log common.Logger
	uc  usecase.CarritoCompra
}

// Create implements CarritoCompraHandler.
func (cc *CarritoCompraHandler) Create(c *fiber.Ctx) error {
	// Decode
	dto := domain.DTOCarritoCompra{}
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
	_, err := cc.uc.Create(common.Context(c), dto)
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

	// Response
	return c.Status(fiber.StatusCreated).JSON(common.Response[any]{
		Status:  "success",
		Code:    fiber.StatusCreated,
		Message: "Success",
	})
}

// Delete implements CarritoCompraHandler.
func (cc *CarritoCompraHandler) Delete(c *fiber.Ctx) error {
	// Decode
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:     "error",
			Code:       fiber.StatusBadRequest,
			Message:    "Invalid request body",
			Data:       nil,
			Pagination: nil,
			Query:      nil,
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Use case
	err = cc.uc.Delete(common.Context(c), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:     "error",
			Code:       fiber.StatusBadRequest,
			Message:    "Internal server error",
			Data:       nil,
			Pagination: nil,
			Query:      nil,
			Errors: []common.APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusOK).JSON(common.Response[any]{
		Status:     "success",
		Code:       fiber.StatusOK,
		Message:    "Success",
		Data:       nil,
		Pagination: nil,
		Query:      nil,
		Errors:     nil,
	})
}

// Get implements CarritoCompraHandler.
func (cc *CarritoCompraHandler) Get(c *fiber.Ctx) error {
	// Decode
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
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

	// Use case
	dto, err := cc.uc.GetById(common.Context(c), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Bad request",
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
		Data:    dto,
	})
}

// List implements CarritoCompraHandler.
func (cc *CarritoCompraHandler) List(c *fiber.Ctx) error {
	// Decode
	filters := common.QueryParams{}
	if err := c.QueryParser(&filters); err != nil {
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
	// Use case
	dtos, err := cc.uc.Search(common.Context(c), nil) // TODO: Handle filters
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
		Data:    dtos,
	})
}

// Update implements CarritoCompraHandler.
func (cc *CarritoCompraHandler) Update(c *fiber.Ctx) error {
	// Decode
	dto := domain.DTOCarritoCompra{}
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:     "error",
			Code:       fiber.StatusBadRequest,
			Message:    "Invalid request body",
			Data:       nil,
			Pagination: nil,
			Query:      nil,
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
			Status:     "error",
			Code:       fiber.StatusBadRequest,
			Message:    "Validation error",
			Data:       nil,
			Pagination: nil,
			Query:      nil,
			Errors:     validationErrors,
		})
	}

	// Use case
	_, err := cc.uc.Update(common.Context(c), dto.ID, dto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.Response[any]{
			Status:     "error",
			Code:       fiber.StatusBadRequest,
			Message:    "Internal server error",
			Data:       nil,
			Pagination: nil,
			Query:      nil,
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

func NewCarritoCompraHandler(log common.Logger, uc usecase.CarritoCompra) *CarritoCompraHandler {
	return &CarritoCompraHandler{
		log: log,
		uc:  uc,
	}
}
