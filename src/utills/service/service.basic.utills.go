package utills

import (
	"context"
	"errors"
	"fmt"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	misc_utills "petpal-backend/src/utills/miscellaneous"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateNewServices(createServices *models.CreateService) *models.Service {

	newTimeslots := []models.Timeslot{}
	for i := 0; i < len(createServices.Timeslots); i++ {
		newTimeslot := models.Timeslot{
			TimeslotID: primitive.NewObjectID().Hex(),
			StartTime:  createServices.Timeslots[i].StartTime,
			EndTime:    createServices.Timeslots[i].EndTime,
			Status:     "available",
		}
		newTimeslots = append(newTimeslots, newTimeslot)
		// newTimeslots[i] = newTimeslot
	}

	randomServiceImage, err := misc_utills.RandomServiceImage()
	if err != nil {
		return nil
	}

	return &models.Service{
		ServiceID:          primitive.NewObjectID().Hex(),
		ServiceName:        createServices.ServiceName,
		ServiceType:        createServices.ServiceType,
		ServiceDescription: createServices.ServiceDescription,
		Price:              createServices.Price,
		ServiceImg:         randomServiceImage,
		AverageRating:      0,
		RequireCert:        false,
		Timeslots:          newTimeslots,
	}
}

func AddNewServices(db *models.MongoDB, createServices *models.CreateService, svcp *models.SVCP) (*models.Service, error) {
	// get collection
	svcp_collection := db.Collection("svcp")

	// find user by email
	service := CreateNewServices(createServices)
	svcpID := svcp.SVCPID
	filter := bson.D{{Key: "SVCPID", Value: svcpID}}
	res, err := svcp_collection.UpdateOne(context.Background(), filter, bson.D{{Key: "$push", Value: bson.D{{Key: "services", Value: service}}}})
	if res == nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return service, nil
}

func UpdateUserPet(db *models.MongoDB, pet *models.Pet, user_id string, pet_idx int) (string, error) {
	// get collection
	user_collection := db.Collection("user")

	// find user by email
	filter := bson.D{{Key: "_id", Value: user_id}}
	res, err := user_collection.UpdateOne(context.Background(), filter, bson.D{{
		Key: "$set", Value: bson.D{{
			Key: "pets." + fmt.Sprint(pet_idx), Value: pet,
		}},
	}})
	if res.MatchedCount == 0 {
		return "User not found (id=" + user_id + ")", err
	}
	if err != nil {
		return "Failed to update pet", err
	}

	return "", nil

}

func DuplicateService(db *models.MongoDB, serviceId string, svcpID string) (*models.Service, error) {
	// get collection
	svcp_collection := db.Collection("svcp")

	// find service by find svcp --> services
	service := models.Service{}
	filter := bson.D{{Key: "services", Value: bson.D{{Key: "serviceID", Value: serviceId}}}}
	err := svcp_collection.FindOne(context.Background(), filter).Decode(&service)

	if err != nil {
		return nil, err
	}

	service.ServiceID = primitive.NewObjectID().Hex()
	filter = bson.D{{Key: "SVCPID", Value: svcpID}}
	res, err := svcp_collection.UpdateOne(context.Background(), filter, bson.D{{Key: "$push", Value: bson.D{{Key: "services", Value: service}}}})
	if res == nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func GetServiceByID(db *models.MongoDB, serviceID string) (*models.Service, error) {
	// get collection
	svcp_collection := db.Collection("svcp")

	// find service by find svcp --> services
	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$services"}},
		{{Key: "$match", Value: bson.M{"services.serviceID": serviceID}}},
		{{Key: "$project", Value: bson.D{
			{Key: "services", Value: 1},
		}}},
	}
	// Run the aggregation
	cursor, err := svcp_collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	type respondServoce struct {
		Services models.Service `bson:"services"`
	}

	// Decode the documents
	var results []respondServoce
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, errors.New("no document found with the given serviceID " + serviceID)
	}

	return &results[0].Services, nil
}

// deletes a service from a svcp's service list by service ID
func DeleteService(db *models.MongoDB, serviceID string, svcpID string) error {
	// get collection
	svcp_collection := db.Collection("svcp")

	// find svcp by id
	filter := bson.D{{Key: "SVCPID", Value: svcpID}}
	update := bson.D{{Key: "$pull", Value: bson.D{{Key: "services", Value: bson.D{{Key: "serviceID", Value: serviceID}}}}}}
	res, err := svcp_collection.UpdateOne(context.Background(), filter, update)

	if res.ModifiedCount == 0 {
		return errors.New("no document found with the given svcpID and serviceID " + svcpID + " " + serviceID)
	}

	return err
}

func UpdateService(db *models.MongoDB, serviceID string, svcpID string, updateService *bson.M) error {
	// get collection
	svcp_collection := db.Collection("svcp")

	services, err := GetServiceByID(db, serviceID)
	if err != nil {
		return err
	}

	for key, value := range *updateService {
		err = (*services).UpdateField(key, value)
		if err != nil {
			return err
		}
	}

	// find user by email
	filter := bson.D{{Key: "SVCPID", Value: svcpID}, {Key: "services.serviceID", Value: serviceID}}

	// update only send fields
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "services.$", Value: services}}}}
	res, err := svcp_collection.UpdateOne(context.Background(), filter, update)

	if res.MatchedCount == 0 {
		return errors.New("no document found with the given svcpID " + svcpID)
	}
	if err != nil {
		return err
	}

	return nil
}
