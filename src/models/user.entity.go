// user.go
package models

import (
	"time"
)

// User represents a user entity
type CreateUser struct {
	// Define the 10 fields here
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
}

type User struct {
	Individual
	Username             string    `json:"username"`
	Password             string    `json:"password"`
	Email                string    `json:"email"`
	FullName             string    `json:"fullName"`
	Address              string    `json:"address"`
	DateOfBirth          time.Time `json:"dateOfBirth"`
	PhoneNumber          string    `json:"phoneNumber"`
	ProfilePicture       string    `json:"profilePicture"`
	DefaultAccountNumber string    `json:"defaultAccountNumber"`
	DefaultBank          string    `json:"defaultBank"`
	Pets                 []Pet     `json:"pets"`
}
type LoginReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	LoginType string `json:logintype`
}

type LoginRes struct {
	AccessToken string `่json:accesstoken`
	UserEmail   string `json:"useremail"`
	LoginType   string `json:logintype`
}

func (u *User) editPet(petName string, petDetails Pet) Pet {
	// Mock Function
	return petDetails
}
