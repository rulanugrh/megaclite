package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (r Response) Error() string {
	return r.Message
}

func (r Response) HTTPCode() int {
	return r.Code
}

func InternalServerError(msg string) Response {
	return Response{
		Code:    500,
		Message: msg,
		Data:    nil,
	}
}

func Success(msg string, data any) Response {
	return Response{
		Code:    200,
		Message: msg,
		Data:    data,
	}
}

func Created(msg string, data any) Response {
	return Response{
		Code:    201,
		Message: msg,
		Data:    data,
	}
}

func NotFound(msg string) Response {
	return Response{
		Code:    404,
		Message: msg,
		Data:    nil,
	}
}

func Forbidden(msg string) Response {
	return Response{
		Code:    403,
		Message: msg,
		Data:    nil,
	}
}

func BadRequest(msg string) Response {
	return Response{
		Code:    400,
		Message: msg,
		Data:    nil,
	}
}

func RedirectView(c *fiber.Ctx, msg string, redirect string) error {
	return flash.WithError(c, fiber.Map{
		"type":    "error",
		"message": msg,
	}).Redirect(redirect)
}

func NextView(c *fiber.Ctx, msg string) error {
	return flash.WithError(c, fiber.Map{
		"type":    "error",
		"message": msg,
	}).Next()
}

func SuccessView(c *fiber.Ctx, msg string, data interface{}, redirect string) error {
	return flash.WithError(c, fiber.Map{
		"type":    "error",
		"message": msg,
		"data":    data,
	}).Redirect(redirect)
}
