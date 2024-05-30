package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/martijnwiekens/go-learning/todo/repository"
)

type addTodoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"dueDate"`
	Done        bool   `json:"done"`
}

func main() {
	// Create the database
	log.Default().Println("Creating the database...")
	repository.Setup()

	// Create the API
	router := gin.Default()
	router.GET("/", outputPage)
	router.GET("/api/v1/todo", getTodos)
	router.POST("/api/v1/todo", addTodo)
	router.GET("/api/v1/todo/:id", getTodo)
	router.PATCH("/api/v1/todo/:id", updateTodo)
	router.DELETE("/api/v1/todo/:id", deleteTodo)

	// Start the API
	log.Default().Print("Starting the API...")
	router.Run("localhost:8081")
}

func outputPage(c *gin.Context) {
	// Return the index.html file
	c.File("./index.html")
}

func getTodos(c *gin.Context) {
	// Get the todos from the database
	res, err := repository.GetTodos()
	if err != nil {
		// Something went wrong
		log.Default().Println(err)
		c.AbortWithStatus(500)
		return
	}

	// Return JSON
	c.JSON(200, res)
}

func addTodo(c *gin.Context) {
	// Retrieve the todo from the request
	var input addTodoInput
	err := c.BindJSON(&input)
	if err != nil {
		// Something went wrong
		log.Default().Println(err)
		c.AbortWithStatus(400)
		return
	}
	fmt.Println(input)

	// Add the todo to the database
	err = repository.AddTodo(repository.Todo{
		Id:          uuid.NewString(),
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Done:        input.Done,
	})
	if err != nil {
		// Something went wrong
		log.Default().Println(err)
		c.AbortWithStatus(500)
		return
	}

	// Return JSON
	c.JSON(200, gin.H{
		"message": "Todo added",
	})
}

func getTodo(c *gin.Context) {
	// Retrieve the ID from the request
	id := c.Param("id")
	if id == "" {
		// Something went wrong
		log.Default().Println("No ID provided")
		c.AbortWithStatus(400)
		return
	}

	// Get the todo from the database
	res, err := repository.GetTodo(id)
	if err != nil {
		// Something went wrong
		log.Default().Println(err)
		c.AbortWithStatus(500)
		return
	}

	// Return JSON
	c.JSON(200, res)
}

func updateTodo(c *gin.Context) {
	// Retrieve the ID from the request
	id := c.Param("id")
	if id == "" {
		// Something went wrong
		log.Default().Println("No ID provided")
		c.AbortWithStatus(400)
		return
	}

	// Retrieve the todo from the request
	var input addTodoInput
	err := c.BindJSON(&input)
	if err != nil {
		// Something went wrong
		log.Default().Println(err)
		c.AbortWithStatus(400)
		return
	}
	fmt.Println(input)

	// Get the item from the database
	res, err := repository.GetTodo(id)
	if err != nil {
		// Something went wrong
		log.Default().Println(err)
		c.AbortWithStatus(500)
		return
	}

	// Update the todo in the database
	err = repository.UpdateTodo(repository.Todo{
		Id:          id,
		Title:       res.Title,
		Description: res.Description,
		DueDate:     res.DueDate,
		Done:        input.Done,
	})
	if err != nil {
		// Something went wrong
		log.Default().Println(err)
		c.AbortWithStatus(500)
		return
	}

	// Return JSON
	c.JSON(200, gin.H{
		"message": "Todo updated",
	})
}

func deleteTodo(c *gin.Context) {
	// Retrieve the ID from the request
	id := c.Param("id")
	if id == "" {
		// Something went wrong
		log.Default().Println("No ID provided")
		c.AbortWithStatus(400)
		return
	}

	// Delete the todo from the database
	err := repository.DeleteTodo(repository.Todo{
		Id: id,
	})
	if err != nil {
		// Something went wrong
		log.Default().Println(err)
		c.AbortWithStatus(500)
		return
	}

	// Return JSON
	c.JSON(200, gin.H{
		"message": "Todo deleted",
	})
}
