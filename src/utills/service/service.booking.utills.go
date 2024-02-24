package utills

import (
	"context"
	"petpal-backend/src/models"
	// "github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func InsertBooking(db *models.MongoDB, booking *models.Booking) (*models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Insert the booking into the collection
	_, err := collection.InsertOne(context.Background(), booking)
	if err != nil {
		return nil, err
	}

	// Return the inserted booking
	return booking, nil
}
