package utills

import (
	"context"
	"petpal-backend/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConfirmBookingPayment(db *models.MongoDB, bookingID string) (*models.Booking, error) {
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
	timeNow := time.Now()
	timeDiff := timeNow.Sub(booking.BookingTimestamp)
	const twentyFourHours = 24 * time.Hour
	if timeDiff > twentyFourHours {
		booking.Cancel.CancelStatus = true
		booking.Cancel.CancelTimestamp = timeNow
		booking.Cancel.CancelReason = "Payment Expired (Not Authorize Payment within 24 hours)"
		booking.Cancel.CancelBy = "Petpal Admin"
	} else {
		booking.Status.PaymentStatus = true
		booking.Status.PaymentTimestamp = timeNow
	}
	// Update the booking in the collection
	_, err = collection.ReplaceOne(context.Background(), filter, booking)
	if err != nil {
		return nil, err
	}
	return &booking, nil
}
