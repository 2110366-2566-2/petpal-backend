// user.go
package models

import (
	"time"
)

// User represents a user entity
type CreateUser struct {
	// Define the 10 fields here
	Username    string `json:"username"`
	Password    string `json:"password"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	FullName    string `json:"fullName"`
	DateOfBirth time.Time `json:"dateOfBirth"`
}

type User struct {
	Individual
	CreateUser
	PhoneNumber          string `json:"phoneNumber"`
	ProfilePicture       string `json:"profilePicture"`
	DefaultAccountNumber string `json:"defaultAccountNumber"`
	DefaultBank          string `json:"defaultBank"`
	Pets                 []Pet  `json:"pets"`
}

func (u *User) editPet(petName string, petDetails Pet) Pet {
	// Mock Function
	return petDetails
}
