// user.go
package models

import (
	"time"
)

// User represents a user entity
type CreateUser struct {
	// Define the 10 fields here
	Username    string 		`json:"username" bson:"username"`
	Password    string 		`json:"password" bson:"password"`
	Address     string 		`json:"address" bson:"address"`
	PhoneNumber string 		`json:"phoneNumber" bson:"phoneNumber"`
	Email       string 		`json:"email" bson:"email"`
	FullName    string 		`json:"fullName" bson:"fullName"`
	DateOfBirth time.Time 	`json:"dateOfBirth" bson:"dateOfBirth"`
}

type User struct {
	Individual
	CreateUser					`json:",inline" bson:",inline"`
	ID					 string `json:"id" bson:"_id"`
	ProfilePicture       string `json:"profilePicture" bson:"profilePicture"`
	DefaultAccountNumber string `json:"defaultAccountNumber" bson:"defaultAccountNumber"`
	DefaultBank          string `json:"defaultBank" bson:"defaultBank"`
	Pets                 []Pet  `json:"pets" bson:"pets"`
}

func (u *User) editPet(petName string, petDetails Pet) Pet {
	// Mock Function
	return petDetails
}
