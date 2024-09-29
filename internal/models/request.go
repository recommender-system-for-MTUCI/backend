package models

type RegistrationRequst struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailRequest struct {
	Code int `json:"code"`
}
