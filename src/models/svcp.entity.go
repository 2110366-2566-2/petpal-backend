// user.go
package models

// User represents a user entity
type SVCP struct {
	Individual
	SVCPImg               string
	SVCPUsername          string
	SVCPPassword          string
	SVCPEmail             string
	isVerified            bool
	SVCPResponsiblePerson string
	defaultBank           string
	defaultAccountNumber  string // ควรเปลี่ยนเป็น date หรือเปล่า
	license               string
	location              string
	SVCPAdditionalImg     Pet
}

func (e *SVCP) validate(SVCPImg string) bool {
	// Mock Function
	return true
}
