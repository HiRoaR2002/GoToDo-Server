package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
	Body  string `json:"body"`
}

var todos []Todo
var dataFile = "todos.json"

func main() {
	fmt.Println("Hello from server")

	app := fiber.New()

	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allows all origins
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Load todos from the JSON file
	err := loadTodos()
	if err != nil {
		log.Fatalf("Error loading todos: %v", err)
	}

	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		err := c.BodyParser(todo)
		if err != nil {
			return err
		}
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)
		err = saveTodos()
		if err != nil {
			return c.Status(500).SendString("Error saving todos")
		}
		return c.JSON(todos)
	})

	app.Patch("/api/todos/:id/done", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(401).SendString("Invalid ID")
		}

		for i, t := range todos {
			if t.ID == id {
				todos[i].Done = true
				break
			}
		}
		err = saveTodos()
		if err != nil {
			return c.Status(500).SendString("Error saving todos")
		}
		return c.JSON(todos)
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.JSON(todos)
	})

	log.Fatal(app.Listen(":4000"))
}

func loadTodos() error {
	// Check if the file exists
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		// If the file does not exist, initialize with an empty slice
		todos = []Todo{}
		return nil
	}

	file, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return err
	}

	// If the file is empty, initialize with an empty slice
	if len(file) == 0 {
		todos = []Todo{}
		return nil
	}

	err = json.Unmarshal(file, &todos)
	if err != nil {
		return err
	}

	return nil
}

func saveTodos() error {
	file, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dataFile, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
