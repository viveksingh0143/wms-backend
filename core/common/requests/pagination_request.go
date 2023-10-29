package requests

import "strings"

type Pagination struct {
	Start    int `form:"_start,default=0"`
	End      int `form:"_end,default=0"`
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=0"`
}

type Sorting struct {
	OrderBy string `form:"_sort"`
	Order   string `form:"_order,default=desc"`
	Desc    bool
}

func (p *Pagination) CalculatePageAndPageSize() {
	if p.PageSize < 0 {
		// continue
	} else if p.PageSize == 0 {
		p.PageSize = p.End - p.Start
		if p.PageSize > 0 {
			p.Page = p.Start / p.PageSize
		} else {
			p.Page = 0
		}
	} else {
		p.Start = (p.Page - 1) * p.PageSize
		p.End = p.Start + p.PageSize
	}
}

func (s *Sorting) CalculateSorting() {
	if strings.ToLower(s.Order) == "asc" {
		s.Desc = false
	} else {
		s.Desc = true
	}
}
