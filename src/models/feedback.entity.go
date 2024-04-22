// user.go
package models

// User represents a user entity
type Feedback struct {
	FeedbackID string  `json:"feedbackID" bson:"feedbackID"`
	Rating     float64 `json:"rating" bson:"rating"`
	Content    string  `json:"content" bson:"content"`
}
