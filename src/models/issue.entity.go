package models

import "time"

type Issue struct {
	IssueID        string    `json:"issueID" bson:"_id"`
	IssueDate      time.Time `json:"issueDate" bson:"issueDate"`
	IsResolved     bool      `json:"isResolved" bson:"isResolved"`
	WorkingAdminID string    `json:"workingAdminID" bson:"workingAdminID"`
	ResolveDate    time.Time `json:"resolveDate,omitempty" bson:"resolveDate,omitempty"`
	CreateIssue
}

type CreateIssue struct {
	ReporterID          string `json:"reporterID" bson:"reporterID"`
	ReporterType        string `json:"reporterType" bson:"reporterType"` // user, svcp
	Details             string `json:"details" bson:"details"`
	AttachedImg         []byte `json:"attachedImg" bson:"attachedImg"`
	IssueType           string `json:"issueType" bson:"issueType"` // "refund", "system", "service"
	AssociatedBookingID string `json:"associatedBookingID" bson:"associatedBookingID"`
}
