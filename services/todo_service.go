package services

import (
	"go-rest-api/models"
	"go-rest-api/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTodos() ([]models.Todo, error) {
	return repositories.GetAllTodos()
}

func CreateTodo(todo models.Todo) (*models.Todo, error) {
	return repositories.AddTodo(todo)
}

func RemoveTodo(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return repositories.DeleteTodo(objectID)
}
