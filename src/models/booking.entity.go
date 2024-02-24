// user.go
package models

import "time"

// User represents a user entity
type Booking struct {
	UserID            string        `json:"userID" bson:"userID"`
	SVCPID            string        `json:"SVCPID" bson:"SVCPID"`
	ServiceID         string        `json:"serviceID" bson:"serviceID"`
	TimeslotID        string        `json:"timeslotID" bson:"timeslotID"`
	BookingStatus     BookingStatus `json:"bookingStatus" bson:"bookingStatus"`
	BookingTimestamp  time.Time     `json:"bookingTimestamp" bson:"bookingTimestamp"`
	TotalBookingPrice float64       `json:"totalBookingPrice" bson:"totalBookingPrice"`
	Feedback          Feedback      `json:"feedback" bson:"feedback"`
}

type BookingCreate struct {
	SVCPID     string `json:"SVCPID" bson:"SVCPID"`
	ServiceID  string `json:"serviceID" bson:"serviceID"`
	TimeslotID string `json:"timeslotID" bson:"timeslotID"`
}

type BookingStatus string

const (
	BookingPending   BookingStatus = "pending payment"            //waiting for user to pay
	BookingComfirmed BookingStatus = "service provided comfirmed" //svcp has confirmed waiting for user to pay
	BookingPaid      BookingStatus = "payment confirmed"          //user has paid
	BookingCompleted BookingStatus = "completed"                  //service has been provided

	BookingCanceledUser BookingStatus = "cancelled by user"             //user has cancelled
	BookingCanceledSvcp BookingStatus = "cancelled by service provider" //svcp has cancelled

	BookingRescheduled BookingStatus = "rescheduled" //user has rescheduled

	BookingExpiredPaid      BookingStatus = "expired from unpaid"                                //user has not paid in time
	BookingExpiredComfirmed BookingStatus = "expired from pending service provider confirmation" //svcp has not confirmed in time
)
