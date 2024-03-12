package models

import "time"

// Q for Services's User name, Service name, Service type
// SortBy for sorting by price, rating

// @Param q query string false "Search query "
// @Param location query string false "Location"
// @Param start_time query string false "start_time"
// @Param end_time query string false "end_time"
// @Param start_price_range query string false "Start price range"
// @Param end_price_range query string false "End price range"
// @Param min_rating query string false "Minimum rating"
// @Param max_rating query string false "Maximum rating"
// @Param page_number query string false "Page number"
// @Param page_size query string false "Page size"
// @Param sort_by query  false "Sort by (price, rating)"
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
	SortBy          string    `json:"sort_by" bson:"sort_by"`
}
