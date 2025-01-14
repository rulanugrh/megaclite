package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/sujit-baniya/flash"
)

type CommonMiddlewareInterface interface {
	ViewMiddleware(c *fiber.Ctx) error
	APIMiddleware(c *fiber.Ctx) error
}

type common struct {
	conf *config.App
}

func NewCommonMiddleware(conf *config.App) CommonMiddlewareInterface {
	return &common{
		conf: conf,
	}
}
func (cm *common) ViewMiddleware(c *fiber.Ctx) error {
	fm := fiber.Map{
		"type": "error",
	}
	session, err := cm.conf.Store.Get(c)
	if session.Get("Authorization") == nil {
		fm["message"] = "Your are not authorized"

		c.Locals("protected", false)
		return flash.WithError(c, fm).Next()
	}
	if err != nil {
		fm["message"] = "Cannot Create New Session"
		c.Locals("protected", false)
		return flash.WithError(c, fm).Next()
	}

	c.Locals("Authorization", session.Get("Authorization"))
	c.Locals("protected", true)
	return c.Next()
}

func (cm *common) APIMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(400).JSON(web.Unauthorized("Sorry you're not authorized"))
	}

	validToken := verifyToken(token)
	if !validToken {
		return c.Status(403).JSON(web.Forbidden("Sorry your token is invalid"))
	}

	return c.Next()
}
