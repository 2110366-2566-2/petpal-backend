// user.go
package models

// User represents a user entity
type SVCP struct {
	Individual
	SVCPID                string `json:"SVCPID" bson:"SVCPID"`
	SVCPImg               string `json:"SVCPImg" bson:"SVCPImg"`
	SVCPUsername          string `json:"SVCPUsername" bson:"SVCPUsername"`
	SVCPPassword          string `json:"SVCPPassword" bson:"SVCPPassword"`
	SVCPEmail             string `json:"SVCPEmail" bson:"SVCPEmail"`
	isVerified            bool   `json:"isVerified"`
	SVCPResponsiblePerson string `json:"SVCPResponsiblePerson" bson:"SVCPResponsiblePerson"`
	defaultBank           string `json:"defaultBank"`
	defaultAccountNumber  string `json:"defaultAccountNumber"`
	license               string `json:"license"`
	location              string `json:"location"`
	SVCPAdditionalImg     string `json:"SVCPAdditionalImg" bson:"SVCPAdditionalImg"`
}

func (e *SVCP) validate(SVCPImg string) bool {
	// Mock Function
	return true
}
