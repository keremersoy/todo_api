package routes

import (
	"net/http"
	"rest_api/app"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func addTodo(context *fiber.Ctx) error {
	var newTodo todo

	err := context.BodyParser(&newTodo)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Request failed"})
		return err
	}

	newTodo.ID = int(uuid.New().ID())
	if _, err := app.DB.NewInsert().Model(&newTodo).Exec(context.Context()); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not create todo"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Todo has been added"})
	return nil
}

func getTodos(context *fiber.Ctx) error {
	todos := []todo{}

	if err := app.DB.NewSelect().Model(&todos).Scan(context.Context()); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get todos"})
		return err
	}

	context.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "Todos fetched successfully", "data": todos})
	return nil
}

func removeTodo(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	if _, err := app.DB.NewDelete().Model((*todo)(nil)).Where("id = ?", id).Exec(context.Context()); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not delete todo"})
		return err

	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Todo deleted successfully"})
	return nil
}

func toggleTodoStatus(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	var todoModel = todo{}
	if err := app.DB.NewSelect().Model(&todoModel).Where("id = ?", id).Scan(context.Context()); err != nil {
		context.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "Todo not found"})
		return err
	}

	todoModel.Completed = !todoModel.Completed
	if _, err := app.DB.NewUpdate().Model(&todoModel).Where("id = ?", id).Exec(context.Context()); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not toggle todo status"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Todo status toggled successfully"})
	return nil
}

func getTodo(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "ID cannot be empty"})
		return nil
	}

	var todoModel todo
	if err := app.DB.NewSelect().Where("id = ?", id).Model(&todoModel).Scan(context.Context()); err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get the todo"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Todo fetched successfully", "data": todoModel})
	return nil
}

func SetupTodoRestrictedRoutes(app *fiber.App) {
	api := app.Group("/todos")
	api.Get("/", getTodos)
	api.Post("/", addTodo)
	api.Get("/:id", getTodo)
	api.Patch("/:id", toggleTodoStatus)
	api.Delete("/:id", removeTodo)
}
