package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/config"
	"github.com/sujit-baniya/flash"
)

type CommonMiddlewareInterface interface {
	ViewMiddleware(c *fiber.Ctx) error
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
