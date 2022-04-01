package token

import (
	"context"
	"os"
	"shared"
	"time"

	"auth_api/users"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type GetResponseData struct {
	JWT string `json:"jwt,omitempty"`
}

func Get(db *mongo.Database) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		res := new(shared.Response[GetResponseData])

		body := new(Credential)
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

		user := new(users.User)
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

		payload := NewPayload(user.ID, time.Hour * 24 * 7)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		res.Data.JWT, err = token.SignedString([]byte(secret))
		if err != nil {
			return c.SendStatus(500)
		}

		return res.Send(c)
	}
}
