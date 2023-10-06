package requests

type Pagination struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}

type Sorting struct {
	OrderBy string `form:"order_by"`
	Desc    bool   `form:"desc,default=false"`
}
