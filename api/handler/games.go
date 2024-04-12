package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	ErrGameNameExists = "Name already exists"
)

type DbGamesInstance struct {
	db *gorm.DB
}

type Game struct {
	gorm.Model
	Name        string `json:"name" validate:"required,min=2"`
	Image       string `json:"image"`
	StockTotal  uint   `json:"stock_total" validate:"gt=0"`
	PricePerDay uint   `json:"price_per_day" validate:"gt=0"`
}

var validateGame = validator.New()

type GamesHandler interface {
	GetGames(c *fiber.Ctx) error
	PostGames(c *fiber.Ctx) error
}

func NewGamesHandler(_db *gorm.DB) GamesHandler {
	return &DbGamesInstance{db: _db}
}

func (db *DbGamesInstance) GetGames(c *fiber.Ctx) error {
	games := []Game{}

	allGamesReturn := db.db.Find(&games)

	if allGamesReturn.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(allGamesReturn.Error)
	}

	return c.Status(fiber.StatusOK).JSON(games)
}

func (db *DbGamesInstance) PostGames(c *fiber.Ctx) error {
	var game Game

	if err := c.BodyParser(&game); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	validationErrors := validateGame.Struct(&game)
	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			validationErrors.Error(),
		)
	}

	var existingGame Game
	if err := db.db.Where("name = ?", game.Name).First(&existingGame).Error; err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusConflict).JSON(ErrGameNameExists)
	}

	db.db.Create(&game)

	return c.Status(fiber.StatusCreated).JSON(&game)
}
