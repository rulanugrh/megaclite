package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/service"
)

type LabelingInterface interface {
	Create(c *fiber.Ctx) error
	FindByID(c *fiber.Ctx) error
	FindByCategory(c *fiber.Ctx) error
	UpdateLabel(c *fiber.Ctx) error
}

type labeling struct {
	service service.LabelingInterface
}

func NewLabelMailHandler(service service.LabelingInterface) LabelingInterface {
	return &labeling{
		service: service,
	}
}

func (l *labeling) Create(c *fiber.Ctx) error {
	// Parsing request to body parse
	var request domain.MailLabel
	// Checking error while parser request
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parser Body Request"))
	}

	// Parsing fronm request to create service layer
	data, err := l.service.Create(request)
	// Check error while create request
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return response json while success created
	return c.Status(201).JSON(web.Created("Success Create New Label Mail", data))
}

func (l *labeling) FindByID(c *fiber.Ctx) error {
	// get parameter from url
	getID := c.Params("id")
	// convert string to integer
	id, err := strconv.Atoi(getID)
	// checking error while parsing request id
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing ID"))
	}

	// process get data in service layer
	data, err := l.service.FindByID(uint(id))
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return success if success get data
	return c.Status(200).JSON(web.Success("Success Get Mail with This ID", data))
}

func (l *labeling) FindByCategory(c *fiber.Ctx) error {
	// get parameter from url
	getCategory := c.Params("category")
	getUID := c.Params("user_id")

	// convert user id to integer
	userID, err := strconv.Atoi(getUID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing User ID"))
	}

	// convert category id to integer
	categoryID, err := strconv.Atoi(getCategory)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Category ID"))
	}

	// process get data find by category into service layer
	data, err := l.service.FindByCategory(uint(categoryID), uint(userID))
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return success get data
	return c.Status(200).JSON(web.Success("Success Get Mail with This Category", data))
}

func (l *labeling) UpdateLabel(c *fiber.Ctx) error {
	// parsng request to body parser
	var request domain.MailLabel
	// checking error while parsing request
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parser Body Request"))
	}

	// get user id from parameter
	getID := c.Params("id")
	// conver id into integer from parameter
	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Request ID"))
	}

	getCategoryID := c.Params("categoryID")
	categoryID, err := strconv.Atoi(getCategoryID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Request ID"))
	}

	// process update label to service layer
	data, err := l.service.UpdateLabel(uint(id), uint(categoryID))
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return success update label mail
	return c.Status(200).JSON(web.Success("Success Update Label Mail", data))
}
