// user.go
package models

// User represents a user entity
type Service struct {
	ServiceType        string     `json:"serviceType" bson:"serviceType"`
	ServiceDescription string     `json:"serviceDescription" bson:"serviceDescription"`
	ServiceImg         string     `json:"serviceImg" bson:"serviceImg"`
	AverageRating      float32    `json:"averageRating" bson:"averageRating"`
	RequireCert        bool       `json:"requireCert" bson:"requireCert"`
	Timeslots          []Timeslot `json:"timeslots" bson:"timeslots"`
}

func (u *Service) createTimeslot(timeslotDetails Timeslot) Timeslot {
	timeslot := Timeslot{
		StartTime: timeslotDetails.StartTime,
		EndTime:   timeslotDetails.EndTime,
	}
	return timeslot
}

func (u *Service) editService(serviceDetails Service) Service {
	// Mock Function
	u = &serviceDetails
	return *u
}
func (u *Service) calculateAvgRating(newRating int) int {
	return 0
}
