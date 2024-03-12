package utills

import (
	"context"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func SearchServices(db *models.MongoDB, searchHistory *models.SearchHistory, id string, is_user bool) (*[]bson.M, error) {
	// get collection
	svcp_collection := db.Collection("svcp")
	// Collection of SVCP is like this
	// {
	// 	"IndividualID": "65f05a3c4dccea9ed58bb99b",
	// 	"SVCPID": "65f05a3c4dccea9ed58bb99b",
	// 	"SVCPImg": "",
	// 	"SVCPUsername": "test5@gmail.com",
	// 	"SVCPPassword": "$2a$10$cHBAnisxdc0o26xUUk0L8u3U3XYUdVi7FHdBISvCBQyuZHHjTurHS",
	// 	"SVCPEmail": "test5@gmail.com",
	// 	"isVerified": false,
	// 	"SVCPResponsiblePerson": "",
	// 	"defaultBank": "",
	// 	"defaultAccountNumber": "",
	// 	"license": "",
	// 	"location": "",
	// 	"description": "",
	// 	"SVCPAdditionalImg": "",
	// 	"SVCPServiceType": "E",
	// 	"services": [
	// 	  {
	// 		"serviceID": "65f05b384dccea9ed58bb99d",
	// 		"serviceName": "WWWWWWW",
	// 		"serviceType": "E",
	// 		"serviceDescription": "ASDASD",
	// 		"serviceImg": "",
	// 		"averageRating": 0,
	// 		"requireCert": false,
	// 		"timeslots": [],
	// 		"price": 5
	// 	  },
	// 	  {
	// 		"serviceID": "65f05b394dccea9ed58bb99e",
	// 		"serviceName": "WWWWWWW",
	// 		"serviceType": "E",
	// 		"serviceDescription": "ASDASD",
	// 		"serviceImg": "",
	// 		"averageRating": 0,
	// 		"requireCert": false,
	// 		"timeslots": [],
	// 		"price": 5
	// 	  },
	// 	  {
	// 		"serviceID": "65f05b3b4dccea9ed58bb99f",
	// 		"serviceName": "WWWWWWW",
	// 		"serviceType": "E",
	// 		"serviceDescription": "ASDASD",
	// 		"serviceImg": "",
	// 		"averageRating": 0,
	// 		"requireCert": false,
	// 		"timeslots": [],
	// 		"price": 5
	// 	  },
	// 	  {
	// 		"serviceID": "65f05b3d4dccea9ed58bb9a0",
	// 		"serviceName": "WWWWWWW",
	// 		"serviceType": "E",
	// 		"serviceDescription": "ASDASD",
	// 		"serviceImg": "",
	// 		"averageRating": 0,
	// 		"requireCert": false,
	// 		"timeslots": [],
	// 		"price": 5
	// 	  },
	// 	  {
	// 		"serviceID": "65f05b3e4dccea9ed58bb9a1",
	// 		"serviceName": "WWWWWWW",
	// 		"serviceType": "E",
	// 		"serviceDescription": "ASDASD",
	// 		"serviceImg": "",
	// 		"averageRating": 0,
	// 		"requireCert": false,
	// 		"timeslots": [],
	// 		"price": 5
	// 	  }
	// 	]
	//   }

	// find user by email
	filter := bson.D{
		// {Key: "services.service_name", Value: bson.D{{Key: "$regex", Value: searchHistory.Q}}},
		// {Key: "services.service_type", Value: bson.D{{Key: "$regex", Value: searchHistory.Q}}},
		{Key: "location", Value: bson.D{{Key: "$regex", Value: searchHistory.Location}}},
		// {Key: "services.timeslots.start_time", Value: bson.D{{Key: "$gte", Value: searchHistory.StartTime}}},
		// {Key: "services.timeslots.end_time", Value: bson.D{{Key: "$lte", Value: searchHistory.EndTime}}},
		// {Key: "services.price", Value: bson.D{{Key: "$gte", Value: searchHistory.StartPriceRange}, {Key: "$lte", Value: searchHistory.EndPriceRange}}},
		// {Key: "services.rating", Value: bson.D{{Key: "$gte", Value: searchHistory.MinRating}, {Key: "$lte", Value: searchHistory.MaxRating}}},
	}

	// Define the options to sort and paginate the results
	opts := options.Find()
	if searchHistory.SortBy == "timeslots" {
		opts.SetSort(bson.D{{Key: "services.timeslots.start_time", Value: 1}})
	}
	opts.SetSort(bson.D{{Key: "services." + searchHistory.SortBy, Value: 1}})
	opts.SetSkip(int64((searchHistory.PageNumber - 1) * searchHistory.PageSize))
	opts.SetLimit(int64(searchHistory.PageSize))
	cursor, err := svcp_collection.Find(context.Background(), filter, opts)
	if err != nil {
		return nil, err
	}

	// Define a variable to store the result
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}
