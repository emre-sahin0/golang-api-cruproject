package main

import (
	"go-rest-api/database"
	"go-rest-api/handlers"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
	// Veritabanını başlat
	database.Connect()

	e := echo.New()

	// Route'lar
	e.GET("/", handlers.RenderTodosPage)
	e.POST("/todos", handlers.AddTodo)
	e.POST("/todos/delete/:id", handlers.DeleteTodo)
	e.POST("/todos/update/:id", handlers.UpdateTodo)
	e.POST("/todos/complete/:id", handlers.MarkComplete)

	// Sunucuyu başlat
	log.Println("Server started at http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
