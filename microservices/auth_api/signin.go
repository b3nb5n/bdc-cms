package main

import (
	"context"
	"os"
	"shared"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type SigninBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SigninResponseData struct {
	JWT string `json:"jwt,omitempty"`
}

func Signin(db *mongo.Database) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res := new(shared.Response[SigninResponseData])

		body := new(SigninBody)
		c.BodyParser(&body)
		err := validate.Struct(body)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				return c.SendStatus(500)
			}

			return c.SendStatus(400)
		}

		queryCtx, cancelQueryCtx := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancelQueryCtx()
		queryResult := db.Collection("users").FindOne(queryCtx, bson.M{"data.email": body.Email})
		if err = queryResult.Err(); err != nil {
			switch err {
			case mongo.ErrNoDocuments:
				res.Error = "No accounts associated with this email"
				return res.Send(c.Status(404))
			default:
				return c.SendStatus(500)
			}
		}

		user := new(User)
		err = queryResult.Decode(&user)
		if err != nil {
			return c.SendStatus(500)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Data.Password), []byte(body.Password))
		if err != nil {
			res.Error = "Wrong password"
			return res.Send(c.Status(401))
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			res.Send(c.Status(500))
			panic("No configured jwt secret")
		}

		token := jwt.New(jwt.SigningMethodHS256)
		res.Data.JWT, err = token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(500)
		}

		return res.Send(c)
	}
}
