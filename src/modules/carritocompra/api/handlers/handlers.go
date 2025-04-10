package handlers

import "github.com/gofiber/fiber/v2"

type CarritoCompraHandler interface {
	// Create
	Create(c *fiber.Ctx) error
	// Read
	Get(c *fiber.Ctx) error
	// Update
	Update(c *fiber.Ctx) error
	// Delete
	Delete(c *fiber.Ctx) error
	// List
	List(c *fiber.Ctx) error
}
