package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// JWTClaims 定义JWT令牌的声明
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"` // 用户角色: student, instructor, admin
	Exp      int64  `json:"exp"`
}

// GenerateJWT 生成JWT令牌
func GenerateJWT(userID, username, role string) (string, error) {
	// 从配置中获取密钥和过期时间
	secret := "your-secret-key-here" // 应该从配置中读取
	expirationSeconds := 86400       // 24小时

	// 创建头部
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	// 序列化头部
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %v", err)
	}

	// 创建载荷
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Exp:      time.Now().Add(time.Duration(expirationSeconds) * time.Second).Unix(),
	}

	// 序列化载荷
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Base64编码头部和载荷
	headerBase64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// 创建签名
	signatureInput := headerBase64 + "." + payloadBase64
	signature := hmacSha256(signatureInput, secret)

	// 组合JWT
	jwt := signatureInput + "." + signature

	return jwt, nil
}

// ValidateJWT 验证JWT令牌并返回声明
func ValidateJWT(tokenString string) (*JWTClaims, error) {
	secret := "your-secret-key-here" // 应该从配置中读取

	// 分割令牌
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	headerBase64 := parts[0]
	payloadBase64 := parts[1]
	signature := parts[2]

	// 验证签名
	signatureInput := headerBase64 + "." + payloadBase64
	expectedSignature := hmacSha256(signatureInput, secret)

	if signature != expectedSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	// 解码载荷
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	// 解析声明
	var claims JWTClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	// 验证过期时间
	if time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

// hmacSha256 计算HMAC SHA-256签名并返回Base64编码的字符串
func hmacSha256(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
