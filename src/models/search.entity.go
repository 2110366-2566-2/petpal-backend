package models

import "time"

// Q for Services's User name, Service name, Service type
// SortBy for sorting by price, rating
type SearchHistory struct {
	Q               string    `json:"q" bson:"q"`
	Location        string    `json:"name" bson:"name"`
	StartTime       time.Time `json:"gender" bson:"gender"`
	EndTime         time.Time `json:"age" bson:"age"`
	StartPriceRange float64   `json:"type" bson:"type"`
	EndPriceRange   float64   `json:"healthInformation" bson:"healthInformation"`
	MinRating       float64   `json:"certificate" bson:"certificate"`
	MaxRating       float64   `json:"behaviouralNotes" bson:"behaviouralNotes"`
	PageNumber      int       `json:"vaccinations" bson:"vaccinations"`
	PageSize        int       `json:"dietyPreferences" bson:"dietyPreferences"`
	SortBy          string    `json:"breed" bson:"breed"`
}
