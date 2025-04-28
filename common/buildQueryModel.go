package common

import (
	"go-web/appError"

	"gorm.io/gorm"
)

// 创建一个包含 Model 和 条件 的 GORM 查询
func BuildQueryModel(db *gorm.DB, model interface{}, conditions map[string]interface{}) (*gorm.DB, error) {
	if model == nil {
		// 返回一个错误：model 不能为空
		return nil, appError.ErrModelNil
	}

	query := db.Model(model)

	if conditions != nil {
		query = BuildQueryWithConditions(query, conditions)
	}

	return query, nil
}
