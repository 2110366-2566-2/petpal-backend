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

	// BookingCreate.SVCPName = svcp.SVCPUsername
	// BookingCreate.AverageRating = foundService.AverageRating
	// BookingCreate.ServiceImg = foundService.ServiceImg
	// BookingCreate.ServiceDescription = foundService.ServiceDescription
	// BookingCreate.SVCPName = svcp.SVCPUsername
	BookingCreate.ServiceName = foundService.ServiceName
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

func GetABookingDetail(db *models.MongoDB, bookingID string) (*models.BookingFull, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by bookingID
	booking := models.BookingFull{}

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

	booking.SVCPName = svcp.SVCPUsername
	booking.ServiceName = foundService.ServiceName
	booking.AverageRating = foundService.AverageRating
	booking.ServiceImg = foundService.ServiceImg
	booking.ServiceDescription = foundService.ServiceDescription

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
func CancelBooking(db *models.MongoDB, bookingID string, Cancel models.BookingCancel) (*models.Booking, error) {

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
	booking.Cancel = Cancel

	// Update the booking in the collection
	_, err = collection.ReplaceOne(context.Background(), filter, booking)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

// func GetBookingsAByUser(db *models.MongoDB, userID string) ([]models.BookingWithId, error) {
// 	// Get the booking collection
// 	collection := db.Collection("booking")

// 	// Find the booking by userID
// 	filter := bson.D{{Key: "userID", Value: userID}}
// 	cursor, err := collection.Find(context.Background(), filter)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var bookings []models.BookingWithId
// 	if err = cursor.All(context.Background(), &bookings); err != nil {
// 		return nil, err
// 	}

// 	return bookings, nil
// }

func GetAllBookingsByUser(db *models.MongoDB, userID string) ([]models.BookingShowALL, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by userID
	filter := bson.D{{Key: "userID", Value: userID}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var bookings []models.BookingShowALL
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

// func BookingArrayGetTimeSlot(db *models.MongoDB, bookingArray []models.BookingWithId) []models.Timeslot {
// 	var timeslotArray []models.Timeslot
// 	for _, b := range bookingArray {
// 		timeslot, err := BookingGetTimeSlot(db, b.BookingID)
// 		if err != nil {
// 			timeslotArray = append(timeslotArray, models.Timeslot{})
// 			continue
// 		}
// 		timeslotArray = append(timeslotArray, *timeslot)
// 	}
// 	return timeslotArray
// }

func AllBookFilter(db *models.MongoDB, bookingArray []models.BookingShowALL, Bookfilter models.RequestBookingAll) []models.BookingShowALL {

	var filteredBooking []models.BookingShowALL
	//TimeslotdArray := BookingArrayGetTimeSlot(db, bookingArray)

	for _, b := range bookingArray {
		if !Bookfilter.StartAfter.IsZero() {
			if b.StartTime.Before(Bookfilter.StartAfter) {
				continue
			}
		}
		if Bookfilter.ReservationType != "" {
			if Bookfilter.ReservationType == "incoming" {
				if b.StartTime.Before(time.Now()) {
					continue
				}
			} else if Bookfilter.ReservationType == "outgoing" {
				if b.StartTime.After(time.Now()) {
					continue
				}
			}
		}

		if Bookfilter.CancelStatus != 2 && (b.Cancel.CancelStatus != (Bookfilter.CancelStatus == 1)) {
			continue
		}
		if Bookfilter.PaymentStatus != 2 && (b.Status.PaymentStatus != (Bookfilter.PaymentStatus == 1)) {
			continue
		}
		if Bookfilter.SvcpConfirmed != 2 && (b.Status.SvcpConfirmed != (Bookfilter.SvcpConfirmed == 1)) {
			continue
		}
		if Bookfilter.SvcpCompleted != 2 && (b.Status.SvcpCompleted != (Bookfilter.SvcpCompleted == 1)) {
			continue
		}
		if Bookfilter.UserCompleted != 2 && (b.Status.UserCompleted != (Bookfilter.UserCompleted == 1)) {
			continue
		}

		filteredBooking = append(filteredBooking, b)

	}

	return filteredBooking
}

func FillSVCPDetail(db *models.MongoDB, bookingArray []models.BookingShowALL) []models.BookingShowALL {
	collectionSVCP := db.Collection("svcp")
	for i, b := range bookingArray {
		var svcp models.SVCP = models.SVCP{}
		filterSvcp := bson.D{{Key: "SVCPID", Value: b.SVCPID}}
		err := collectionSVCP.FindOne(context.Background(), filterSvcp).Decode(&svcp)
		if err != nil {
			bookingArray[i].SVCPName = "svcp not found"
		} else {
			bookingArray[i].SVCPName = svcp.SVCPUsername
		}

	}

	return bookingArray
}
