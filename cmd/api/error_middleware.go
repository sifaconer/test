package api

import (
	"api-test/src/common"

	"github.com/gofiber/fiber/v2"
)

func (r *Rest) ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err == nil {
			return nil
		}

		var statusCode int
		var errorResponse any

		switch e := err.(type) {
		case *fiber.Error:
			statusCode = e.Code
			errorResponse = common.Response[any]{
				Status:  "error",
				Code:    e.Code,
				Message: e.Message,
				Errors:  []common.APIError{{Message: e.Message}},
			}

		default:
			return err
		}

		return c.Status(statusCode).JSON(errorResponse)
	}
}
