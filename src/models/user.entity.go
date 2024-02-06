// user.go
package models

// User represents a user entity
type CreateUser struct {
	// Define the 10 fields here
	Username    string
	Password    string
	Address     string
	PhoneNumber string
	Email       string
	FullName    string
	DateOfBirth string // ควรเปลี่ยนเป็น date หรือเปล่า
}

type User struct {
	Individual
	CreateUser
	PhoneNumber          string
	ProfilePicture       string
	DefaultAccountNumber string
	DefaultBank          string
	Pet                  *Pet
}

func (u *User) editPet(petName string, petDetails Pet) Pet {
	// Mock Function
	return *u.Pet
}
