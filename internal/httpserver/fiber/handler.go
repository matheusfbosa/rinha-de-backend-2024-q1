package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matheusfbosa/rinha-de-backend-2024-q1/customer"
)

func Handlers(service customer.UseCase) *fiber.App {
	app := fiber.New()
	app.Post("/clientes/:id/transacoes", MakeTransaction(service))
	app.Get("/clientes/:id/extrato", GetStatement(service))
	return app
}

func MakeTransaction(s customer.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tr customer.Transaction
		if err := c.BodyParser(&tr); err != nil {
			return c.Status(fiber.StatusBadRequest).
				SendString(customer.ErrInvalidTransaction.Error())
		}

		tr.CustomerID = c.Params("id")
		err := tr.Validate()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		balance, err := s.MakeTransaction(&tr)
		if err != nil {
			return handleError(c, err)
		}

		return c.JSON(balance)
	}
}

func GetStatement(s customer.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		customerID := c.Params("id")
		statement, err := s.GetBankStatement(customerID)
		if err != nil {
			return handleError(c, err)
		}

		return c.JSON(statement)
	}
}

func handleError(c *fiber.Ctx, err error) error {
	switch err {
	case customer.ErrInvalidTransaction:
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	case customer.ErrCustomerNotFound:
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	case customer.ErrInsufficientFunds:
		return c.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	default:
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
}
