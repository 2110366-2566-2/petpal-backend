package utills

import (
	"context"
	"fmt"
	"petpal-backend/src/models"
	email_utils "petpal-backend/src/utills/email"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func ConfirmBookingPayment(db *models.MongoDB, bookingID string, userID string) (*models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")
	// Find the booking by bookingID
	var booking models.Booking = models.Booking{}
	// Convert bookingID to ObjectID
	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
	if err != nil {
		return nil, err
	}
	if bookingID == "" {
		return nil, fmt.Errorf("Missing BookingID")
	}
	if booking.UserID != userID {
		return nil, fmt.Errorf("This booking is not belong to you")
	}
	if booking.Cancel.CancelStatus {
		return nil, fmt.Errorf("This booking is already cancelled")
	}
	timeNow := time.Now()
	timeDiff := timeNow.Sub(booking.BookingTimestamp)
	const twentyFourHours = 24 * time.Hour
	if timeDiff > twentyFourHours {
		UpdateExpiredBookingPayment(booking)
	} else {
		err := SendMoneyToBank(db, "Petpal", booking.TotalBookingPrice)
		if err != nil {
			return nil, err
		}
		booking.Status.PaymentStatus = true
		booking.Status.PaymentTimestamp = timeNow
	}
	_, err = collection.UpdateOne(context.Background(), filter, bson.D{{Key: "$set", Value: booking}})
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func UpdateExpiredBookingPayment(booking models.Booking) models.Booking {
	timeNow := time.Now()
	timeDiff := timeNow.Sub(booking.BookingTimestamp)
	const twentyFourHours = 24 * time.Hour
	if (timeDiff > twentyFourHours) && (!booking.Status.PaymentStatus) {
		booking.Cancel.CancelStatus = true
		booking.Cancel.CancelTimestamp = timeNow
		booking.Cancel.CancelReason = "Payment Expired (Not Authorize Payment within 24 hours)"
		booking.Cancel.CancelBy = "Petpal Admin"
	}
	return booking
}

func CheckUpdateExpiredBookingPayment(db *models.MongoDB, bookingID string) error {
	// Get the booking collection
	collection := db.Collection("booking")
	// Find the booking by bookingID
	var booking models.Booking = models.Booking{}
	// Convert bookingID to ObjectID
	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
	if err != nil {
		return err
	}
	if booking.Cancel.CancelStatus {
		return fmt.Errorf("This booking is already cancelled")
	}
	booking = UpdateExpiredBookingPayment(booking)

	_, err = collection.UpdateOne(context.Background(), filter, bson.D{{Key: "$set", Value: booking}})
	if err != nil {
		return err
	}

	return nil
}

func UpdateBookingSVCPCompleted(db *models.MongoDB, bookingID string) error {
	// Get the booking collection
	collection := db.Collection("booking")
	// Find the booking by bookingID
	var booking models.Booking = models.Booking{}
	// Convert bookingID to ObjectID
	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
	if err != nil {
		return err
	}
	if booking.Cancel.CancelStatus {
		return fmt.Errorf("This booking is already cancelled")
	}
	booking.Status.SvcpCompleted = true
	booking.Status.SvcpCompletedTimestamp = time.Now()
	booking.Status.SvcpConfirmed = true
	booking.Status.SvcpConfirmedTimestamp = time.Now()

	_, err = collection.UpdateOne(context.Background(), filter, bson.D{{Key: "$set", Value: booking}})
	if err != nil {
		return err
	}
	return nil
}

func RefundBooking(db *models.MongoDB, bookingID string) error {
	// Get the booking collection
	collection := db.Collection("booking")
	// Find the booking by bookingID
	var booking models.Booking = models.Booking{}
	// Convert bookingID to ObjectID
	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
	if err != nil {
		return err
	}
	if booking.Cancel.CancelStatus {
		return fmt.Errorf("This booking is already cancelled")
	}
	err = SendMoneyToUser(db, booking.UserID, booking.TotalBookingPrice)
	if err != nil {
		return err
	}
	booking.Cancel.CancelStatus = true
	booking.Cancel.CancelTimestamp = time.Now()
	booking.Cancel.CancelReason = "Refund"
	booking.Cancel.CancelBy = "User"
	booking.Status.UserRefund = true
	booking.Status.UserRefundTimestamp = time.Now()

	_, err = collection.UpdateOne(context.Background(), filter, bson.D{{Key: "$set", Value: booking}})
	if err != nil {
		return err
	}
	return nil
}

func CalculateFee(money float64) float64 {
	return money * 0.97
}

func SendMoneyToSVCP(db *models.MongoDB, SVCPID string, money float64) error {
	// Get the booking collection
	collection := db.Collection("svcp")
	// Find the booking by bookingID
	var SVCP models.SVCP
	// Convert bookingID to ObjectID
	filter := bson.D{{Key: "_id", Value: SVCPID}}
	err := collection.FindOne(context.Background(), filter).Decode(&SVCP)
	if err != nil {
		return err
	}

	err = SendMoneyToBank(db, SVCP.DefaultBank, money)
	if err != nil {
		return err
	}
	return nil
}

func SendMoneyToUser(db *models.MongoDB, userID string, money float64) error {
	// Get the booking collection
	collection := db.Collection("user")
	// Find the booking by bookingID
	var user models.User
	// Convert bookingID to ObjectID
	filter := bson.D{{Key: "_id", Value: userID}}
	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return err
	}

	err = SendMoneyToBank(db, user.DefaultBank, money)
	if err != nil {
		return err
	}
	return nil
}

func SendMoneyToBank(db *models.MongoDB, bankID string, money float64) error {
	// monet to string
	moneyStr := fmt.Sprintf("%f", money)

	email_utils.SendEmailWithGmail("petpalgoldenarmorwarrior@gmail.com", "PetPal Confirm payment", "คุณได้จ่ายเงิน "+moneyStr+" บาท ไปยังธนาคาร "+bankID+" สำเร็จแล้ว")
	return nil
}
