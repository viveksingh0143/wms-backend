package rmbatch

type Filter struct {
	Query     string `form:"query" db:"batch_no,po_category,so_number" whereType:"like" binding:"omitempty,max=100"`
	ID        uint   `form:"id" db:"id" binding:"omitempty,gt=0"`
	StoreID   uint   `form:"store_id" db:"store_id" binding:"omitempty,gt=0"`
	ProductID uint   `form:"product_id" db:"product_id" binding:"omitempty,gt=0"`
}
