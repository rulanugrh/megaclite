package handler

import "github.com/gofiber/fiber/v2"

type CategoryInterface interface {
	Create(c *fiber.Ctx) error
}
