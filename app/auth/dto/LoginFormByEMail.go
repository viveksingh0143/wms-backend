package dto

type LoginFormByEMail struct {
	EMail      string `json:"email" validate:"required,email,min=4,max=100"`
	Password   string `json:"password" validate:"required,min=4,max=100"`
	RememberMe bool   `json:"remember_me"`
}
