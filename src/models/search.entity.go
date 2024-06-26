package models

import "time"

// Q for Services's User name, Service name, Service type
// SortBy for sorting by price, rating, name
type SearchFilter struct {
	Q               string    `json:"q" bson:"q"`
	ServicesType    string    `json:"services_type" bson:"services_type"`
	Address         string    `json:"address" bson:"address"`
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

type SearchHistory struct {
	Date         time.Time    `json:"date" bson:"date"`
	SearchFilter SearchFilter `json:"search_filters" bson:"search_filters"`
}

type SearchResult struct {
	Services        Service `json:"services"`
	Address         string  `json:"address"`
	SVCPUsername    string  `json:"SVCPUsername"`
	SVCPServiceType string  `json:"SVCPServiceType"`
	SVCPID          string  `json:"SVCPID"`
}
