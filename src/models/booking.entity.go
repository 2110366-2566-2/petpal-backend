// user.go
package models

import "time"

// User represents a user entity
type Booking struct {
	UserID     string `json:"userID" bson:"userID"`
	SVCPID     string `json:"SVCPID" bson:"SVCPID"`
	ServiceID  string `json:"serviceID" bson:"serviceID"`
	TimeslotID string `json:"timeslotID" bson:"timeslotID"`

	BookingTimestamp  time.Time `json:"bookingTimestamp" bson:"bookingTimestamp"`
	TotalBookingPrice float64   `json:"totalBookingPrice" bson:"totalBookingPrice"`

	ServiceName string    `json:"serviceName" bson:"serviceName"`
	StartTime   time.Time `json:"startTime" bson:"startTime"`
	EndTime     time.Time `json:"endTime" bson:"endTime"`

	Cancel BookingCancel `json:"cancel" bson:"cancel"`

	Status BookingStatus `json:"status" bson:"status"`

	Feedback Feedback `json:"feedback" bson:"feedback"`
}
type BookingCancel struct {
	CancelStatus    bool      `json:"cancelStatus" bson:"cancelStatus"`       //true if cancelled
	CancelTimestamp time.Time `json:"cancelTimestamp" bson:"cancelTimestamp"` //time of cancellation
	CancelReason    string    `json:"cancelReason" bson:"cancelReason"`       //reason for cancellation
	CancelBy        string    `json:"cancelBy" bson:"cancelBy"`               //who cancelled
}
type BookingStatus struct {
	RescheduleStatus bool `json:"rescheduleStatus" bson:"rescheduleStatus"` //true if rescheduled

	PaymentStatus          bool      `json:"paymentStatus" bson:"paymentStatus"`
	PaymentTimestamp       time.Time `json:"paymentTimestamp" bson:"paymentTimestamp"`
	SvcpConfirmed          bool      `json:"svcpConfirmed" bson:"svcpConfirmed"`
	SvcpConfirmedTimestamp time.Time `json:"svcpConfirmedTimestamp" bson:"svcpConfirmedTimestamp"`
	SvcpCompleted          bool      `json:"svcpCompleted" bson:"svcpCompleted"`
	SvcpCompletedTimestamp time.Time `json:"svcpCompletedTimestamp" bson:"svcpCompletedTimestamp"`
	UserCompleted          bool      `json:"userCompleted" bson:"userCompleted"`
	UserCompletedTimestamp time.Time `json:"userCompletedTimestamp" bson:"userCompletedTimestamp"`
}

// type BookingWithId struct {
// 	BookingID  string `json:"bookingID" bson:"_id"`
// 	UserID     string `json:"userID" bson:"userID"`
// 	SVCPID     string `json:"SVCPID" bson:"SVCPID"`
// 	ServiceID  string `json:"serviceID" bson:"serviceID"`
// 	TimeslotID string `json:"timeslotID" bson:"timeslotID"`
// 	// BookingStatus     BookingStatus `json:"bookingStatus" bson:"bookingStatus"`
// 	BookingTimestamp  time.Time `json:"bookingTimestamp" bson:"bookingTimestamp"`
// 	TotalBookingPrice float64   `json:"totalBookingPrice" bson:"totalBookingPrice"`
// 	Feedback          Feedback  `json:"feedback" bson:"feedback"`
// }

type BookingShowALL struct {
	BookingID         string    `json:"bookingID" bson:"_id"`
	UserID            string    `json:"userID" bson:"userID"`
	SVCPID            string    `json:"SVCPID" bson:"SVCPID"`
	ServiceID         string    `json:"serviceID" bson:"serviceID"`
	TimeslotID        string    `json:"timeslotID" bson:"timeslotID"`
	BookingTimestamp  time.Time `json:"bookingTimestamp" bson:"bookingTimestamp"`
	TotalBookingPrice float64   `json:"totalBookingPrice" bson:"totalBookingPrice"`

	ServiceName string    `json:"serviceName" bson:"serviceName"`
	SVCPName    string    `json:"SVCPName" bson:"SVCPName"`
	StartTime   time.Time `json:"startTime" bson:"startTime"`
	EndTime     time.Time `json:"endTime" bson:"endTime"`

	Cancel BookingCancel `json:"cancel" bson:"cancel"`

	Status BookingStatus `json:"status" bson:"status"`

	Feedback Feedback `json:"feedback" bson:"feedback"`

	// BookingStatus BookingStatus `json:"bookingStatus" bson:"bookingStatus"`
}

