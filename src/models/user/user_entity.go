// user.go
package models

// User represents a user entity
type User struct {
	username             string
	password             string
	address              string
	phoneNumber          string
	email                string
	profilePicture       string
	fullName             string
	dateOfBirth          string // ควรเปลี่ยนเป็น date หรือเปล่า
	defaultAccountNumber string
	defaultBank          string
	// Waiting for implement eneity pet
}
