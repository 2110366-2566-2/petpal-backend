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

type UpdateSVCP struct {
	// Define the 10 fields here
	SVCPUsername      string `json:"SVCPUsername" bson:"SVCPUsername"`
	Description       string `json:"description" bson:"description"`
	PhoneNumber       string `json:"phoneNumber" bson:"phoneNumber"`
	Address           string `json:"address" bson:"address"`
	SVCPAdditionalImg []byte `json:"SVCPAdditionalImg" bson:"SVCPAdditionalImg"`
	SVCPImg           []byte `json:"SVCPImg" bson:"SVCPImg"`
}

type SVCP struct {
	Individual
	SVCPID                string    `json:"SVCPID" bson:"SVCPID,omitempty"`
	SVCPImg               []byte    `json:"SVCPImg" bson:"SVCPImg"`
	SVCPUsername          string    `json:"SVCPUsername" bson:"SVCPUsername"`
	SVCPPassword          string    `json:"SVCPPassword" bson:"SVCPPassword"`
	SVCPEmail             string    `json:"SVCPEmail" bson:"SVCPEmail"`
	IsVerified            bool      `json:"isVerified" bson:"isVerified"`
	SVCPResponsiblePerson string    `json:"SVCPResponsiblePerson" bson:"SVCPResponsiblePerson"`
	DefaultBank           string    `json:"defaultBank" bson:"defaultBank"`
	DefaultAccountNumber  string    `json:"defaultAccountNumber" bson:"defaultAccountNumber"`
	License               string    `json:"license" bson:"license"`
	Address               string    `json:"address" bson:"address"`
	PhoneNumber           string    `json:"phoneNumber" bson:"phoneNumber"`
	Description           string    `json:"description" bson:"description"`
	SVCPAdditionalImg     []byte    `json:"SVCPAdditionalImg" bson:"SVCPAdditionalImg"`
	SVCPServiceType       string    `json:"SVCPServiceType" bson:"SVCPServiceType"`
	Services              []Service `json:"services" bson:"services"`
}

func (e *SVCP) validate(SVCPImg string) bool {
	// Mock Function
	return true
}
func (e *SVCP) ImplementCurrentEntity() {} // Empty stub if no shared methods

func (e *SVCP) RemoveSensitiveData() {
	// remove password
	e.SVCPPassword = ""
}

func (e *SVCP) UpdateField(key string, value any) {
	// UpdateField
	// get the field and update it
	// return the updated service
	if key == "SVCPUsername" {
		e.SVCPUsername = value.(string)
	}
	if key == "description" {
		e.Description = value.(string)
	}
	if key == "phoneNumber" {
		e.PhoneNumber = value.(string)
	}
	if key == "address" {
		e.Address = value.(string)
	}
	if key == "SVCPAdditionalImg" {
		e.SVCPAdditionalImg = value.([]byte)
	}
	if key == "SVCPImg" {
		e.SVCPImg = value.([]byte)
	}
}
