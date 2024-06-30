package main

import (
	"fmt"

	"log"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Server started.")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	app := fiber.New()

	todos := []Todo{}

	app.Get("/api/getTodos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo Body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)

	})

	app.Put("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	
	app.Delete("/api/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		log.Printf("Received DELETE request for ID: %s", id)
		
		for i, item := range todos {
			if fmt.Sprint(item.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				log.Printf("Deleted TODO with ID: %s", id)
				return c.Status(200).JSON(todos)
			}
		}
		log.Printf("TODO with ID: %s not found", id)
		return c.Status(404).JSON(fiber.Map{"error": "TODO not found"})
	})
	
	log.Fatal(app.Listen(":" +port))
}

