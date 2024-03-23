package utills

import (
	"context"
	"log"
	"petpal-backend/src/configs"
	"petpal-backend/src/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDB creates a new MongoDB instance with the provided connection string.
func NewMongoDB() (*models.MongoDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//"mongodb://inwza:strongpassword@localhost:27017/"
	mongoClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(configs.GetInstance().GetDB_URI()),
	)

	if err != nil {
		log.Fatalf("connection error :%v", err)
		return nil, err
	}

	return &models.MongoDB{Client: mongoClient, DbName: "petpal"}, nil
}

// NewMongoDB creates a new MongoDB instance with the provided connection string.
func NewMongoDBTest() (*models.MongoDB, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//"mongodb://inwza:strongpassword@localhost:27017/"
	mongoClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(configs.GetInstance().GetDB_URI()),
	)


	if err != nil {
		log.Fatalf("connection error :%v", err)
		return nil, err
	}

	return &models.MongoDB{Client: mongoClient, DbName: "test"}, nil
}
