package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"movie/database"
	"movie/mail"
	"movie/models"
	"movie/pkg/utils"
	"os"
)

func User(c *fiber.Ctx) error {
	c.Accepts("application/json")

	collection := database.InitDB().Db.Collection("user")
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	hashedPassword, err := utils.Hash(user.Password)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	user.Name = utils.TrimSpace(user.Name)
	user.Email = utils.TrimSpace(user.Email)
	user.Password = hashedPassword

	filter := bson.D{{Key: "email", Value: user.Email}}
	var result models.User
	err = collection.FindOne(c.Context(), filter).Decode(&result)
	if err == nil {
		return c.Status(401).SendString("email already in use")
	}

	insertionResult, err := collection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	filter = bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdUser := &models.User{}
	createdRecord.Decode(createdUser)
	OTP := utils.GenerateOTP()

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	from := os.Getenv("NAMECHEAP_EMAIL")
	password := os.Getenv("NAMECHEAP_PASSWORD")

	OTPsender := mail.NewEmailSender("Uchiha Foundation", from, password)
	to := []string{user.Email}
	err = OTPsender.SendEmail("Hello "+user.Name, "Hello Your OTP code: "+OTP, to, nil, nil, nil)
	if err != nil {
		return c.Status(500).SendString("OTP not sent to your email")
	}

	return c.Status(201).SendString("OTP sent to your email")

}
