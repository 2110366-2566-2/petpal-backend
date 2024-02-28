package models

type Admin struct {
	AdminID       string `json:"adminID" bson:"adminID"`
	AdminName     string `json:"adminName" bson:"adminName"`
	AdminEmail    string `json:"adminEmail" bson:"adminEmail"`
	AdminPassword string `json:"adminPassword" bson:"adminPassword"`
}
