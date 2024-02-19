// user.go
package models

import "time"

// User represents a user entity
type Timeslot struct {
	TimeslotID  string    `json:"timeslotID" bson:"timeslotID"`
	StartTime   time.Time `json:"startTime" bson:"startTime"`
	EndTime     time.Time `json:"endTime" bson:"endTime"`
}

func (u *Timeslot) editTimeslot(timeslotDetails Timeslot) Timeslot {
	// Mock Function
	return *u
}
