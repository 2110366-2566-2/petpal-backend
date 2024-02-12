package models

// Pet represents a pet entity
type Pet struct {
	OwnerUsername     string `json:ownerUsername"`
	Name              string `json:name"`
	Gender            string `json:gender"`
	Age               int    `json:age"`
	Pet_type          string `json:type"`
	HealthInformation string `json:healthInformation"`
	Certificate       string `json:certificate"`
	BehaviouralNotes  string `json:behaviouralNotes"`
	Vaccinations      string `json:vaccinations"`
	DietyPreferences  string `json:dietyPreferences"`
}
