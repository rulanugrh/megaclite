package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/config"
	"github.com/sujit-baniya/flash"
)

func ViewMiddleware(c *fiber.Ctx) error {
	fm := fiber.Map{
		"type": "error",
	}
	session, err := config.Store.Get(c)
	if err != nil {
		fm["message"] = "Cannot Create New Session"

		return flash.WithError(c, fm).Redirect("/")
	}

	if session.Get("Authorization") == nil {
		fm["message"] = "Your are not authorized"

		return flash.WithError(c, fm).Redirect("/")
	}

	c.Locals("Authorization", session.Get("Authorization"))
	return c.Next()
}
