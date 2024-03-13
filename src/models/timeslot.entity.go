// user.go
package models

import "time"

// User represents a user entity
type CreateTimeslot struct {
	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime"`
}

type Timeslot struct {
	TimeslotID string    `json:"timeslotID" bson:"timeslotID"`
	StartTime  time.Time `json:"startTime" bson:"startTime"`
	EndTime    time.Time `json:"endTime" bson:"endTime"`
	Status     string    `json:"status" bson:"status"`
}

func (u *Timeslot) editTimeslot(timeslotDetails Timeslot) Timeslot {
	// Mock Function
	return *u
}
