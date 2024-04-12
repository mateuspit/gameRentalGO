package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	ErrGameAlreadyRented      = "Game not found in stock"
	ErrRentalNotFinalized     = "Rental not finalized"
	ErrRentalNotFound         = "Rental not found"
	ErrGameNotFound           = "Game not found"
	ErrRentalAlreadyFinalized = "Rental already finalized"
)

type Rental struct {
	gorm.Model
	CustomerId    uint   `json:"customer_id" validate:"gt=0"`
	GameId        uint   `json:"game_id" validate:"gt=0"`
	RentDate      string `json:"rent_date"`
	DaysRented    uint   `json:"days_rented" validate:"required,gt=0"`
	ReturnDate    string `json:"return_date"`
	OriginalPrice uint   `json:"original_price"`
	DelayFee      uint   `json:"delay_fee"`
}

var validateRental = validator.New()

type RentalsHandler interface {
	GetRentals(c *fiber.Ctx) error
	PostRentals(c *fiber.Ctx) error
	FinalizeRentalById(c *fiber.Ctx) error
	DeleteRentalById(c *fiber.Ctx) error
}

type DbRentalsInstance struct {
	db *gorm.DB
}

func NewRentalsHandler(_db *gorm.DB) RentalsHandler {
	return &DbRentalsInstance{db: _db}
}

func (db *DbRentalsInstance) GetRentals(c *fiber.Ctx) error {
	rentals := []Rental{}

	allRentalsReturn := db.db.Find(&rentals)

	if allRentalsReturn.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(allRentalsReturn.Error)
	}

	return c.Status(fiber.StatusOK).JSON(&rentals)
}

func (db *DbRentalsInstance) DeleteRentalById(c *fiber.Ctx) error {
	rental := Rental{}

	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(err.Error())
	}

	if err := db.db.Where("id = ?", id).First(&rental).Error; err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(ErrRentalNotFound)
	}

	if rental.RentDate == "" {
		return c.Status(fiber.StatusConflict).JSON(ErrRentalNotFinalized)
	}

	db.db.Delete(&rental)

	return c.SendStatus(fiber.StatusOK)
}

func calculateDelayFee(delayDays float64, originalPrice uint) uint {
	const delayFeeRate = 0.01 // Taxa de atraso de 1%
	return uint(delayDays * delayFeeRate * float64(originalPrice))
}

func (db *DbRentalsInstance) FinalizeRentalById(c *fiber.Ctx) error {
	rental := Rental{}

	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(err.Error())
	}

	if err := db.db.Where("id = ?", id).First(&rental).Error; err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(ErrRentalNotFound)
	}

	if rental.ReturnDate != "" {
		return c.Status(fiber.StatusConflict).JSON(ErrRentalAlreadyFinalized)
	}

	rental.ReturnDate = time.Now().Format("2006-01-02")

	rentDate, _ := time.Parse("2006-01-02", rental.RentDate)
	expectedReturnDate := rentDate.AddDate(0, 0, int(rental.DaysRented))

	if time.Now().After(expectedReturnDate) {

		delayDays := time.Now().Sub(expectedReturnDate).Hours() / 24
		rental.DelayFee = calculateDelayFee(delayDays, rental.OriginalPrice)

	} else {
		rental.DelayFee = 0
	}

	db.db.Save(&rental)

	return c.Status(fiber.StatusOK).JSON(&rental)
}

func (db *DbRentalsInstance) PostRentals(c *fiber.Ctx) error {
	var rental Rental
	var game Game
	var customer Customer

	if err := c.BodyParser(&rental); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	validationErrors := validateRental.Struct(&rental)
	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			validationErrors.Error(),
		)
	}

	if err := db.db.Where("id = ?", rental.GameId).First(&game).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrGameNotFound)
	}

	if err := db.db.Where("id = ?", rental.CustomerId).First(&customer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrCustomerNotFound)
	}

	if game.StockTotal == 0 {
		return c.Status(fiber.StatusConflict).JSON(ErrGameAlreadyRented)
	}

	rental.OriginalPrice = game.PricePerDay * rental.DaysRented
	rental.RentDate = time.Now().Format("2006-01-02")

	db.db.Create(&rental)

	game.StockTotal -= 1
	db.db.Save(&game)

	return c.Status(fiber.StatusCreated).JSON(&rental)
}
