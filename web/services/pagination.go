package services

import "gorm.io/gorm"

type Page struct {
	Number int
	Size   int
}

func Paginate(p *Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p == nil {
			return db
		}

		if p.Number == 0 {
			p.Number = 1
		}

		offset := (p.Number - 1) * p.Size
		return db.Offset(offset).Limit(p.Size)
	}
}
