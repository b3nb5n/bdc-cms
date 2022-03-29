package main

import (
	"context"
	"shared"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

type SignupResponseData struct {
	ID shared.Snowflake `json:"id"`
}

type SignupResponseError struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

func Signup(db *mongo.Database) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res := new(shared.Response[SignupResponseData])

		data := new(UserData)
		c.BodyParser(data)
		err := validate.Struct(data)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				return c.SendStatus(500)
			}

			return c.SendStatus(400)
		}

		readCtx, cancelReadCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelReadCtx()
		readResult := db.Collection("users").FindOne(readCtx, bson.M{"data.email": data.Email})
		if readResult.Err() != mongo.ErrNoDocuments {
			return c.SendStatus(409)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.SendStatus(500)
		}

		data.Password = string(hash)
		user, err := shared.NewResource(*data)
		if err != nil {
			return c.SendStatus(500)
		}

		writeCtx, cancelWriteCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelWriteCtx()
		writeResult, err := db.Collection("users").InsertOne(writeCtx, user)
		if err != nil {
			return c.SendStatus(500)
		}

		if id, ok := writeResult.InsertedID.(int64); ok {
			res.Data.ID = shared.Snowflake(id)
			return res.Send(c)
		}

		return c.SendStatus(500)
	}
}
