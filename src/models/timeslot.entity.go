// user.go
package models

// User represents a user entity
type Timeslot struct {
	SVCPID      string
	serviceType string
	startTime   string
	endTime     string
}

func (u *Timeslot) editTimeslot(timeslotDetails Timeslot) Timeslot {
	// Mock Function
	return *u
}
