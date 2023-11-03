package main

import (
	"context"
	"encoding/json"
	"fmt"
	"rest_api/app"
	model "rest_api/models"
	"rest_api/routes"
	"rest_api/storage"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config := storage.NewConfig()

	db := storage.NewConnection(&config)
	defer db.Close()

	app.Initialise(db)

	ctx := context.Background()

	if !fiber.IsChild() {

		conf_str, _ := json.MarshalIndent(config, "", " ")
		fmt.Println(string(conf_str))

		model.Initialise(db, ctx)
	}
	app := fiber.New()
	app.Use(logger.New())

	routes.GenerateKey()
	routes.SetupUserRoutes(app)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    routes.PrivateKey.Public(),
		},
	}))

	routes.SetupTodoRestrictedRoutes(app)
	routes.SetupUserRestrictedRoutes(app)

	app.Listen("localhost:9090")
}
