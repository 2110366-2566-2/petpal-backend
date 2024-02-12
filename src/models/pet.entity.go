package models

// Pet represents a pet entity
type Pet struct {
	ownerUsername     string
	name              string
	gender            string
	age               int
	pet_type          string
	healthInformation string
	certificates      string
	behaviouralNotes  string
	vaccinations      string
	dietyPreferences  string
}
