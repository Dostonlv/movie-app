package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"movie/database"
	"movie/middleware"
	"movie/models"
	"movie/router"
)

func main() {
	app := fiber.New()

	api := app.Group("/api", logger.New())
	user := api.Group("/user")
	user.Post("/create", middleware.UserValidator, router.User)
	database.InitDB()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/sign-in", func(c *fiber.Ctx) error {
		body := new(models.SignIn)
		c.BodyParser(&body)
		if body.Email == "" || body.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Please enter all fields")
		}

		return c.Status(fiber.StatusOK).JSON("Success!")
	})

	log.Fatal(app.Listen(":3000"))
}
