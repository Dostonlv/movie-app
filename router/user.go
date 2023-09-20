package router

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"movie/database"
	"movie/mail"
	"movie/models"
	"movie/pkg/utils"
	"os"
	"time"
)

func CreateUser(c *fiber.Ctx) error {
	c.Accepts("application/json")

	collectionUser := database.InitDB().Db.Collection("user")
	collectionEmailVerificationToken := database.InitDB().Db.Collection("email_verification_token")
	user := new(models.User)
	emailVerifyToken := new(models.EmailVerificationToken)
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
	err = collectionUser.FindOne(c.Context(), filter).Decode(&result)
	if err == nil {
		return c.Status(401).SendString("email already in use")
	}

	insertionResult, err := collectionUser.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	filter = bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := collectionUser.FindOne(c.Context(), filter)

	createdUser := &models.User{}
	createdRecord.Decode(createdUser)
	OTP := utils.GenerateOTP()

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	from := os.Getenv("NAMECHEAP_EMAIL")
	password := os.Getenv("NAMECHEAP_PASSWORD")

	OTPsender := mail.NewEmailSender("to-kioname", from, password)
	to := []string{user.Email}
	err = OTPsender.SendEmail("Hello "+user.Name, "Hello Your OTP code: "+OTP, to, nil, nil, nil)
	if err != nil {
		return c.Status(500).SendString("OTP not sent to your email")
	}
	emailVerifyToken.OwnerID = insertionResult.InsertedID
	hashedOTP, err := utils.Hash(OTP)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	emailVerifyToken.Token = hashedOTP
	emailVerifyToken.CreatedAt = time.Now()

	opts := options.Index().SetExpireAfterSeconds(60)
	_, err = collectionEmailVerificationToken.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{"Expires", 1}},
		Options: opts,
	})
	if err != nil {
		return c.Status(500).SendString("index not created")
	}
	emailVerifyToken.Expires = emailVerifyToken.CreatedAt.Add(1 * time.Minute)

	_, err = collectionEmailVerificationToken.InsertOne(c.Context(), emailVerifyToken)
	if err != nil {
		return c.Status(500).SendString("OTP not sent to your email")
	}

	return c.Status(201).SendString("OTP sent to your email")

}

func VerifyEmail(c *fiber.Ctx) error {
	c.Accepts("application/json")

	collectionUser := database.InitDB().Db.Collection("user")
	collectionEmailVerificationToken := database.InitDB().Db.Collection("email_verification_token")
	emailVerify := new(models.EmailVerify)
	if err := c.BodyParser(emailVerify); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// check email otp expired email_verification_token
	var expired models.EmailVerificationToken
	filter := bson.D{{Key: "_id", Value: emailVerify.ID}}
	err := collectionEmailVerificationToken.FindOne(c.Context(), filter).Decode(&expired)
	if err != nil {
		return c.Status(500).SendString("error from  expired decode")
	}

	// find user by id
	var result models.User
	err = collectionUser.FindOne(c.Context(), filter).Decode(&result)
	if err != nil {
		return c.Status(401).SendString("user not found")
	}

	// find email verification token by id
	filter = bson.D{{Key: "_id", Value: emailVerify.ID}}
	var resultEmailVerificationToken models.EmailVerificationToken
	err = collectionEmailVerificationToken.FindOne(c.Context(), filter).Decode(&resultEmailVerificationToken)
	if err != nil {
		return c.Status(401).SendString("user not found")
	}

	// compare otp
	err = utils.CompareHash(resultEmailVerificationToken.Token, emailVerify.OTP)
	if err != nil {
		return c.Status(401).SendString("wrong otp")
	}

	// if otp is expired user is_verify not update

	// check expired.Expires > time.Now() update else otp expired
	if expired.Expires.After(time.Now()) {
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_verify", Value: true}}}}
		_, err = collectionUser.UpdateOne(c.Context(), filter, update)
		if err != nil {
			return c.Status(500).SendString("something went wrong")
		}
	} else {
		return c.Status(500).SendString("OTP expired")
	}

	return c.Status(200).SendString("email verified")

}
