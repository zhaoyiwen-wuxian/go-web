package service

import (
	"context"
	"go-web/appError"
	"go-web/appResponse"
	buildconditionsmap "go-web/buildConditionsMap"
	"go-web/common"
	"go-web/models"
	"go-web/redisutil"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SendMessage 发送消息
func SendMessage(db *gorm.DB, from uint, to uint, content string, msgType string, extra string) error {
	message := models.Message{
		Type:      msgType,
		From:      strconv.Itoa(int(from)),
		To:        strconv.Itoa(int(to)),
		Content:   content,
		Timestamp: time.Now().Unix(),
		Extra:     extra,
	}

	if err := common.Create(db, &message); err != nil {
		return appError.NewAppError(500, "发送消息失败", err)
	}

	// 缓存消息
	if message.ID > 0 {
		if err := redisutil.CacheMessage(context.Background(), message.ID, message, 24*time.Hour); err != nil {
			log.Printf("CacheMessage failed: messageID=%d, err=%v", message.ID, err)
		}
	}

	return nil
}

// GetMessages 获取与 userID 有关的消息（发出或接收），并从 Redis 获取缓存优先
func GetMessages(db *gorm.DB, redisCtx context.Context, userID uint, c *gin.Context) (appResponse.PaginatedResult, error) {
	var dbMessages []models.Message
	var finalMessages []models.Message

	userIDStr := strconv.Itoa(int(userID))

	// 构建 from=userID 或 to=userID 的查询条件
	cond := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{"messages": {"from", "to"}},
		ExactFields: []string{"from", "to"},
		FieldValues: map[string]interface{}{
			"from": userIDStr,
			"to":   userIDStr,
		},
	}

	// 排序条件（按时间倒序）
	sortFields := []common.SortField{
		{Field: "timestamp", Direction: "DESC"},
	}

	// 使用分页查询所有消息
	rust, err := common.QueryAllWithPagination(db, &models.Message{}, &dbMessages, cond, sortFields, c, nil)
	if err != nil {
		return appResponse.PaginatedResult{}, appError.NewAppError(500, "获取消息失败", err)
	}

	// 遍历查询结果，优先从缓存获取
	for _, msg := range dbMessages {
		cached, err := redisutil.GetCachedMessage(redisCtx, msg.ID)
		if err == nil {
			finalMessages = append(finalMessages, cached)
		} else {
			finalMessages = append(finalMessages, msg)
			_ = redisutil.CacheMessage(redisCtx, msg.ID, msg, 24*time.Hour)
		}

		// 若设置了 limit，提前截断
		limit := rust.PageSize
		if limit > 0 && len(finalMessages) >= limit {
			break
		}
	}

	// 计算最终消息的总数并返回
	messger_len := len(finalMessages)

	// 返回分页结果，注意 TotalCount 需要转换为 int64
	return appResponse.PaginatedResult{
		Data:       finalMessages,
		Page:       rust.Page,
		PageSize:   rust.PageSize,
		TotalCount: int64(messger_len), // Convert to int64
	}, nil
}

// 获取群组消息
func GetGroupMessages(db *gorm.DB, groupID uint, c *gin.Context) (appResponse.PaginatedResult, error) {
	var finalMessages []models.Message
	var dbMessages []models.Message

	// 创建群组消息查询条件
	conditionConfig := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{
			"message": {"to", "content", "created_at"},
		},
		ExactFields: []string{"message.to"},
		FieldValues: map[string]interface{}{
			"message.to": groupID, // 群消息查询条件
		},
	}

	// 执行数据库查询
	result, err := common.QueryAllWithPagination(db, models.Message{}, &dbMessages, conditionConfig, nil, c, nil)
	if err != nil {
		return appResponse.PaginatedResult{}, err
	}

	// 或者根据需要传递实际的上下文

	// 提前限制数量
	limit := result.PageSize
	if limit > 0 && len(dbMessages) > limit {
		dbMessages = dbMessages[:limit]
	}

	for _, msg := range dbMessages {
		cached, err := redisutil.GetCachedMessage(c, msg.ID)
		if err == nil {
			finalMessages = append(finalMessages, cached)
		} else {
			finalMessages = append(finalMessages, msg)
			_ = redisutil.CacheMessage(c, msg.ID, msg, 24*time.Hour)
		}
	}

	decryptedMessage, err := DecryptMessage(finalMessages, []byte(redisutil.MessageKey))
	if err != nil {
		log.Println("消息解密错误:", err)
		return appResponse.PaginatedResult{}, err
	}

	return appResponse.PaginatedResult{
		Data:       decryptedMessage,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalCount: result.TotalCount,
	}, nil
}

// 获取个人聊天记录
func GetPersonalMessages(db *gorm.DB, userID uint, otherUserID uint, c *gin.Context) (appResponse.PaginatedResult, error) {
	var finalMessages []models.Message
	var dbMessages []models.Message

	// 创建个人消息查询条件
	conditionConfig := buildconditionsmap.ConditionConfig{
		TableFields: map[string][]string{
			"message": {"from", "to", "content", "created_at"},
		},
		ExactFields: []string{"message.from", "message.to"},
		FieldValues: map[string]interface{}{
			"message.from": userID,      // 当前用户发送的消息
			"message.to":   otherUserID, // 接收消息的用户
		},
	}

	// 执行数据库查询
	result, err := common.QueryAllWithPagination(db, models.Message{}, &dbMessages, conditionConfig, nil, c, nil)
	if err != nil {
		return appResponse.PaginatedResult{}, err
	}
	// 或者根据需要传递实际的上下文
	// 提前限制数量
	limit := result.PageSize
	if limit > 0 && len(dbMessages) > limit {
		dbMessages = dbMessages[:limit]
	}

	for _, msg := range dbMessages {
		cached, err := redisutil.GetCachedMessage(c, msg.ID)
		if err == nil {
			finalMessages = append(finalMessages, cached)
		} else {
			finalMessages = append(finalMessages, msg)
			_ = redisutil.CacheMessage(c, msg.ID, msg, 24*time.Hour)
		}
	}

	decryptedMessage, err := DecryptMessage(finalMessages, []byte(redisutil.MessageKey))
	if err != nil {
		log.Println("消息解密错误:", err)
		return appResponse.PaginatedResult{}, err
	}

	return appResponse.PaginatedResult{
		Data:       decryptedMessage,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalCount: result.TotalCount,
	}, nil
}
