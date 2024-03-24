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
	User          User            `json:"user" bson:"user"`
	SearchHistory []SearchHistory `json:"search_history" bson:"search_history"`
}

func (u *User) UpdateField(key string, value any) User {
	// UpdateField
	// get the field and update it
	// return the updated service
	if key == "id" {
		u.ID = value.(string)
	} else if key == "username" {
		u.Username = value.(string)
	} else if key == "password" {
		u.Password = value.(string)
	} else if key == "email" {
		u.Email = value.(string)
	} else if key == "fullName" {
		u.FullName = value.(string)
	} else if key == "address" {
		u.Address = value.(string)
	} else if key == "dateOfBirth" {
		u.DateOfBirth = value.(time.Time)
	} else if key == "phoneNumber" {
		u.PhoneNumber = value.(string)
	} else if key == "profilePicture" {
		u.ProfilePicture = value.([]byte)
	} else if key == "defaultAccountNumber" {
		u.DefaultAccountNumber = value.(string)
	} else if key == "defaultBank" {
		u.DefaultBank = value.(string)
	} else if key == "pets" {
		u.Pets = value.([]Pet)
	}
	return *u

}
func (u *User) editPet(petName string, petDetails Pet) Pet {
	// Mock Function
	return petDetails
}

func (u *User) ImplementCurrentEntity() {} // Empty stub if no shared methods

func (u *User) RemoveSensitiveData() {
	// remove password
	u.Password = ""
}
