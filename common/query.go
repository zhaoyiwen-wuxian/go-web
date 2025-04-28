package common

import (
	buildconditionsmap "go-web/buildConditionsMap"
	"reflect"

	"gorm.io/gorm"
)

// 2. 通用查询单个数据：根据条件查询单一数据
func QueryOne(db *gorm.DB, model interface{}, out interface{}, config buildconditionsmap.ConditionConfig, joins []string) error {
	conditions := buildconditionsmap.BuildConditionsMapWithConfig(model, config)
	query, err := BuildQueryModel(db, model, conditions)
	if err != nil {
		return err
	}

	if joins != nil {
		for _, join := range joins {
			query = query.Joins(join)
		}
	}
	// 查询单条数据
	if err := query.First(out).Error; err != nil {
		return err
	}

	return nil
}

// 3. 通用插入：插入新数据
func Create(db *gorm.DB, model interface{}) error {
	if err := db.Create(model).Error; err != nil {
		return err
	}
	return nil
}

// 4. 通用更新：支持根据条件更新数据
func Update(db *gorm.DB, model interface{}, conditions map[string]interface{}, updates interface{}) error {
	query, err := BuildQueryModel(db, model, conditions)
	if err != nil {
		return err
	}

	// 更新字段
	updateMap, err := StructToMap(updates) // 确保你有这个函数来处理结构体到 map 的转换
	if err != nil {
		return err
	}
	if err := query.UpdateColumns(updateMap).Error; err != nil {
		return err
	}

	return nil
}

// 5. 通用删除：根据条件删除记录
func Delete(db *gorm.DB, model interface{}, conditions map[string]interface{}) error {
	query, err := BuildQueryModel(db, model, conditions)
	if err != nil {
		return err
	}

	// 删除记录
	if err := query.Delete(model).Error; err != nil {
		return err
	}

	return nil
}

// 6. 不分页查询，获取全部数据
func QueryAll(db *gorm.DB, model interface{}, out interface{}, config buildconditionsmap.ConditionConfig, sortFields []SortField, joins []string) error {
	conditions := buildconditionsmap.BuildConditionsMapWithConfig(model, config)
	query, err := BuildQueryModel(db, model, conditions)
	if err != nil {
		return err
	}

	// 排序
	if sortFields != nil {
		query = ApplySorting(query, sortFields)

	}

	// 使用 Joins 支持链表查询
	if joins != nil {
		for _, join := range joins {
			query = query.Joins(join)
		}
	}

	// 查询全部数据
	if err := query.Find(out).Error; err != nil {
		return err
	}

	return nil
}

// 7. 批量插入：插入多条记录
func BulkCreate(db *gorm.DB, models interface{}) error {
	if err := db.Create(models).Error; err != nil {
		return err
	}
	return nil
}

// 8. 批量更新：更新多条记录
func BulkUpdate(db *gorm.DB, model interface{}, conditions map[string]interface{}, updates map[string]interface{}) error {
	query, err := BuildQueryModel(db, model, conditions)
	if err != nil {
		return err
	}

	if err := query.UpdateColumns(updates).Error; err != nil {
		return err
	}
	return nil
}

// 9. 支持事务的批量操作：例如批量插入或更新
func ExecuteTransaction(db *gorm.DB, operations func(tx *gorm.DB) error) error {
	tx := db.Begin()

	// 确保在操作中出现错误时回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 执行传入的操作
	if err := operations(tx); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// 获取数据库中的字段名（辅助方法，方便根据字段名生成查询）
func GetColumnNames(db *gorm.DB, model interface{}) ([]string, error) {
	var columns []string

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	tableName := db.NamingStrategy.TableName(modelType.Name())
	query := `
	SELECT COLUMN_NAME
	FROM INFORMATION_SCHEMA.COLUMNS
	WHERE TABLE_NAME = ? AND TABLE_SCHEMA = DATABASE()`
	if err := db.Raw(query, tableName).Scan(&columns).Error; err != nil {
		return nil, err
	}

	return columns, nil
}

// 单一数据更新：根据ID更新数据
func UpdateOneByID(db *gorm.DB, model interface{}, id interface{}, updates map[string]interface{}) error {
	// 根据ID查询记录
	query := db.Model(model).Where("id = ?", id)

	// 更新数据
	if err := query.UpdateColumns(updates).Error; err != nil {
		return err
	}

	return nil
}

// 执行复杂的 SQL 查询
func ExecuteComplexQuery(db *gorm.DB, query string, args []interface{}, out interface{}) error {
	if err := db.Raw(query, args...).Scan(out).Error; err != nil {
		return err
	}
	return nil
}
