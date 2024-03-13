package models

import "time"

// Q for Services's User name, Service name, Service type
// SortBy for sorting by price, rating, name
type SearchHistory struct {
	Q               string    `json:"q" bson:"q"`
	Location        string    `json:"location" bson:"location"`
	StartTime       time.Time `json:"start_time" bson:"start_time"`
	EndTime         time.Time `json:"end_time" bson:"end_time"`
	StartPriceRange float64   `json:"start_price_range" bson:"start_price_range"`
	EndPriceRange   float64   `json:"end_price_range" bson:"end_price_range"`
	MinRating       float64   `json:"min_rating" bson:"min_rating"`
	MaxRating       float64   `json:"max_rating" bson:"max_rating"`
	PageNumber      int       `json:"page_number" bson:"page_number"`
	PageSize        int       `json:"page_size" bson:"page_size"`
	SortBy          string    `json:"sort_by" bson:"sort_by"` // price, rating, name
	Descending      bool      `json:"descending" bson:"descending"`
}

type SearchResult struct {
	Services        Service `json:"services"`
	Location        string  `json:"location"`
	SVCPUsername    string  `json:"SVCPUsername"`
	SVCPServiceType string  `json:"SVCPServiceType"`
}
