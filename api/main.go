package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/sayantan-s/easy-send/config"
	"github.com/sayantan-s/easy-send/router"
)

func main(){
	PORT := config.GetConfig("PORT")
	ServerDomain := fmt.Sprintf("localhost:%s", PORT)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	
	router.SetupRoutes(app)
	
	app.Get("/status", func (c *fiber.Ctx)  error{
		return c.SendString("OK")
	})

	log.Fatal(app.Listen(ServerDomain))
}