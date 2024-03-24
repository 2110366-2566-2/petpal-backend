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

func (u *Service) UpdateField(key string, value any) error {
	// UpdateField
	// get the field and update it
	// return the updated service
	if key == "serviceName" {
		u.ServiceName = value.(string)
	}
	if key == "servicesImg" {
		e := value.(string)
		u.ServiceImg = []byte(e)
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
	if key == "timeslots" {
		currentTimeslots := u.Timeslots
		var Timeslots []Timeslot
		for i, timeslot := range value.([]interface{}) {
			var singleTimeslot Timeslot
			if i < len(currentTimeslots) {
				singleTimeslot = currentTimeslots[i]
			} else {
				singleTimeslot.SetDefult()
			}

			for k, v := range timeslot.(map[string]interface{}) {
				err := singleTimeslot.UpdateField(k, v)
				if err != nil {
					return err
				}
			}
			Timeslots = append(Timeslots, singleTimeslot)
		}
		u.Timeslots = Timeslots
	}
	return nil
}

func (u *Service) calculateAvgRating(newRating int) int {
	return 0
}
