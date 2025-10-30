package signature

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"
)

func TestCmcShanghaiSigner(t *testing.T) {
	signer := &CmcShanghaiSigner{}

	// 测试用例1: ValidateConfig - 合法配置（16字节密钥）
	t.Run("ValidateConfig With Valid Config and 16-byte Key", func(t *testing.T) {
		apiKey := "1234567890123456" // 16字节的AES密钥
		config := SignatureConfig{
			AppID:      "test-app-001",
			TimeOffset: 0,
		}

		err := signer.ValidateConfig(apiKey, config)
		if err != nil {
			t.Fatalf("Expected no error for valid config, but got: %v", err)
		}

		t.Log("Valid config with 16-byte key passed validation")
	})

	// 测试用例2: ValidateConfig - 合法配置（24字节密钥）
	t.Run("ValidateConfig With Valid Config and 24-byte Key", func(t *testing.T) {
		apiKey := "123456789012345678901234" // 24字节的AES密钥
		config := SignatureConfig{
			AppID:      "test-app-001",
			TimeOffset: 0,
		}

		err := signer.ValidateConfig(apiKey, config)
		if err != nil {
			t.Fatalf("Expected no error for valid config, but got: %v", err)
		}

		t.Log("Valid config with 24-byte key passed validation")
	})

	// 测试用例3: ValidateConfig - 合法配置（32字节密钥）
	t.Run("ValidateConfig With Valid Config and 32-byte Key", func(t *testing.T) {
		apiKey := "12345678901234567890123456789012" // 32字节的AES密钥
		config := SignatureConfig{
			AppID:      "test-app-001",
			TimeOffset: 0,
		}

		err := signer.ValidateConfig(apiKey, config)
		if err != nil {
			t.Fatalf("Expected no error for valid config, but got: %v", err)
		}

		t.Log("Valid config with 32-byte key passed validation")
	})

	// 测试用例4: ValidateConfig - 密钥长度不正确
	t.Run("ValidateConfig With Invalid Key Length", func(t *testing.T) {
		apiKey := "short" // 5字节，不符合要求
		config := SignatureConfig{
			AppID:      "test-app-001",
			TimeOffset: 0,
		}

		err := signer.ValidateConfig(apiKey, config)
		if err == nil {
			t.Fatal("Expected error for invalid key length, but got nil")
		}

		t.Logf("Correctly rejected invalid key length with error: %v", err)
	})

	// 测试用例5: ValidateConfig - 缺少 AppID
	t.Run("ValidateConfig Without AppID", func(t *testing.T) {
		apiKey := "1234567890123456" // 16字节的AES密钥
		config := SignatureConfig{
			AppID:      "",
			TimeOffset: 0,
		}

		err := signer.ValidateConfig(apiKey, config)
		if err == nil {
			t.Fatal("Expected error for missing AppID, but got nil")
		}

		expectedError := "signature_app_id is required for cmc_sh signature algorithm"
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s', got '%s'", expectedError, err.Error())
		}

		t.Logf("Correctly rejected missing AppID with error: %v", err)
	})

	// 测试用例6: 基本签名
	t.Run("Basic Sign", func(t *testing.T) {
		apiKey := "1234567890123456" // 16字节的AES密钥
		config := SignatureConfig{
			AppID:      "test-app-001",
			TimeOffset: 0,
		}

		token, err := signer.Sign(apiKey, config)
		if err != nil {
			t.Fatalf("Sign failed: %v", err)
		}

		if token == "" {
			t.Fatal("Token is empty")
		}

		t.Logf("Generated token: %s", token)

		// 验证token格式：应该是 {app_id}.{encrypted_base64}
		if len(token) < len(config.AppID)+2 {
			t.Fatal("Token format incorrect: too short")
		}

		// 提取app_id前缀和加密部分
		expectedPrefix := config.AppID + "."
		if token[:len(expectedPrefix)] != expectedPrefix {
			t.Fatalf("Token prefix mismatch: expected '%s', got '%s'", expectedPrefix, token[:len(expectedPrefix)])
		}

		// 提取加密的base64部分
		encryptedToken := token[len(expectedPrefix):]

		// 验证加密部分可以解密
		decoded, err := base64.StdEncoding.DecodeString(encryptedToken)
		if err != nil {
			t.Fatalf("Failed to decode base64: %v", err)
		}

		decrypted, err := signer.aesDecryptECB(decoded, apiKey)
		if err != nil {
			t.Fatalf("Failed to decrypt: %v", err)
		}

		t.Logf("Decrypted content: %s", decrypted)

		// 验证JSON格式：只包含time字段
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(decrypted), &payload); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		if _, ok := payload["time"]; !ok {
			t.Fatal("Payload missing 'time' field")
		}

		t.Logf("Token format verified: %s.{encrypted_data}", config.AppID)
	})

	// 测试用例7: 带时间偏移
	t.Run("Sign With Time Offset", func(t *testing.T) {
		apiKey := "1234567890123456"
		config := SignatureConfig{
			AppID:      "test-app-002",
			TimeOffset: 60, // 60秒偏移
		}

		token, err := signer.Sign(apiKey, config)
		if err != nil {
			t.Fatalf("Sign failed: %v", err)
		}

		// 提取加密部分（跳过app_id前缀）
		expectedPrefix := config.AppID + "."
		encryptedToken := token[len(expectedPrefix):]

		// 解密验证
		decoded, _ := base64.StdEncoding.DecodeString(encryptedToken)
		decrypted, _ := signer.aesDecryptECB(decoded, apiKey)

		var payload map[string]interface{}
		json.Unmarshal([]byte(decrypted), &payload)

		timestamp := int64(payload["time"].(float64))
		expectedTime := time.Now().UnixMilli() + config.TimeOffset*1000

		// 允许1秒的误差
		if timestamp < expectedTime-1000 || timestamp > expectedTime+1000 {
			t.Fatalf("Timestamp offset incorrect: got %d, expected around %d", timestamp, expectedTime)
		}

		t.Logf("Time offset verification passed, timestamp: %d", timestamp)
	})

	// 测试用例8: 缺少AppID应该报错（上海移动签名算法要求AppID必填）
	t.Run("Sign Without AppID Should Fail", func(t *testing.T) {
		apiKey := "1234567890123456"
		config := SignatureConfig{
			AppID:      "",
			TimeOffset: 0,
		}

		_, err := signer.Sign(apiKey, config)
		if err == nil {
			t.Fatal("Expected error for missing AppID, but got nil")
		}

		expectedError := "signature_app_id is required for cmc_sh algorithm"
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s', got '%s'", expectedError, err.Error())
		}

		t.Logf("Correctly rejected missing AppID with error: %v", err)
	})

	// 测试用例9: 空配置应该报错（缺少AppID）
	t.Run("Sign With Empty Config Should Fail", func(t *testing.T) {
		apiKey := "1234567890123456"
		config := SignatureConfig{} // 完全空配置，所有字段都是零值

		_, err := signer.Sign(apiKey, config)
		if err == nil {
			t.Fatal("Expected error for missing AppID, but got nil")
		}

		expectedError := "signature_app_id is required for cmc_sh algorithm"
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s', got '%s'", expectedError, err.Error())
		}

		t.Logf("Correctly rejected empty config with error: %v", err)
	})

	// 测试用例10: 错误的密钥长度
	t.Run("Invalid Key Length", func(t *testing.T) {
		apiKey := "short" // 太短的密钥
		config := SignatureConfig{
			AppID:      "test-app-003",
			TimeOffset: 0,
		}

		_, err := signer.Sign(apiKey, config)
		if err == nil {
			t.Fatal("Expected error for invalid key length, but got nil")
		}

		t.Logf("Invalid key length error: %v", err)
	})
}

