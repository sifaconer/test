package handlers

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/admin/domain"
	"api-test/src/modules/admin/usecase"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	log    common.Logger
	uc     usecase.Auth
	config config.Config
}

func (a *authHandler) Login(c *fiber.Ctx) error {
	// Decode
	dto := domain.DTOLogin{}
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
	auth, err := a.uc.Login(c.Context(), domain.DTOUserDirectory{
		Email:    dto.Email,
		Password: dto.Password,
	})
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
		Message: "Login successful",
		Data:    auth,
	})
}

func (a *authHandler) Register(c *fiber.Ctx) error {
	// Decode
	dto := domain.DTORegister{}
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
	auth, err := a.uc.Register(c.Context(), domain.DTOUserDirectory{
		Email:    dto.Email,
		Password: dto.Password,
		Name:     dto.Name,
	})
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
		Message: "Register successful",
		Data:    auth,
	})
}

func (a *authHandler) Logout(c *fiber.Ctx) error {
	return nil
}

func (a *authHandler) Refresh(c *fiber.Ctx) error {
	// Decode
	dto := domain.DTOAuth{}
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
	auth, err := a.uc.Refresh(common.Context(c), &dto)
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
		Message: "Refresh successful",
		Data:    auth,
	})
}

// ForgotPassword maneja el olvido de la contraseña.
func (a *authHandler) ForgotPassword(c *fiber.Ctx) error {
	return nil
}

// ResetPassword maneja el cambio de la contraseña.
func (a *authHandler) ResetPassword(c *fiber.Ctx) error {
	return nil
}

func NewAuthHandler(log common.Logger, config config.Config, uc usecase.Auth) *authHandler {
	return &authHandler{
		log:    log,
		config: config,
		uc:     uc,
	}
}

var _ AuthHandler = (*authHandler)(nil)
