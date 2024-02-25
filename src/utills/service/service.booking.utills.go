package utills

import (
	"context"
	"errors"
	"petpal-backend/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func InsertBooking(db *models.MongoDB, BookingCreate *models.Booking) (*models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")
	collectionSVCP := db.Collection("svcp")

	var svcp models.SVCP = models.SVCP{}
	filter := bson.D{{Key: "SVCPID", Value: BookingCreate.SVCPID}}
	err := collectionSVCP.FindOne(context.Background(), filter).Decode(&svcp)
	if err != nil {
		return nil, err
	}

	// Check if the service exists in the service provider
	var foundService models.Service
	err = errors.New("service not found")
	for _, s := range svcp.Services {
		println(s.ServiceID, BookingCreate.ServiceID)
		if s.ServiceID == BookingCreate.ServiceID {
			foundService = s
			err = nil
			break
		}
	}
	if err != nil {
		return nil, err
	}

	// Check if the timeslot exists in the service
	var foundtimeslot models.Timeslot
	err = errors.New("timeslot not found")
	for _, t := range foundService.Timeslots {
		if t.TimeslotID == BookingCreate.TimeslotID {
			foundtimeslot = t
			err = nil
			break
		}
	}

	if err != nil {
		return nil, err
	}

	BookingCreate.TotalBookingPrice = foundService.Price

	// Check if the timeslot has already passed
	if foundtimeslot.StartTime.Before(BookingCreate.BookingTimestamp) {
		// return nil, errors.New("timeslot has already passed")
	}

	// Insert the booking into the collection
	_, err = collection.InsertOne(context.Background(), BookingCreate)
	if err != nil {
		return nil, err
	}

	// Return the inserted booking
	return BookingCreate, nil
}

func ChangeBookingStatus(db *models.MongoDB, bookingID string, status models.BookingStatus) (*models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by bookingID
	var booking models.Booking = models.Booking{}

	// Convert bookingID to ObjectID
	objID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	err = collection.FindOne(context.Background(), filter).Decode(&booking)
	if err != nil {
		return nil, err
	}

	// Update the booking status
	booking.BookingStatus = status

	// Update the booking in the collection
	_, err = collection.ReplaceOne(context.Background(), filter, booking)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func GetBooking(db *models.MongoDB, bookingID string) (*models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by bookingID
	var booking models.Booking = models.Booking{}

	objID, err := primitive.ObjectIDFromHex(bookingID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	err = collection.FindOne(context.Background(), filter).Decode(&booking)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func GetBookingsByUser(db *models.MongoDB, userID string, status models.BookingStatus) ([]models.BookingWithId, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by userID
	filter := bson.D{{Key: "userID", Value: userID}, {Key: "bookingStatus", Value: status}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var bookings []models.BookingWithId
	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func GetAllBookingsByUser(db *models.MongoDB, userID string) ([]models.BookingWithId, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by userID
	filter := bson.D{{Key: "userID", Value: userID}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var bookings []models.BookingWithId
	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func GetBookingsBySVCP(db *models.MongoDB, SVCPID string, status models.BookingStatus) ([]models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by SVCPID
	filter := bson.D{{Key: "SVCPID", Value: SVCPID}, {Key: "bookingStatus", Value: status}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var bookings []models.Booking
	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func GetAllBookingsBySVCP(db *models.MongoDB, SVCPID string) ([]models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by SVCPID
	filter := bson.D{{Key: "SVCPID", Value: SVCPID}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var bookings []models.Booking
	if err = cursor.All(context.Background(), &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}
