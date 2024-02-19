// user.go
package models

import "time"

// User represents a user entity
type Booking struct {
	UserID            string    `json:"userID" bson:"userID"`
	SVCPID            string    `json:"SVCPID" bson:"SVCPID"`
	ServiceID         string    `json:"serviceID" bson:"serviceID"`
	TimeslotID        string    `json:"timeslotID" bson:"timeslotID"`
	BookingStatus     string    `json:"bookingStatus" bson:"bookingStatus"`
	BookingTimestamp  time.Time `json:"bookingTimestamp" bson:"bookingTimestamp"`
	TotalBookingPrice float64   `json:"totalBookingPrice" bson:"totalBookingPrice"`
	Feedback          Feedback  `json:"feedback" bson:"feedback"`
}
