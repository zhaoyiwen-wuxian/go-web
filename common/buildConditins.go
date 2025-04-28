package common

import (
	"fmt"

	"gorm.io/gorm"
)

func BuildConditions(db *gorm.DB, exact, gt, lt, like map[string]interface{}) *gorm.DB {
	query := db
	for k, v := range exact {
		query = query.Where(k+" = ?", v)
	}
	for k, v := range gt {
		query = query.Where(k+" > ?", v)
	}
	for k, v := range lt {
		query = query.Where(k+" < ?", v)
	}
	for k, v := range like {
		query = query.Where(k+" LIKE ?", "%"+fmt.Sprintf("%v", v)+"%")
	}
	return query
}
