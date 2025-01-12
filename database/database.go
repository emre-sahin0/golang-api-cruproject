package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var TodoCollection *mongo.Collection

func Connect() {
	// MongoDB'ye bağlan
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MongoDB'ye bağlanılamadı:", err)
	}

	TodoCollection = Client.Database("go_rest_api_db").Collection("todos")

	log.Println("MongoDB'ye başarıyla bağlanıldı!")
}
