package utills

import (
	"context"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// search result struct

func SearchServices(db *models.MongoDB, searchHistory *models.SearchHistory, id string, is_user bool) ([]models.SearchResult, error) {
	// get collection
	var sortCiteria string
	var isDesc int32
	svcp_collection := db.Collection("svcp")

	if searchHistory.SortBy == "price" {
		sortCiteria = "services.price"
	} else if searchHistory.SortBy == "rating" {
		sortCiteria = "services.averageRating"
	} else {
		sortCiteria = "services.serviceName"
	}
	if searchHistory.Descending {
		isDesc = 1
	} else {
		isDesc = -1
	}

	// Define the filter to find the documents
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$services"}},
		{{Key: "$match", Value: bson.D{
			{Key: "$and", Value: bson.A{
				bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "services.serviceName", Value: bson.D{{Key: "$regex", Value: searchHistory.Q}, {Key: "$options", Value: "i"}}}},
					bson.D{{Key: "services.serviceDescription", Value: bson.D{{Key: "$regex", Value: searchHistory.Q}, {Key: "$options", Value: "i"}}}},
				}}},
				bson.D{{Key: "services.serviceType", Value: bson.D{{Key: "$regex", Value: searchHistory.ServicesType}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "location", Value: bson.D{{Key: "$regex", Value: searchHistory.Location}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "services.price", Value: bson.D{{Key: "$gte", Value: searchHistory.StartPriceRange}, {Key: "$lte", Value: searchHistory.EndPriceRange}}}},
				bson.D{{Key: "services.averageRating", Value: bson.D{{Key: "$gte", Value: searchHistory.MinRating}, {Key: "$lte", Value: searchHistory.MaxRating}}}},
			}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: sortCiteria, Value: isDesc}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "services", Value: 1},
			{Key: "location", Value: 1},
			{Key: "SVCPUsername", Value: 1},
			{Key: "SVCPServiceType", Value: 1},
		}}},
	}

	// Run the aggregation
	cursor, err := svcp_collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	// Decode the documents
	var results []models.SearchResult
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return results, nil
}
