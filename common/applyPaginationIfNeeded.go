package common

import (
	"go-web/appResponse"
	buildconditionsmap "go-web/buildConditionsMap"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取上下文中的分页参数，并应用分页
func ApplyPaginationIfNeeded(db *gorm.DB, c *gin.Context) (*gorm.DB, int, int, error) {
	// 从上下文中获取分页参数
	page, _ := c.Get("page")
	pageSize, _ := c.Get("pageSize")

	// 如果分页参数存在，应用分页
	if page != nil && pageSize != nil {
		offset := (page.(int) - 1) * pageSize.(int)
		db = db.Offset(offset).Limit(pageSize.(int))
	}

	return db, page.(int), pageSize.(int), nil
}

// 通用查询：返回查询的分页数据
func QueryAllWithPagination(db *gorm.DB, model interface{}, out interface{}, config buildconditionsmap.ConditionConfig, sortFields []SortField, c *gin.Context, joins []string) (appResponse.PaginatedResult, error) {
	// 先构建查询模型
	conditions := buildconditionsmap.BuildConditionsMapWithConfig(model, config)
	query, err := BuildQueryModel(db, model, conditions)
	if err != nil {
		return appResponse.PaginatedResult{}, err
	}
	if joins != nil {
		for _, join := range joins {
			query = query.Joins(join)
		}
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return appResponse.PaginatedResult{}, err
	}
	// 应用分页
	query, page, pageSize, err := ApplyPaginationIfNeeded(query, c)
	if err != nil {
		return appResponse.PaginatedResult{}, err
	}

	// 排序
	if sortFields != nil {
		query = ApplySorting(query, sortFields)
	}

	// 执行查询
	if err := query.Find(out).Error; err != nil {
		return appResponse.PaginatedResult{}, err
	}

	return appResponse.PaginatedResult{
		Data:       out,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
	}, nil
}
