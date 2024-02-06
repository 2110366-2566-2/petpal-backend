// user.go
package models

// User represents a user entity
type User struct {
	Individual
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
	pet                  Pet
}

func (u *User) editPet(petName string, petDetails Pet) Pet {
	// Mock Function
	return u.pet
}
