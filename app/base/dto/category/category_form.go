package category

import (
	"star-wms/core/common/dto"
	"star-wms/core/types"
)

type Form struct {
	dto.BaseDto
	Name     string       `json:"name" validate:"required,min=3,max=100"`
	Slug     string       `json:"slug"`
	FullPath string       `json:"full_path"`
	Status   types.Status `json:"status" validate:"required,gt=0"`
	Children []*Form      `json:"children"`
	Parent   *Form        `json:"parent" validationTag:"parent.id" validate:"omitempty,validRelationID,structonly"`
}
