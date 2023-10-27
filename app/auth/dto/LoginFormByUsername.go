package dto

type LoginFormByUsername struct {
	Username   string `json:"username" validate:"required,min=4,max=100"`
	Password   string `json:"password" validate:"required,min=4,max=100"`
	RememberMe bool   `json:"remember_me"`
}
