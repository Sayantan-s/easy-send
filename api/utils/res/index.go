package res

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SuccessTemplate struct{
	StatusCode int
	Data interface{}
	Message string
}

type FalureTemplate struct{
	StatusCode int
	Error interface{}
	Message string
}

func Success(c *fiber.Ctx, r SuccessTemplate) error{
	return c.Status(r.StatusCode).JSON(fiber.Map{
		"requestId": uuid.New(),
		"status": "success",
		"message": r.Message,
		"data": r.Data,
	})
}

func Failure(c *fiber.Ctx, r FalureTemplate) error{
	return c.Status(r.StatusCode).JSON(fiber.Map{
		"requestId": uuid.New(),
		"status": "failure",
		"message": r.Message,
		"error": r.Error,
	})
}