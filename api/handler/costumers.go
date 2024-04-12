package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	ErrCustomerCPFExists = "Customer already exists"
	ErrCustomerNotFound  = "Customer not found"
)

type Customer struct {
	gorm.Model
	Name     string `json:"name" validate:"required,min=2"`
	Phone    string `json:"phone" validate:"min=10,max=11"`
	CPF      string `json:"cpf" validate:"len=11"`
	Birthday string `json:"birthday" validate:"required,datetime=2006-01-02"`
}

var validateCustomer = validator.New()

type CustomersHandler interface {
	GetCustomers(c *fiber.Ctx) error
	GetCustomerById(c *fiber.Ctx) error
	PostCustomers(c *fiber.Ctx) error
	UpdateCustomerById(c *fiber.Ctx) error
}

type DbCustomersInstance struct {
	db *gorm.DB
}

func NewCustomersHandler(_db *gorm.DB) CustomersHandler {
	return &DbCustomersInstance{db: _db}
}

func (db *DbCustomersInstance) GetCustomers(c *fiber.Ctx) error {
	customers := []Customer{}

	allCustomersReturn := db.db.Find(&customers)

	if allCustomersReturn.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(allCustomersReturn.Error)
	}

	return c.Status(fiber.StatusOK).JSON(&customers)
}

func (db *DbCustomersInstance) UpdateCustomerById(c *fiber.Ctx) error {
	customer := Customer{}
	customerUpdate := Customer{}

	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(err.Error())
	}

	if err := db.db.Where("id = ?", id).First(&customer).Error; err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(ErrCustomerNotFound)
	}

	if err := c.BodyParser(&customerUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	validationErrors := validateCustomer.Struct(&customerUpdate)
	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			validationErrors.Error(),
		)
	}

	var existingCustomer Customer
	if err := db.db.Where("cpf = ? AND id != ?", customerUpdate.CPF, id).First(&existingCustomer).Error; err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusConflict).JSON(ErrCustomerCPFExists)
	}

	customer.Name = customerUpdate.Name
	customer.CPF = customerUpdate.CPF
	customer.Phone = customerUpdate.Phone
	customer.Birthday = customerUpdate.Birthday

	db.db.Save(&customer)

	return c.Status(fiber.StatusOK).JSON(customer)
}

func (db *DbCustomersInstance) GetCustomerById(c *fiber.Ctx) error {
	customer := Customer{}

	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(err.Error())
	}

	if err := db.db.Where("id = ?", id).First(&customer).Error; err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(ErrCustomerCPFExists)
	}

	return c.Status(fiber.StatusOK).JSON(customer)

}

func (db *DbCustomersInstance) PostCustomers(c *fiber.Ctx) error {
	var customer Customer

	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	validationErrors := validateCustomer.Struct(&customer)
	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			validationErrors.Error(),
		)
	}

	var existingCustomer Customer
	if err := db.db.Where("cpf = ?", customer.CPF).First(&existingCustomer).Error; err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusConflict).JSON(ErrCustomerCPFExists)
	}

	db.db.Create(&customer)

	return c.Status(fiber.StatusCreated).JSON(&customer)
}
