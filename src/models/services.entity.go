// user.go
package models

// Create Service
type CreateService struct {
	ServiceName        string  `json:"serviceName" bson:"serviceName"`
	ServiceType        string  `json:"serviceType" bson:"serviceType"`
	ServiceDescription string  `json:"serviceDescription" bson:"serviceDescription"`
	Price              float64 `json:"price" bson:"price"`
	Timeslots          []CreateTimeslot
}

// User represents a user entity
type Service struct {
	ServiceID          string     `json:"serviceID" bson:"serviceID,omitempty"`
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
	u = &serviceDetails
	return *u
}
func (u *Service) UpdateField(key string, value any) Service {
	// UpdateField
	// get the field and update it
	// return the updated service
	if key == "serviceName" {
		u.ServiceName = value.(string)
	}
	if key == "serviceType" {
		u.ServiceType = value.(string)
	}
	if key == "serviceDescription" {
		u.ServiceDescription = value.(string)
	}
	if key == "averageRating" {
		u.AverageRating = value.(float64)
	}
	if key == "requireCert" {
		u.RequireCert = value.(bool)
	}
	if key == "price" {
		u.Price = value.(float64)
	}
	return *u

}

func (u *Service) calculateAvgRating(newRating int) int {
	return 0
}