type BookingFull struct {
	BookingID  string `json:"bookingID" bson:"_id"`
	UserID     string `json:"userID" bson:"userID"`
	SVCPID     string `json:"SVCPID" bson:"SVCPID"`
	ServiceID  string `json:"serviceID" bson:"serviceID"`
	TimeslotID string `json:"timeslotID" bson:"timeslotID"`

	BookingTimestamp  time.Time `json:"bookingTimestamp" bson:"bookingTimestamp"`
	TotalBookingPrice float64   `json:"totalBookingPrice" bson:"totalBookingPrice"`

	SVCPName           string  `json:"SVCPName" bson:"SVCPName"`
	ServiceName        string  `json:"serviceName" bson:"serviceName"` //service name
	AverageRating      float64 `json:"averageRating" bson:"averageRating"`
	ServiceImg         []byte  `json:"serviceImg" bson:"serviceImg"`
	ServiceDescription string  `json:"serviceDescription" bson:"serviceDescription"`

	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime"`

	// BookingStatus BookingStatus `json:"bookingStatus" bson:"bookingStatus"`

	Cancel BookingCancel `json:"cancel" bson:"cancel"`

	Status BookingStatus `json:"status" bson:"status"`

	Feedback Feedback `json:"feedback" bson:"feedback"`
}

type BookingFullNoID struct {
	UserID     string `json:"userID" bson:"userID"`
	SVCPID     string `json:"SVCPID" bson:"SVCPID"`
	ServiceID  string `json:"serviceID" bson:"serviceID"`
	TimeslotID string `json:"timeslotID" bson:"timeslotID"`

	BookingTimestamp  time.Time `json:"bookingTimestamp" bson:"bookingTimestamp"`
	TotalBookingPrice float64   `json:"totalBookingPrice" bson:"totalBookingPrice"`

	SVCPName           string  `json:"SVCPName" bson:"SVCPName"`
	AverageRating      float64 `json:"averageRating" bson:"averageRating"`
	ServiceImg         []byte  `json:"serviceImg" bson:"serviceImg"`
	ServiceDescription string  `json:"serviceDescription" bson:"serviceDescription"`

	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime"`

	Cancel BookingCancel `json:"cancel" bson:"cancel"`

	Status BookingStatus `json:"status" bson:"status"`

	Feedback Feedback `json:"feedback" bson:"feedback"`
}

type BookingRequest struct {
	ServiceID  string `json:"serviceID" bson:"serviceID"`
	TimeslotID string `json:"timeslotID" bson:"timeslotID"`
}

type BookingIndex struct {
	UserID     string `json:"userID" bson:"userID"`
	SVCPID     string `json:"SVCPID" bson:"SVCPID"`
	ServiceID  string `json:"serviceID" bson:"serviceID"`
	TimeslotID string `json:"timeslotID" bson:"timeslotID"`
}

// type BookingStatus string

// const (
// 	BookingPending   BookingStatus = "pending payment"            //waiting for user to pay 					1
// 	BookingPaid      BookingStatus = "payment paid"               //user has paid wait svcp comfirm				2 or 3
// 	BookingComfirmed BookingStatus = "service provided comfirmed" //svcp has confirmed waiting for user to pay 	3 or 2
// 	BookingCompleted BookingStatus = "completed"                  //service has been provided					4

// 	BookingCanceledUser BookingStatus = "cancelled by user"             //user has cancelled
// 	BookingCanceledSvcp BookingStatus = "cancelled by service provider" //svcp has cancelled

// 	//BookingRescheduled BookingStatus = "rescheduled" //user has rescheduled

// 	// may be un used
// 	BookingExpiredPaid      BookingStatus = "expired from unpaid"                                //user has not paid in time
// 	BookingExpiredComfirmed BookingStatus = "expired from pending service provider confirmation" //svcp has not confirmed in time
// )

// var (
// 	BookingStatusDone    = []BookingStatus{BookingCompleted, BookingCanceledUser, BookingCanceledSvcp, BookingExpiredPaid, BookingExpiredComfirmed}
// 	BookingStatusNotdone = []BookingStatus{BookingPending, BookingPaid, BookingComfirmed}
// )

type RequestCancelBooking struct {
	BookingID    string `json:"bookingID" bson:"_id"`
	CancelReason string `json:"cancelReason" bson:"cancelReason"`
}
type RequestBookingId struct {
	BookingID string `json:"bookingID"`
}

type RequestBookingRescheduled struct {
	BookingID  string `json:"bookingID"`
	TimeslotID string `json:"timeslotID" bson:"timeslotID"`
}

type RequestBookingAll struct {
	StartAfter      time.Time `json:"startAfter" `
	ReservationType string    `json:"reservationType"`
	CancelStatus    int       `json:"cancelStatus" bson:"cancelStatus"`
	PaymentStatus   int       `json:"paymentStatus" bson:"paymentStatus"`
	SvcpConfirmed   int       `json:"svcpConfirmed" bson:"svcpConfirmed"`
	SvcpCompleted   int       `json:"svcpCompleted" bson:"svcpCompleted"`
	UserCompleted   int       `json:"userCompleted" bson:"userCompleted"`
}
type BookkingDetailRes struct {
	Message string      `json:"message"`
	Result  BookingFull `json:"result"`
}

type BookingBasicRes struct {
	Message string  `json:"message"`
	Result  Booking `json:"result"`
}

// type BookingWithIdRes struct {
// 	Message string        `json:"message"`
// 	Result  BookingWithId `json:"result"`
// }

type BookingWithIdArrayRes struct {
	Message string           `json:"message"`
	Result  []BookingShowALL `json:"result"`
}

type PromptpayQr struct {
	QrImage []byte `json:"qrImage"`
}
