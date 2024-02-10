// user.go
package models

// User represents a user entity
type Service struct {
	SVCPID             string
	serviceType        string
	serviceDescription string
	serviceImage       string
	averageRating      string
	requireCert        string
}

func (u *Service) createTimeslot(timeslotDetails Timeslot) Timeslot {
	timeslot := Timeslot{
		SVCPID:      timeslotDetails.SVCPID,
		serviceType: timeslotDetails.serviceType,
		startTime:   timeslotDetails.startTime,
		endTime:     timeslotDetails.endTime,
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
