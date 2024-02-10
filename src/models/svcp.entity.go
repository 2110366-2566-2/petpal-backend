// user.go
package models

// User represents a user entity
type SVCP struct {
	Individual
	SVCPID                string `json:"SVCPID"`
	SVCPImg               string `json:"SVCPImg"`
	SVCPUsername          string `json:"SVCPUsername"`
	SVCPPassword          string `json:"SVCPPassword"`
	SVCPEmail             string `json:"SVCPEmail"`
	isVerified            bool   `json:"isVerified"`
	SVCPResponsiblePerson string `json:"SVCPResponsiblePerson"`
	defaultBank           string `json:"defaultBank"`
	defaultAccountNumber  string `json:"defaultAccountNumber"`
	license               string `json:"license"`
	location              string `json:"location"`
	SVCPAdditionalImg     string `json:"SVCPAdditionalImg"`
}

func (e *SVCP) validate(SVCPImg string) bool {
	// Mock Function
	return true
}
