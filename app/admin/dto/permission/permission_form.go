package permission

type Form struct {
	ID         uint   `json:"id" binding:"-"`
	GroupName  string `json:"group_name" validate:"required,min=4,max=100"`
	ModuleName string `json:"module_name" validate:"required,min=4,max=100"`
	ReadPerm   bool   `json:"read_perm"`
	CreatePerm bool   `json:"create_perm"`
	UpdatePerm bool   `json:"update_perm"`
	DeletePerm bool   `json:"delete_perm"`
	ImportPerm bool   `json:"import_perm"`
	ExportPerm bool   `json:"export_perm"`
}
