package utills

import (
	"context"
	"errors"
	"petpal-backend/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func InsertBooking(db *models.MongoDB, BookingCreate *models.BookingFullNoID) (*models.BookingFullNoID, error) {
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
		//println(s.ServiceID, BookingCreate.ServiceID)
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

	BookingCreate.SVCPName = svcp.SVCPUsername
	BookingCreate.AverageRating = foundService.AverageRating
	BookingCreate.ServiceImg = foundService.ServiceImg
	BookingCreate.ServiceDescription = foundService.ServiceDescription
	BookingCreate.StartTime = foundtimeslot.StartTime
	BookingCreate.EndTime = foundtimeslot.EndTime

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

func GetBooking(db *models.MongoDB, bookingID string) (*models.BookingWithId, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by bookingID
	var booking models.BookingWithId = models.BookingWithId{}

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

func GetBookingsByUser(db *models.MongoDB, userID string) ([]models.BookingWithId, error) {
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

func GetBookingsBySVCP(db *models.MongoDB, SVCPID string) ([]models.Booking, error) {
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

func ChangeBookingScheduled(db *models.MongoDB, bookingID string, newTimeslotID string) (*models.Booking, error) {

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

	if booking.TimeslotID == newTimeslotID {
		return nil, errors.New("new timeslot is the same as the old timeslot")
	}

	collectionSVCP := db.Collection("svcp")
	var svcp models.SVCP = models.SVCP{}
	filterSvcp := bson.D{{Key: "SVCPID", Value: booking.SVCPID}}
	err = collectionSVCP.FindOne(context.Background(), filterSvcp).Decode(&svcp)
	if err != nil {
		return nil, err
	}

	// Check if the service exists in the service provider
	var foundService models.Service
	err = errors.New("service not found")
	for _, s := range svcp.Services {
		//println(s.ServiceID, booking.ServiceID)
		if s.ServiceID == booking.ServiceID {
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
		if t.TimeslotID == newTimeslotID {
			foundtimeslot = t
			err = nil
			break
		}
	}

	if err != nil {
		return nil, err
	}

	// Check if the timeslot has already passed
	if foundtimeslot.StartTime.Before(booking.BookingTimestamp) {
		// return nil, errors.New("timeslot has already passed")
	}

	// Update the booking status
	booking.TimeslotID = newTimeslotID

	// Update the booking in the collection
	_, err = collection.ReplaceOne(context.Background(), filter, booking)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func BookingRefund() {
	//todo
	println("booking refund payment to user")

}

func BookingGetTimeSlot(db *models.MongoDB, bookingID string) (*models.Timeslot, error) {
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

	collectionSVCP := db.Collection("svcp")
	var svcp models.SVCP = models.SVCP{}
	filterSvcp := bson.D{{Key: "SVCPID", Value: booking.SVCPID}}
	err = collectionSVCP.FindOne(context.Background(), filterSvcp).Decode(&svcp)
	if err != nil {
		return nil, err
	}

	// Check if the service exists in the service provider
	var foundService models.Service
	err = errors.New("service not found")
	for _, s := range svcp.Services {
		//println(s.ServiceID, booking.ServiceID)
		if s.ServiceID == booking.ServiceID {
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
		if t.TimeslotID == booking.TimeslotID {
			foundtimeslot = t
			err = nil
			break
		}
	}

	if err != nil {
		return nil, err
	}

	return &foundtimeslot, nil
}
func BookingArrayGetTimeSlot(db *models.MongoDB, bookingArray []models.BookingWithId) []models.Timeslot {
	var timeslotArray []models.Timeslot
	for _, b := range bookingArray {
		timeslot, err := BookingGetTimeSlot(db, b.BookingID)
		if err != nil {
			timeslotArray = append(timeslotArray, models.Timeslot{})
			continue
		}
		timeslotArray = append(timeslotArray, *timeslot)
	}
	return timeslotArray
}

func AllBookFilter(db *models.MongoDB, bookingArray []models.BookingWithId, Bookfilter models.RequestBookingAll) []models.BookingWithId {

	var filteredBooking []models.BookingWithId
	TimeslotdArray := BookingArrayGetTimeSlot(db, bookingArray)

	for i, b := range bookingArray {
		if !Bookfilter.TimeslotStartAfter.IsZero() {
			if TimeslotdArray[i].StartTime.Before(Bookfilter.TimeslotStartAfter) {
				continue
			}
		}

		if Bookfilter.ReservationType != "" {
			if Bookfilter.ReservationType == "incoming" {
				if !TimeslotdArray[i].StartTime.Before(time.Now()) {
					continue
				}
			} else if Bookfilter.ReservationType == "outgoing" {
				if !TimeslotdArray[i].StartTime.After(time.Now()) {
					continue
				}
			}
		}

		filteredBooking = append(filteredBooking, b)

	}

	return filteredBooking
}
