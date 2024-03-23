package models

type Admin struct {
	Individual
	AdminID   string `json:"adminID" bson:"adminID"`
	FullName  string `json:"fullName" bson:"fullName"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	Username  string `json:"username" bson:"username"`
}

type CreateAdmin struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (e *Admin) RemoveSensitiveData() {
	e.Password = ""
}