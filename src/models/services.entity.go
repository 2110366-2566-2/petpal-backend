// user.go
package models

// User represents a user entity
type Service struct {
	ServiceID          string     `json:"serviceID" bson:"serviceID"`
	ServiceName        string     `json:"serviceName" bson:"serviceName"`
	ServiceType        string     `json:"serviceType" bson:"serviceType"`
	ServiceDescription string     `json:"serviceDescription" bson:"serviceDescription"`
	ServiceImg         []byte     `json:"serviceImg" bson:"serviceImg"`
	AverageRating      float64    `json:"averageRating" bson:"averageRating"`
	RequireCert        bool       `json:"requireCert" bson:"requireCert"`
	Timeslots          []Timeslot `json:"timeslots" bson:"timeslots"`
	Price              float64    `json:"price" bson:"price"`
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
