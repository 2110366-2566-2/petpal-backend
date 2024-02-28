package models

type Issue struct {
	IssueID      string `json:"issueID" bson:"issueID"`
	ReporterID   string `json:"reporterID" bson:"reporterID"`
	ReporterType string `json:"reporterType" bson:"reporterType"`
	IssueTitle   string `json:"issueTitle" bson:"issueTitle"`
	IssueData    string `json:"issueData" bson:"issueData"`
	IsSolved     bool   `json:"isSolved" bson:"isSolved"`
}
