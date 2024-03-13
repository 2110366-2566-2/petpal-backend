// user.go
package models

import (
	"time"
)

// User represents a user entity
type CreateUser struct {
	// Define the 10 fields here
	Username    string    `json:"username" bson:"username"`
	Password    string    `json:"password" bson:"password"`
	Email       string    `json:"email" bson:"email"`
	FullName    string    `json:"fullName" bson:"fullName"`
	Address     string    `json:"address" bson:"address"`
	DateOfBirth time.Time `json:"dateOfBirth" bson:"dateOfBirth"`
	PhoneNumber string    `json:"phoneNumber" bson:"phoneNumber"`
}

type User struct {
	Individual
	ID                   string    `json:"id" bson:"_id,omitempty"`
	Username             string    `json:"username" bson:"username"`
	Password             string    `json:"password" bson:"password"`
	Email                string    `json:"email" bson:"email"`
	FullName             string    `json:"fullName" bson:"fullName"`
	Address              string    `json:"address" bson:"address"`
	DateOfBirth          time.Time `json:"dateOfBirth" bson:"dateOfBirth"`
	PhoneNumber          string    `json:"phoneNumber" bson:"phoneNumber"`
	ProfilePicture       []byte    `json:"profilePicture" bson:"profilePicture"`
	DefaultAccountNumber string    `json:"defaultAccountNumber" bson:"defaultAccountNumber"`
	DefaultBank          string    `json:"defaultBank" bson:"defaultBank"`
	Pets                 []Pet     `json:"pets" bson:"pets"`
}

type UserSearchHistory struct {
	User          User           `json:"user" bson:"user"`
	SearchHistory []SearchFilter `json:"search_history" bson:"search_history"`
}

func (u *User) editPet(petName string, petDetails Pet) Pet {
	// Mock Function
	return petDetails
}

func (u *User) ImplementCurrentEntity() {} // Empty stub if no shared methods

func (u *User) RemoveSensitiveData() {
	// remove password
	u.Password = ""
	u.DefaultBank = ""
	u.DefaultAccountNumber = ""
}
