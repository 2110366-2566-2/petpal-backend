package utills

import (
	"context"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	// "github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func SearchServices(db *models.MongoDB, searchHistory *models.SearchHistory, id string, is_user bool) (*[]models.Service, error) {
	// get collection
	svcp_collection := db.Collection("svcp")

	// Q := searchHistory.Q
	// Location := searchHistory.Location
	// StartTime := searchHistory.StartTime
	// EndTime := searchHistory.EndTime
	// StartPriceRange := searchHistory.StartPriceRange
	// EndPriceRange := searchHistory.EndPriceRange
	// MinRating := searchHistory.MinRating
	// MaxRating := searchHistory.MaxRating
	// PageNumber := searchHistory.PageNumber
	// PageSize := searchHistory.PageSize
	// SortBy := searchHistory.SortBy

	// find user by email
	filter := bson.D{
		// {Key: "services.service_name", Value: bson.D{{Key: "$regex", Value: bookingCreate.Q}}},
		// {Key: "location", Value: bookingCreate.Location},
		// {Key: "services.start_time", Value: bson.D{{Key: "$gte", Value: bookingCreate.StartTime}}},
		// {Key: "services.end_time", Value: bson.D{{Key: "$lte", Value: bookingCreate.EndTime}}},
		// {Key: "services.price", Value: bson.D{{Key: "$gte", Value: bookingCreate.StartPriceRange}, {Key: "$lte", Value: bookingCreate.EndPriceRange}}},
		// {Key: "services.rating", Value: bson.D{{Key: "$gte", Value: bookingCreate.MinRating}, {Key: "$lte", Value: bookingCreate.MaxRating}}},
	}
	cursor, err := svcp_collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var services []models.Service
	if err = cursor.All(context.Background(), &services); err != nil {
		return nil, err
	}

	return &services, nil
}
