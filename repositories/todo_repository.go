package repositories

import (
	"context"
	"go-rest-api/database"
	"go-rest-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllTodos() ([]models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.TodoCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []models.Todo
	if err = cursor.All(ctx, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func AddTodo(todo models.Todo) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.TodoCollection.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)
	return &todo, nil
}

func DeleteTodo(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}

	_, err := database.TodoCollection.DeleteOne(ctx, filter)
	return err
}
