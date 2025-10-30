package signature

import (
	"time"
)

// Signer 签名器接口
type Signer interface {
	// Sign 执行签名
	// apiKey: 原始API密钥
	// config: 签名配置参数
	// Returns: 签名后的新密钥
	Sign(apiKey string, config SignatureConfig) (string, error)

	// ValidateConfig 验证配置参数是否完整
	// apiKey: 原始API密钥
	// config: 签名配置参数
	// Returns: 如果配置不完整或无效，返回错误信息
	ValidateConfig(apiKey string, config SignatureConfig) error

	// GetName 返回签名器名称
	GetName() string

	// GetCacheDuration 返回缓存有效期
	GetCacheDuration() time.Duration

	// SupportCache 是否支持缓存
	SupportCache() bool
}

// SignatureConfig 签名配置
type SignatureConfig struct {
	// AppID 签名应用ID
	AppID string `json:"app_id"`

	// TimeOffset 时间偏移（秒）
	TimeOffset int64 `json:"time_offset"`

	// Extra 额外参数
	Extra map[string]interface{} `json:"extra,omitempty"`
}

// TokenCache token缓存结构
type TokenCache struct {
	Token      string
	ExpiryTime time.Time
}
