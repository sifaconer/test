package api

import "github.com/gofiber/fiber/v2"

func (r *Rest) LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		method := c.Method()
		path := c.Path()

		err := c.Next()
		if err != nil {
			status := c.Response().StatusCode()
			r.log.Error(c.Context(), method, "path", path, "status", status, "error", err.Error())
			return err
		}

		status := c.Response().StatusCode()
		if status >= 400 {
			r.log.Error(c.Context(), method, "path", path, "status", status)
		} else {
			r.log.Info(c.Context(), method, "path", path, "status", status)
		}

		return err
	}
}