func TestSignatureFactory(t *testing.T) {
	t.Run("Get CMC Shanghai Signer", func(t *testing.T) {
		signer, err := GetSigner("cmc_sh")
		if err != nil {
			t.Fatalf("Failed to get signer: %v", err)
		}

		if signer.GetName() != "cmc_sh" {
			t.Fatalf("Signer name mismatch: expected 'cmc_sh', got '%s'", signer.GetName())
		}

		t.Logf("Signer retrieved successfully: %s", signer.GetName())
	})

	t.Run("SignWithCache Validates Config", func(t *testing.T) {
		apiKey := "1234567890123456"
		config := SignatureConfig{
			AppID:      "", // 空 AppID，应该失败
			TimeOffset: 0,
		}

		// 应该在签名前验证失败
		_, err := SignWithCache("cmc_sh", apiKey, config)
		if err == nil {
			t.Fatal("Expected validation error for missing AppID, but got nil")
		}

		expectedError := "signature_app_id is required for cmc_sh signature algorithm"
		if err.Error() != expectedError {
			t.Fatalf("Expected error '%s', got '%s'", expectedError, err.Error())
		}

		t.Logf("SignWithCache correctly rejected invalid config: %v", err)
	})

	t.Run("Sign With Cache", func(t *testing.T) {
		apiKey := "1234567890123456"
		config := SignatureConfig{
			AppID:      "test-cache-app",
			TimeOffset: 0,
		}

		// 第一次生成
		token1, err := SignWithCache("cmc_sh", apiKey, config)
		if err != nil {
			t.Fatalf("First sign failed: %v", err)
		}

		// 第二次应该从缓存获取
		token2, err := SignWithCache("cmc_sh", apiKey, config)
		if err != nil {
			t.Fatalf("Second sign failed: %v", err)
		}

		if token1 != token2 {
			t.Fatal("Cached token mismatch")
		}

		t.Logf("Cache test passed, token: %s", token1)

		// 清空缓存
		ClearCache("cmc_sh", apiKey, config)

		// 清空后再次生成（时间戳不同，token应该不同）
		time.Sleep(10 * time.Millisecond)
		token3, err := SignWithCache("cmc_sh", apiKey, config)
		if err != nil {
			t.Fatalf("Third sign failed: %v", err)
		}

		// 注意：由于时间戳不同，token应该不同
		t.Logf("Token after cache clear: %s", token3)
	})

	t.Run("List All Signers", func(t *testing.T) {
		signers := ListSigners()
		if len(signers) == 0 {
			t.Fatal("No signers registered")
		}

		found := false
		for _, name := range signers {
			if name == "cmc_sh" {
				found = true
				break
			}
		}

		if !found {
			t.Fatal("cmc_sh signer not found in list")
		}

		t.Logf("Available signers: %v", signers)
	})

	t.Run("Get Invalid Signer", func(t *testing.T) {
		_, err := GetSigner("non_existent_signer")
		if err == nil {
			t.Fatal("Expected error for non-existent signer")
		}

		t.Logf("Expected error: %v", err)
	})
}

// Benchmark测试
func BenchmarkCmcShanghaiSigner(b *testing.B) {
	signer := &CmcShanghaiSigner{}
	apiKey := "1234567890123456"
	config := SignatureConfig{
		AppID:      "bench-app",
		TimeOffset: 0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = signer.Sign(apiKey, config)
	}
}

func BenchmarkSignWithCache(b *testing.B) {
	apiKey := "1234567890123456"
	config := SignatureConfig{
		AppID:      "bench-cache-app",
		TimeOffset: 0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = SignWithCache("cmc_sh", apiKey, config)
	}
}
