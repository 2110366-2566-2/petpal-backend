// user.go
package models

// User represents a user entity
type Individual struct {
	IndividualID int
}

func (e *Individual) editProfile(individualID int, profileDetails string) Individual {
	// Mock Function <-------------------------------
	// What is profileDetails Na
	e.IndividualID = individualID
	return *e
}

func (e *Individual) viewChatDashBoard(individualCredential string) []string {
	dashboardItems := []string{"Chat 1", "Chat 2", "Chat 3"}
	return dashboardItems
}
