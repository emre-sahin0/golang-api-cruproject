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
	e.GET("/", handlers.RenderTodosPage)                 // Ana sayfa
	e.POST("/todos", handlers.AddTodo)                   // Yeni todo ekleme
	e.POST("/todos/delete/:id", handlers.DeleteTodo)     // Todo silme
	e.POST("/todos/update/:id", handlers.UpdateTodo)     // Todo güncelleme
	e.POST("/todos/complete/:id", handlers.MarkComplete) // Todo'yu tamamlama

	// Sunucuyu başlat
	log.Println("Server started at http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
