package api

import (
	"boardcamp/api/handler"

	"github.com/gofiber/fiber/v2"
)

func HealthMessage(c *fiber.Ctx) error {
	return c.SendString("I'm okay :)")
}

func SetupRoutes(app *fiber.App, games handler.GamesHandler, customer handler.CustomersHandler, rental handler.RentalsHandler) {
	app.Get("/api/health", HealthMessage)

	app.Post("/api/games", games.PostGames)
	app.Get("/api/games", games.GetGames)

	app.Post("/api/customers", customer.PostCustomers)
	app.Get("/api/customers", customer.GetCustomers)
	app.Get("/api/customers/:id", customer.GetCustomerById)
	app.Put("/api/customers/:id", customer.UpdateCustomerById)

	app.Post("/api/rentals", rental.PostRentals)
	app.Get("/api/rentals", rental.GetRentals)
	app.Post("/api/rentals/:id/return", rental.FinalizeRentalById)
	app.Delete("/api/rentals/:id", rental.DeleteRentalById)

}
