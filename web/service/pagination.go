package service

import (
	"gorm.io/gorm"
)

type Page struct {
	PageNr   int
	PageSize int
}

func Paginate(p *Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p == nil {
			return db
		}

		if p.PageNr == 0 {
			p.PageNr = 1
		}

		offset := (p.PageNr - 1) * p.PageSize
		return db.Offset(offset).Limit(p.PageSize)
	}
}
