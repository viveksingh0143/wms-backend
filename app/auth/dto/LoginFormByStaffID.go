package dto

type LoginFormByStaffID struct {
	StaffID    string `json:"staff_id" validate:"required,min=4,max=100"`
	Password   string `json:"password" validate:"required,min=4,max=100"`
	RememberMe bool   `json:"remember_me"`
}
