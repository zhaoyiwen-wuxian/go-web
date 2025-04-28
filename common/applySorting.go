package common

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type SortField struct {
	Field     string // 排序字段
	Direction string // 排序方向（"ASC" 或 "DESC"）
}

// 默认排序方向为升序
const DefaultSortDirection = "ASC"

// 构建排序条件（ORDER BY）
func BuildOrderClause(sortFields []SortField) string {
	var orderClauses []string

	// 遍历排序字段，构建排序条件
	for _, sortField := range sortFields {
		// 如果没有指定排序方向，默认使用升序
		direction := strings.ToUpper(sortField.Direction)
		if direction != "ASC" && direction != "DESC" {
			direction = "ASC"
		}
		// 格式化排序字段
		orderClauses = append(orderClauses, fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}

	// 将多个排序条件用逗号连接
	return strings.Join(orderClauses, ", ")
}

// 应用排序到查询
func ApplySorting(db *gorm.DB, sortFields []SortField) *gorm.DB {
	if len(sortFields) > 0 {
		// 生成排序条件并应用到查询中
		orderClause := BuildOrderClause(sortFields)
		db = db.Order(orderClause)
	}
	return db
}
