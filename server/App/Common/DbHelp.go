package Common

import "github.com/jinzhu/gorm"

type DbHelp struct{}

func (DbHelp) ModelByPage(tel *gorm.DB, limit int, page int) *gorm.DB {
	return tel.Offset((page - 1) * limit).Limit(limit)
}
