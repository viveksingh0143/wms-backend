package ability

type Form struct {
	ID     uint   `json:"id" binding:"-"`
	Name   string `json:"name" validate:"required,min=4,max=100"`
	Module string `json:"module" validate:"required,min=4,max=100"`
}
