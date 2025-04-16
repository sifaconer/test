package common

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


type GenericHandler[DTO any, ID any] struct {
	Log    Logger
	UseCase UseCase[DTO, ID]
	ParseID func(string) (ID, error) 
}


func (h *GenericHandler[DTO, ID]) Create(c *fiber.Ctx) error {
	// Decode
	var dto DTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}
	
	// Validate
	if validationErrors := Validate(dto); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Validation error",
			Errors:  validationErrors,
		})
	}

	// Use case
	result, err := h.UseCase.Create(Context(c), dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusInternalServerError,
			Message: "Internal server error",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusCreated).JSON(Response[any]{
		Status:  "success",
		Code:    fiber.StatusCreated,
		Message: "Resource created successfully",
		Data:    result,
	})
}

func (h *GenericHandler[DTO, ID]) Get(c *fiber.Ctx) error {
	// Decode
	idParam := c.Params("id")
	id, err := h.ParseID(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ID format",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Use case
	result, err := h.UseCase.GetById(Context(c), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusInternalServerError,
			Message: "Error retrieving resource",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusOK).JSON(Response[any]{
		Status:  "success",
		Code:    fiber.StatusOK,
		Message: "Resource retrieved successfully",
		Data:    result,
	})
}

func (h *GenericHandler[DTO, ID]) Search(c *fiber.Ctx) error {
	// Parse query parameters
	filters := QueryParams{}
	if err := c.QueryParser(&filters); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid query parameters",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Use case
	result, err := h.UseCase.Search(Context(c), &filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusInternalServerError,
			Message: "Error retrieving resources",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusOK).JSON(Response[any]{
		Status:  "success",
		Code:    fiber.StatusOK,
		Message: "Resources retrieved successfully",
		Data:    result,
	})
}

func (h *GenericHandler[DTO, ID]) Update(c *fiber.Ctx) error {
	// Decode ID
	idParam := c.Params("id")
	id, err := h.ParseID(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ID format",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Decode body
	var dto DTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid request body",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Validate
	if validationErrors := Validate(dto); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Validation error",
			Errors:  validationErrors,
		})
	}

	// Use case
	result, err := h.UseCase.Update(Context(c), id, dto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusInternalServerError,
			Message: "Error updating resource",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusOK).JSON(Response[any]{
		Status:  "success",
		Code:    fiber.StatusOK,
		Message: "Resource updated successfully",
		Data:    result,
	})
}

func (h *GenericHandler[DTO, ID]) Delete(c *fiber.Ctx) error {
	// Decode
	idParam := c.Params("id")
	id, err := h.ParseID(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusBadRequest,
			Message: "Invalid ID format",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Use case
	err = h.UseCase.Delete(Context(c), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response[any]{
			Status:  "error",
			Code:    fiber.StatusInternalServerError,
			Message: "Error deleting resource",
			Errors: []APIError{
				{
					Message: err.Error(),
				},
			},
		})
	}

	// Response
	return c.Status(fiber.StatusOK).JSON(Response[any]{
		Status:  "success",
		Code:    fiber.StatusOK,
		Message: "Resource deleted successfully",
	})
}

func ParseID[ID any](idStr string) (ID, error) {
	var id ID
	if idStr == "" {
		return id, errors.New("ID cannot be empty")
	}

	// parse segun el tipo de ID (int64, uuid, string)
	switch any(id).(type) {
	case int64:
		parsedID, err := ParseInt64ID(idStr)
		if err != nil {
			return id, err
		}
		return parsedID.(ID), nil
	case string:
		parsedID, err := ParseStringID(idStr)
		if err != nil {
			return id, err
		}
		return parsedID.(ID), nil
	case uuid.UUID:
		parsedID, err := ParseUUID(idStr)
		if err != nil {
			return id, err
		}
		return parsedID.(ID), nil
	default:
		return id, errors.New("unsupported ID type")
	}
}

func ParseInt64ID(idStr string) (any, error) {
	return strconv.ParseInt(idStr, 10, 64)
}

func ParseStringID(idStr string) (any, error) {
	return idStr, nil
}

func ParseUUID(idStr string) (any, error) {
	return uuid.Parse(idStr)
}

func NewGenericHandler[DTO any, ID any](log Logger, useCase UseCase[DTO, ID]) *GenericHandler[DTO, ID] {
	return &GenericHandler[DTO, ID]{
		Log:     log,
		UseCase: useCase,
		ParseID: ParseID[ID],
	}
}