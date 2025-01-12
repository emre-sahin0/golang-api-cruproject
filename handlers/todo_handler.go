package handlers

import (
	"context"
	"go-rest-api/database"
	"go-rest-api/models"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddTodo(c echo.Context) error {
	title := c.FormValue("title")
	if title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Title is required"})
	}

	// Yeni bir Todo oluştur
	todo := models.Todo{
		ID:        primitive.NewObjectID(),
		Title:     title,
		Completed: false,
	}

	// MongoDB'ye kaydet
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := database.TodoCollection.InsertOne(ctx, todo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save todo"})
	}

	return c.Redirect(http.StatusSeeOther, "/") // Ana sayfaya yönlendir
}
func DeleteTodo(c echo.Context) error {
	id := c.Param("id")

	// ID'yi kontrol edin
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID is required"})
	}
	log.Println("Received ID:", id)
	// ObjectID'ye dönüştür
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID format"})
	}

	// MongoDB'den sil
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}
	result, err := database.TodoCollection.DeleteOne(ctx, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete todo"})
	}

	// Silinen belge sayısını kontrol edin
	if result.DeletedCount == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Todo not found"})
	}

	return c.Redirect(http.StatusSeeOther, "/") // Ana sayfaya yönlendir
}
func UpdateTodo(c echo.Context) error {
	id := c.Param("id")
	title := c.FormValue("title")
	completed := c.FormValue("completed") == "on"

	// ID kontrolü ve ObjectID'ye dönüştürme
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	// MongoDB'deki Todo'yu güncelle
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"title":     title,
			"completed": completed,
		},
	}

	_, err = database.TodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update todo"})
	}

	return c.Redirect(http.StatusSeeOther, "/") // Ana sayfaya yönlendir
}
func MarkComplete(c echo.Context) error {
	id := c.Param("id")

	// ID kontrolü ve ObjectID'ye dönüştürme
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
	}

	// MongoDB'deki Todo'yu güncelle
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"completed": true,
		},
	}

	_, err = database.TodoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to mark as complete"})
	}

	return c.Redirect(http.StatusSeeOther, "/") // Ana sayfaya yönlendir
}
func RenderTodosPage(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.TodoCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Database error"})
	}
	defer cursor.Close(ctx)

	var todos []models.Todo
	if err := cursor.All(ctx, &todos); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Data conversion error"})
	}

	// Todo'ları ayır
	var pendingTodos, completedTodos []models.Todo
	for _, todo := range todos {
		if todo.Completed {
			completedTodos = append(completedTodos, todo)
		} else {
			pendingTodos = append(pendingTodos, todo)
		}
	}

	data := map[string]interface{}{
		"PendingTodos":   pendingTodos,
		"CompletedTodos": completedTodos,
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Template error"})
	}
	return tmpl.Execute(c.Response().Writer, data)
}
