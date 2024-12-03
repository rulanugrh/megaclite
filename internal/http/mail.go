package handler

import "github.com/gofiber/fiber/v2"

type MailView interface {
	InboxView(c *fiber.Ctx) error
	ArchiveView(c *fiber.Ctx) error
	SpamView(c *fiber.Ctx) error
	SentView(c *fiber.Ctx) error
	TrashView(c *fiber.Ctx) error
}
