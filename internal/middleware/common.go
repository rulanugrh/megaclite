package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/config"
	"github.com/rulanugrh/megaclite/internal/entity/web"
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
	session, _ := cm.conf.Store.Get(c)
	if session.Get("Authorization") == nil {

		c.Locals("protected", false)
		return c.Next()
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
