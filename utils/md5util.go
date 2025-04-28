package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CalculateMD5(data string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	return string(bytes)
}

// 验证密码是否匹配（用户输入密码 vs 数据库存储的哈希密码）
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
