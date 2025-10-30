package signature

import (
	"fmt"
	"sync"
	"time"
)

var (
	signers      = make(map[string]Signer)
	signersMutex sync.RWMutex
	tokenCache   sync.Map
)

func init() {
	// 注册所有签名器
	RegisterSigner(&CmcShanghaiSigner{})
}

// RegisterSigner 注册签名器
func RegisterSigner(signer Signer) {
	signersMutex.Lock()
	defer signersMutex.Unlock()
	signers[signer.GetName()] = signer
}

// GetSigner 获取签名器
func GetSigner(name string) (Signer, error) {
	signersMutex.RLock()
	defer signersMutex.RUnlock()

	signer, exists := signers[name]
	if !exists {
		return nil, fmt.Errorf("signer '%s' not found", name)
	}
	return signer, nil
}

// SignWithCache 使用签名器生成token（带缓存）
func SignWithCache(signerName, apiKey string, config SignatureConfig) (string, error) {
	signer, err := GetSigner(signerName)
	if err != nil {
		return "", err
	}

	// 提前验证配置参数（包括 apiKey）
	if err := signer.ValidateConfig(apiKey, config); err != nil {
		return "", err
	}

	// 如果不支持缓存，直接生成
	if !signer.SupportCache() {
		return signer.Sign(apiKey, config)
	}

	// 构建缓存key
	cacheKey := fmt.Sprintf("%s:%s:%s:%d", signerName, apiKey, config.AppID, config.TimeOffset)

	// 检查缓存
	if cached, ok := tokenCache.Load(cacheKey); ok {
		if tokenData, ok := cached.(TokenCache); ok {
			// 提前1分钟刷新
			if time.Now().Add(1 * time.Minute).Before(tokenData.ExpiryTime) {
				return tokenData.Token, nil
			}
			// 后台异步刷新即将过期的token
			go func() {
				_, _ = refreshToken(signer, cacheKey, apiKey, config)
			}()
			return tokenData.Token, nil
		}
	}

	return refreshToken(signer, cacheKey, apiKey, config)
}

// refreshToken 刷新token
func refreshToken(signer Signer, cacheKey, apiKey string, config SignatureConfig) (string, error) {
	token, err := signer.Sign(apiKey, config)
	if err != nil {
		return "", err
	}

	// 缓存
	expiryTime := time.Now().Add(signer.GetCacheDuration())
	tokenCache.Store(cacheKey, TokenCache{
		Token:      token,
		ExpiryTime: expiryTime,
	})

	return token, nil
}

// ClearCache 清空指定的token缓存
func ClearCache(signerName, apiKey string, config SignatureConfig) {
	cacheKey := fmt.Sprintf("%s:%s:%s:%d", signerName, apiKey, config.AppID, config.TimeOffset)
	tokenCache.Delete(cacheKey)
}

// ClearAllCache 清空所有缓存
func ClearAllCache() {
	tokenCache.Range(func(key, value interface{}) bool {
		tokenCache.Delete(key)
		return true
	})
}

// ListSigners 列出所有可用的签名器
func ListSigners() []string {
	signersMutex.RLock()
	defer signersMutex.RUnlock()

	names := make([]string, 0, len(signers))
	for name := range signers {
		names = append(names, name)
	}
	return names
}

// ValidateSignatureConfig 验证签名配置是否有效（用于保存时校验）
// signerName: 签名器名称
// apiKey: 原始API密钥
// config: 签名配置参数
// Returns: 如果配置无效，返回错误信息
func ValidateSignatureConfig(signerName, apiKey string, config SignatureConfig) error {
	signer, err := GetSigner(signerName)
	if err != nil {
		return err
	}

	return signer.ValidateConfig(apiKey, config)
}
