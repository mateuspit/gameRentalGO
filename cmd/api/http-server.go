package api

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"boardcamp/api"
	"boardcamp/api/handler"
)

func InitMicroService(db *gorm.DB) *fiber.App {
	fiberServer := fiber.New()

	log.Printf("Initialize web server\n")

	games := handler.NewGamesHandler(db)
	customers := handler.NewCustomersHandler(db)
	rentals := handler.NewRentalsHandler(db)

	api.SetupRoutes(fiberServer, games, customers, rentals)

	return fiberServer
}
