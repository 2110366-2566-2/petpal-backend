// user.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user entity
type CreateTimeslot struct {
	StartTime time.Time `json:"startTime" bson:"startTime"`
	EndTime   time.Time `json:"endTime" bson:"endTime"`
}

type Timeslot struct {
	TimeslotID string    `json:"timeslotID" bson:"timeslotID,omitempty"`
	StartTime  time.Time `json:"startTime" bson:"startTime"`
	EndTime    time.Time `json:"endTime" bson:"endTime"`
	Status     string    `json:"status" bson:"status"`
}

func (u *Timeslot) SetDefult() {
	u.Status = "available"
	u.StartTime = time.Now()
	u.EndTime = time.Now().Add(1 * time.Hour)
	u.TimeslotID = primitive.NewObjectID().Hex()
}

func (u *Timeslot) UpdateField(key string, value interface{}) error {
	// UpdateField
	// get the field and update it
	// return the updated service
	var err error
	if key == "startTime" {
		e := value.(string)
		u.EndTime, err = time.Parse(time.RFC3339, e)
		if err != nil {
			return err
		}
	}
	if key == "endTime" {
		e := value.(string)

		u.EndTime, err = time.Parse(time.RFC3339, e)
		if err != nil {
			return err
		}
	}
	if key == "status" {
		u.Status = value.(string)
	}
	return nil
}
