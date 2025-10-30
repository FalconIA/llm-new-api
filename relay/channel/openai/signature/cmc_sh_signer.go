package signature

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// CmcShanghaiSigner 上海移动签名算法
// 算法：AES/ECB/PKCS5Padding
// 生成规则：
// 1. 构建JSON对象 {"time": 当前时间戳(毫秒) + 偏移}
// 2. 使用原始apiKey作为密钥进行AES加密
// 3. Base64编码作为新的apiKey
type CmcShanghaiSigner struct{}

func (s *CmcShanghaiSigner) GetName() string {
	return "cmc_sh"
}

func (s *CmcShanghaiSigner) GetCacheDuration() time.Duration {
	// Token有效期5分钟，缓存4分钟后刷新
	return 4 * time.Minute
}

func (s *CmcShanghaiSigner) SupportCache() bool {
	return true
}

func (s *CmcShanghaiSigner) ValidateConfig(apiKey string, config SignatureConfig) error {
	//  算法要求 apiKey 长度：必须是16、24或32字节（对应AES-128/192/256）
	keyLen := len(apiKey)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return fmt.Errorf("invalid API key length for %s: got %d bytes, expected 16, 24, or 32 bytes (AES-128/192/256)", s.GetName(), keyLen)
	}

	// 算法要求 signature_app_id 必填
	if config.AppID == "" {
		return fmt.Errorf("signature_app_id is required for %s signature algorithm", s.GetName())
	}

	return nil
}

func (s *CmcShanghaiSigner) Sign(apiKey string, config SignatureConfig) (string, error) {
	// 算法要求 signature_app_id 必填
	if config.AppID == "" {
		return "", fmt.Errorf("signature_app_id is required for cmc_sh algorithm")
	}

	// 步骤1: 计算带偏移的时间戳（毫秒）
	timeMillis := time.Now().UnixMilli()
	if config.TimeOffset != 0 {
		timeMillis += config.TimeOffset * 1000 // 秒转毫秒
	}

	// 步骤2: 构建payload（只包含time字段）
	payloadMap := map[string]interface{}{
		"time": timeMillis,
	}

	// 步骤3: JSON序列化
	payloadJSON, err := json.Marshal(payloadMap)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 步骤4: 使用原始apiKey作为密钥进行AES加密
	encryptedBytes, err := s.aesEncryptECB(string(payloadJSON), apiKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %w", err)
	}

	// 步骤5: Base64编码
	encryptedToken := base64.StdEncoding.EncodeToString(encryptedBytes)

	// 步骤6: 拼接应用ID前缀
	token := config.AppID + "." + encryptedToken

	return token, nil
}

// aesEncryptECB AES/ECB/PKCS5Padding加密
func (s *CmcShanghaiSigner) aesEncryptECB(plaintext string, key string) ([]byte, error) {
	keyBytes := []byte(key)

	// 创建AES加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// 对明文进行PKCS5填充
	plainBytes := []byte(plaintext)
	paddedPlaintext := s.pkcs5Padding(plainBytes, aes.BlockSize)

	// 使用ECB模式加密（逐块加密）
	ciphertext := make([]byte, len(paddedPlaintext))
	for i := 0; i < len(paddedPlaintext); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], paddedPlaintext[i:i+aes.BlockSize])
	}

	return ciphertext, nil
}

// pkcs5Padding PKCS5填充（与PKCS7填充逻辑相同，因为AES块大小为16）
func (s *CmcShanghaiSigner) pkcs5Padding(plaintext []byte, blockSize int) []byte {
	// 计算需要填充的字节数
	paddingLen := blockSize - (len(plaintext) % blockSize)

	// 创建填充字节数组，填充值等于填充长度
	paddingBytes := make([]byte, paddingLen)
	for i := 0; i < paddingLen; i++ {
		paddingBytes[i] = byte(paddingLen)
	}

	return append(plaintext, paddingBytes...)
}

// aesDecryptECB AES/ECB/PKCS5Padding解密（用于测试验证）
func (s *CmcShanghaiSigner) aesDecryptECB(ciphertext []byte, key string) (string, error) {
	keyBytes := []byte(key)

	// 创建AES解密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// 使用ECB模式解密（逐块解密）
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(plaintext[i:i+aes.BlockSize], ciphertext[i:i+aes.BlockSize])
	}

	// 移除PKCS5填充
	plaintext, err = s.pkcs5Unpadding(plaintext)
	if err != nil {
		return "", fmt.Errorf("failed to remove padding: %w", err)
	}

	return string(plaintext), nil
}

// pkcs5Unpadding 移除PKCS5填充
func (s *CmcShanghaiSigner) pkcs5Unpadding(plaintext []byte) ([]byte, error) {
	if len(plaintext) == 0 {
		return nil, fmt.Errorf("plaintext is empty")
	}

	paddingLen := int(plaintext[len(plaintext)-1])

	if paddingLen > len(plaintext) || paddingLen == 0 {
		return nil, fmt.Errorf("invalid padding length: %d", paddingLen)
	}

	return plaintext[:len(plaintext)-paddingLen], nil
}
