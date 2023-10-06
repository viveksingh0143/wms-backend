package requests

type RequestIDs struct {
	IDs []uint `json:"ids" validate:"required,min_ids"`
}
