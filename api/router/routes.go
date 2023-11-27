package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayantan-s/easy-send/router/nestedRoutes"
)

func SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	nestedRoutes.GenerationRoutes(api)
}