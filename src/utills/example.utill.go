package utills

import (
	"petpal-backend/src/models"
)

// this is the mock data
var mockData = models.Example{
	Message:      "hello world",
	ErrorMessage: "no error",
	Number:       2024,
}

func ExampleUtill() (models.Example, error) {
	return mockData, nil
}
