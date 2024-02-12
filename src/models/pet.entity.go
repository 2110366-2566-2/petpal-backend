package models

// Pet represents a pet entity
type Pet struct {
	OwnerUsername     string `json:"ownerUsername" bson:"ownerUsername"`
	Name              string `json:"name" bson:"name"`
	Gender            string `json:"gender" bson:"gender"`
	Age               int	 `json:"age" bson:"age"`
	Pet_type          string `json:"type" bson:"type"`
	HealthInformation string `json:"healthInformation" bson:"healthInformation"`
	Certificate		  string `json:"certificate" bson:"certificate"`
	BehaviouralNotes  string `json:"behaviouralNotes" bson:"behaviouralNotes"`
	Vaccinations      string `json:"vaccinations" bson:"vaccinations"`
	DietyPreferences  string `json:"dietyPreferences" bson:"dietyPreferences"`
}
