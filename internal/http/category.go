package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rulanugrh/megaclite/internal/entity/domain"
	"github.com/rulanugrh/megaclite/internal/entity/web"
	"github.com/rulanugrh/megaclite/internal/service"
)

type CategoryInterface interface {
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
}

type category struct {
	service service.CategoryInterface
}

func NewCategoryHandler(service service.CategoryInterface) CategoryInterface {
	return &category{
		service: service,
	}
}

func (ct *category) Create(c *fiber.Ctx) error {
	// parsing request to body parser
	var request domain.Category
	err := c.BodyParser(&request)
	// checking error while request body parser
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Request"))
	}

	// parse reqeust to create from service layer
	data, err := ct.service.Create(request)
	// check error while create
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return success while create new category
	return c.Status(201).JSON(web.Created("Success Create New Category", data))
}

func (ct *category) Delete(c *fiber.Ctx) error {
	// get request id from url parameter
	getID := c.Params("id")
	// convert id string to integer
	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing ID"))
	}

	// process delete by id
	err = ct.service.Delete(uint(id))
	// checking error after deleted
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return success
	return c.Status(200).JSON(web.Success("Success Delete Category", nil))
}

func (ct *category) Update(c *fiber.Ctx) error {
	// parsing body request to entity
	var request domain.Category
	err := c.BodyParser(&request)
	// check error after parsing body request
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing Body Request"))
	}

	// get id from url parameter
	getID := c.Params("id")
	// convert string to integer
	id, err := strconv.Atoi(getID)
	if err != nil {
		return c.Status(500).JSON(web.InternalServerError("Cannot Parsing ID Request"))
	}

	// process update from request and parameter url
	data, err := ct.service.Update(uint(id), request)
	if err != nil {
		return c.Status(400).JSON(web.BadRequest(err.Error()))
	}

	// return success
	return c.Status(200).JSON(web.Success("Success Update Category", data))
}
