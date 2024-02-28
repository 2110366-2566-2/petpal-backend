package models

import (
	"time"
)

type SearchHistory struct {
	Date          time.Time            `json:"date" bson:"date"`
	SearchFilters ServiceSearchFilters `json:"search_filters" bson:"search_filters"`
}

type ServiceSearchFilters struct {
	ServiceName        string            `json:"serviceName" bson:"serviceName" optional:"true"`
	ServiceType        string            `json:"serviceType" bson:"serviceType" optional:"true"`
	ServiceDescription string            `json:"serviceDescription" bson:"serviceDescription" optional:"true"`
	Timeslot           time.Time         `json:"timeslot" bson:"timeslot" optional:"true"`
	AverageRating      NumericComparison `json:"averageRating" bson:"averageRating" optional:"true"`
}

type NumericComparison struct {
	GreaterThan float64 `json:"greaterThan" bson:"$gt" optional:"true"`
	LessThan    float64 `json:"lessThan" bson:"$lt" optional:"true"`
}
