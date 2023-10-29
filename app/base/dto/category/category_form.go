package category

import (
	"star-wms/core/types"
)

type Form struct {
	ID       uint         `json:"id" binding:"-"`
	Name     string       `json:"name" validate:"required,min=4,max=100"`
	Slug     string       `json:"slug"`
	FullPath string       `json:"full_path"`
	Status   types.Status `json:"status" validate:"required,gt=0"`
	Children []*Form      `json:"children"`
	Parent   *Form        `json:"parent" validationTag:"parent.id" validate:"omitempty,validRelationID,structonly"`
}
