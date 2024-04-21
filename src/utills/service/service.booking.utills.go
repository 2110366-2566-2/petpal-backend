package utills

import (
	"context"
	"errors"
	"petpal-backend/src/models"
	payment_utils "petpal-backend/src/utills/payment"
	svcp_utils "petpal-backend/src/utills/serviceprovider"
	user_utils "petpal-backend/src/utills/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertBooking(db *models.MongoDB, BookingCreate *models.Booking, user *models.User) (*models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")
	collectionSVCP := db.Collection("svcp")

	var svcp models.SVCP = models.SVCP{}
	filter := bson.M{
		"services": bson.M{
			"$elemMatch": bson.M{
				"serviceID": BookingCreate.ServiceID,
			},
		},
	}
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
	BookingCreate.SVCPID = svcp.SVCPID

	// Check if the timeslot has already passed
	if foundtimeslot.StartTime.Before(BookingCreate.BookingTimestamp) {
		// return nil, errors.New("timeslot has already passed")
	}
	// Insert the booking into the collection
	_, err = collection.InsertOne(context.Background(), BookingCreate)
	if err != nil {
		return nil, err
	}

	// Send email to notify SVCP
	err = NotifyCreateBooking(BookingCreate, user, &svcp)
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

	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
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

	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
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

	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
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

	// Notification
	svcp, err := svcp_utils.GetSVCPByID(db, booking.SVCPID)
	if err != nil {
		return nil, err
	}

	user, err := user_utils.GetUserByID(db, booking.UserID)
	if err != nil {
		return nil, err
	}

	if booking.Cancel.CancelBy == "user" {
		// Send email to notify SVCP
		err = NotifyCancelBookingToSVCP(&booking, user, svcp)
		if err != nil {
			return nil, err
		}
	} else if booking.Cancel.CancelBy == "svcp" {
		// Send email to notify user
		err = NotifyCancelBookingToUser(&booking, user, svcp)
		if err != nil {
			return nil, err
		}
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

// func GetBookingsBySVCP(db *models.MongoDB, SVCPID string) ([]models.BookingShowALL, error) {
// 	// Get the booking collection
// 	collection := db.Collection("booking")

// 	// Find the booking by SVCPID
// 	filter := bson.D{{Key: "SVCPID", Value: SVCPID}}
// 	cursor, err := collection.Find(context.Background(), filter)
// 	if err != nil {
// 		return nil, err
// 	}

// 	bookings := []models.BookingShowALL{}
// 	if err = cursor.All(context.Background(), &bookings); err != nil {
// 		return nil, err
// 	}

// 	return bookings, nil
// }

func GetAllBookingsBySVCP(db *models.MongoDB, SVCPID string) ([]models.BookingShowALL, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by SVCPID
	filter := bson.D{{Key: "SVCPID", Value: SVCPID}}
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

func ChangeBookingScheduled(db *models.MongoDB, bookingID string, newTimeslotID string, user *models.User) (*models.Booking, error) {

	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by bookingID
	var booking models.Booking = models.Booking{}

	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
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
	booking.Status.SvcpConfirmed = false
	booking.Status.RescheduleStatus = true
	booking.StartTime = foundtimeslot.StartTime
	booking.EndTime = foundtimeslot.EndTime

	// Update the booking in the collection
	_, err = collection.ReplaceOne(context.Background(), filter, booking)
	if err != nil {
		return nil, err
	}

	// Notify SVCP
	err = NotifyRescheduleBooking(&booking, user, &svcp)
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

	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
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

func FillBookingStatusString(db *models.MongoDB, bookingArray []models.BookingShowALL) []models.BookingShowALL {
	timeNow := time.Now()
	const oneHours = 1 * time.Hour
	const twentyFourHours = 72 * time.Hour
	const threeDays = 72 * time.Hour
	for i, b := range bookingArray {
		if !b.Status.PaymentStatus {
			if !b.Status.SvcpCompleted {
				if timeNow.Sub(b.BookingTimestamp) > oneHours {
					payment_utils.UpdateBookingSVCPCompleted(db, b.BookingID)
				}
			}
		}

		if !b.Status.PaymentStatus {
			if b.Cancel.CancelReason == "Payment Expired (Not Authorize Payment within 24 hours)" {
				bookingArray[i].StatusString = "Payment Expired"
			} else if timeNow.Sub(b.BookingTimestamp) > twentyFourHours {
				payment_utils.CheckUpdateExpiredBookingPayment(db, b.BookingID)
				bookingArray[i].StatusString = "Payment Expired"
			} else {
				bookingArray[i].StatusString = "Pending Payment"
			}
		} else if !b.Status.SvcpCompleted {
			if b.Cancel.CancelStatus {
				bookingArray[i].StatusString = "Cancelled"
			} else {
				bookingArray[i].StatusString = "Paid"
			}
		} else if !b.Status.UserCompleted {
			bookingArray[i].StatusString = "Completed"
		} else if timeNow.Sub(b.Status.SvcpCompletedTimestamp) > threeDays {
			bookingArray[i].StatusString = "Completed"
			CompleteBooking(db, b.BookingID, "user")
		} else {
			if b.Status.UserRefund {
				bookingArray[i].StatusString = "Refunded"
			} else {
				bookingArray[i].StatusString = "Service Completed"
			}
		}
	}

	return bookingArray
}

func CompleteBooking(db *models.MongoDB, bookingID string, userType string) (*models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by bookingID
	var booking models.Booking = models.Booking{}

	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
	if err != nil {
		return nil, err
	}

	// Update the booking status
	if userType == "svcp" {
		booking.Status.SvcpCompleted = true
		booking.Status.SvcpCompletedTimestamp = time.Now()
	} else if userType == "user" {
		moneyToSVCP := payment_utils.CalculateFee(booking.TotalBookingPrice)
		err = payment_utils.SendMoneyToSVCP(db, booking.SVCPID, moneyToSVCP)
		if err != nil {
			return nil, err
		}
		booking.Status.UserCompleted = true
		booking.Status.UserCompletedTimestamp = time.Now()
	}
	// Update the booking in the collection
	_, err = collection.ReplaceOne(context.Background(), filter, booking)
	if err != nil {
		return nil, err
	}

	// Notification
	svcp, err := svcp_utils.GetSVCPByID(db, booking.SVCPID)
	if err != nil {
		return nil, err
	}

	user, err := user_utils.GetUserByID(db, booking.UserID)
	if err != nil {
		return nil, err
	}

	if userType == "svcp" {
		// Send email to notify user
		err = NotifyCompleteBookingToUser(&booking, user, svcp)
		if err != nil {
			return nil, err
		}
	} else if userType == "user" {
		// Send email to notify SVCP
		err = NotifyCompleteBookingToSVCP(&booking, user, svcp)
		if err != nil {
			return nil, err
		}
	}

	return &booking, nil
}

func SVCPConfirmBooking(db *models.MongoDB, bookingID string, svcp *models.SVCP) (*models.Booking, error) {
	// Get the booking collection
	collection := db.Collection("booking")

	// Find the booking by bookingID
	var booking models.Booking = models.Booking{}

	filter := bson.D{{Key: "_id", Value: bookingID}}
	err := collection.FindOne(context.Background(), filter).Decode(&booking)
	if err != nil {
		return nil, err
	}

	// Update the booking status
	booking.Status.SvcpConfirmed = true
	booking.Status.SvcpConfirmedTimestamp = time.Now()
	// Update the booking in the collection
	_, err = collection.ReplaceOne(context.Background(), filter, booking)
	if err != nil {
		return nil, err
	}

	// Notification
	user, err := user_utils.GetUserByID(db, booking.UserID)
	if err != nil {
		return nil, err
	}

	// Send email to notify user
	err = NotifyConfirmBookingToUser(&booking, user, svcp)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}
