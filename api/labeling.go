package api

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

func NewLabelMailAPI(service service.LabelingInterface) LabelingInterface {
	return &labeling{
		service: service,
	}
}

// @Summary add mail to label
// @ID adding
// @Tags labelings
// @Accept json
// @Produce json
// @Param request body domain.MailLabelRegister true "body request for add label mail"
// @Router /api/label/add [post]
// @Success 201 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (l *labeling) Create(c *fiber.Ctx) error {
	// Parsing request to body parse
	var request domain.MailLabelRegister
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

// @Summary find label by id
// @ID findByID
// @Tags labelings
// @Accept json
// @Produce json
// @Param id path int true "id label"
// @Router /api/label/find/{id} [get]
// @Success 201 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
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

// @Summary find by category
// @ID findByCategory
// @Tags labelings
// @Accept json
// @Produce json
// @Param user_id path int true "user id"
// @Param category path int true "category id"
// @Route /api/label/get/{user_id}/{category}
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
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

// @Summary update label by user id
// @ID updateLabel
// @Tags labelings
// @Accept json
// @Produce json
// @Param request body domain.MailLabelRegister true "request body for update label"
// @Param id path int true "parameter id"
// @Route /api/label/update/{id} [put]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (l *labeling) UpdateLabel(c *fiber.Ctx) error {
	// parsng request to body parser
	var request domain.MailLabelRegister
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
