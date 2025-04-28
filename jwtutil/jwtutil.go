package jwtutil

import (
	"context"
	"errors"
	"time"

	"go-web/redisutil"
	"go-web/utils"

	"github.com/golang-jwt/jwt/v4"
)

// 设置 JWT 秘钥
var jwtSecret = []byte(utils.Conf.JWT.JWT.Secret)

// Claims 定义了 JWT 的 Payload 部分
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, role string, expiration time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT Token
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("签名错误")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token失效")
	}
	return claims, nil
}

// StoreToken 在 Redis 中缓存 Token
func StoreToken(ctx context.Context, token string, userID uint, expiration time.Duration) error {
	return redisutil.SetToken(ctx, token, userID, expiration)
}

// GetTokenUserID 从 Redis 中获取 Token 对应的 UserID
func GetTokenUserID(ctx context.Context, token string) (string, error) {
	return redisutil.GetToken(ctx, token)
}
