// user.go
package models

// User represents a user entity

type CreateSVCP struct {
	// Define the 10 fields here
	SVCPUsername    string `json:"SVCPUsername" bson:"SVCPUsername"`
	SVCPPassword    string `json:"SVCPPassword" bson:"SVCPPassword"`
	SVCPEmail       string `json:"SVCPEmail" bson:"SVCPEmail"`
	SVCPServiceType string `json:"SVCPServiceType" bson:"SVCPServiceType"`
}
type SVCP struct {
	Individual
	SVCPID                string `json:"SVCPID" bson:"SVCPID"`
	SVCPImg               string `json:"SVCPImg" bson:"SVCPImg"`
	SVCPUsername          string `json:"SVCPUsername" bson:"SVCPUsername"`
	SVCPPassword          string `json:"SVCPPassword" bson:"SVCPPassword"`
	SVCPEmail             string `json:"SVCPEmail" bson:"SVCPEmail"`
	IsVerified            bool   `json:"isVerified" bson:"isVerified"`
	SVCPResponsiblePerson string `json:"SVCPResponsiblePerson" bson:"SVCPResponsiblePerson"`
	DefaultBank           string `json:"defaultBank" bson:"defaultBank"`
	DefaultAccountNumber  string `json:"defaultAccountNumber" bson:"defaultAccountNumber"`
	License               string `json:"license" bson:"license"`
	Location              string `json:"location" bson:"location"`
	Description		      string `json:"description" bson:"description"`
	SVCPAdditionalImg     string `json:"SVCPAdditionalImg" bson:"SVCPAdditionalImg"`
	SVCPServiceType	      string `json:"SVCPServiceType" bson:"SVCPServiceType"`
	Services 			  []Service `json:"services" bson:"services"`
}

func (e *SVCP) validate(SVCPImg string) bool {
	// Mock Function
	return true
}
