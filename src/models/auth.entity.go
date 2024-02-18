// user.go
package models

type LoginReq struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	LoginType string `json:"logintype"`
}

type LoginRes struct {
	AccessToken string `่json:"accesstoken"`
	UserEmail   string `json:"useremail"`
	LoginType   string `json:"logintype"`
}

type BasicRes struct {
	Message string `json:"message"`
}
type BasicErrorRes struct {
	Error string `json:"error"`
}