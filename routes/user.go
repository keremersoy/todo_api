package routes

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"rest_api/app"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type user struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

var (
	PrivateKey *rsa.PrivateKey
)

func GenerateKey() {
	rng := rand.Reader
	var err error
	PrivateKey, err = rsa.GenerateKey(rng, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}
}

func userLogin(context *fiber.Ctx) error {
	var _user user

	err := context.BodyParser(&_user)
	if _user.Email == "" || _user.Password == "" || err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Request failed"})
		return err
	}

	var userModel user
	if err := app.DB.NewSelect().Where("email = ?", _user.Email).Model(&userModel).Scan(context.Context()); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "User was not found"})
		return err
	}

	if _user.Password != userModel.Password {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Wrong password"})
		return err
	}

	claims := jwt.MapClaims{
		"user_id": userModel.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	t, err := token.SignedString(PrivateKey)
	if err != nil {
		context.Status(fiber.StatusInternalServerError)
		return err
	}

	context.Status(http.StatusOK).JSON(fiber.Map{"token": t})
	return nil
}

func userRegister(context *fiber.Ctx) error {
	var newUser user

	err := context.BodyParser(&newUser)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Request failed"})
		return err
	}

	newUser.ID = int(uuid.New().ID())
	if _, err := app.DB.NewInsert().Model(&newUser).Exec(context.Context()); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not register user"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User has been register"})
	return nil
}

func getUserList(context *fiber.Ctx) error {
	users := []user{}

	if err := app.DB.NewSelect().Model(&users).Scan(context.Context()); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get users"})
		return err
	}

	context.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Users fetched successfully", "data": users})
	return nil
}

func getUser(context *fiber.Ctx) error {
	currentUser := context.Locals("user").(*jwt.Token)
	claims := currentUser.Claims.(jwt.MapClaims)
	id := claims["user_id"].(float64)
	fmt.Println(id)
	var userModel user

	if err := app.DB.NewSelect().Where("id = ?", id).Model(&userModel).Scan(context.Context()); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "User was not found"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "User fetched successfully", "data": userModel})
	return nil
}

func SetupUserRoutes(app *fiber.App) {
	api := app.Group("/user")
	api.Post("/login", userLogin)
	api.Post("/register", userRegister)
}

func SetupUserRestrictedRoutes(app *fiber.App) {
	api := app.Group("/user")
	api.Get("/", getUserList)
	api.Get("/:info", getUser)
}
