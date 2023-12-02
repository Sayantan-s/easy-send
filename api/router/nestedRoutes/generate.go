package nestedRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sayantan-s/easy-send/handlers"
)

func GenerationRoutes(router fiber.Router){
	rtr := router.Group("/generate")
	
	rtr.Post("/transcript_CE", handlers.InitiateTranscriptions)
	rtr.Post("/transcript_CE_webhook", handlers.GetGeneratedTranscripts)
}
