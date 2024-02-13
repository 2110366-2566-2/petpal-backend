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
	ID                   string    `json:"id" bson:"_id"`
	Username             string    `json:"username" bson:"username"`
	Password             string    `json:"password" bson:"password"`
	Email                string    `json:"email" bson:"email"`
	FullName             string    `json:"fullName" bson:"fullName"`
	Address              string    `json:"address" bson:"address"`
	DateOfBirth          time.Time `json:"dateOfBirth" bson:"dateOfBirth"`
	PhoneNumber          string    `json:"phoneNumber" bson:"phoneNumber"`
	ProfilePicture       string    `json:"profilePicture" bson:"profilePicture"`
	DefaultAccountNumber string    `json:"defaultAccountNumber" bson:"defaultAccountNumber"`
	DefaultBank          string    `json:"defaultBank" bson:"defaultBank"`
	Pets                 []Pet     `json:"pets" bson:"pets"`
}
type LoginReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	LoginType string `json:"logintype"`
}

type LoginRes struct {
	AccessToken string `่json:"accesstoken"`
	UserEmail   string `json:"useremail"`
	LoginType   string `json:"logintype"`
}

func (u *User) editPet(petName string, petDetails Pet) Pet {
	// Mock Function
	return petDetails
}
