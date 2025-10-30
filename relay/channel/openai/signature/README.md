# 自定义渠道签名算法说明

## 概述

本功能为自定义渠道（Custom Channel）提供了灵活的签名算法支持，通过策略模式设计，可以方便地扩展不同的签名算法。

## 架构设计

```
relay/channel/openai/signature/
├── interface.go          # 签名器接口定义
├── factory.go            # 签名器工厂（注册、获取、缓存管理）
├── cmc_sh_signer.go      # 上海移动签名算法实现
└── cmc_sh_signer_test.go # 单元测试
```

### 核心接口

```go
type Signer interface {
    Sign(apiKey string, config SignatureConfig) (string, error)
    GetName() string
    GetCacheDuration() time.Duration
    SupportCache() bool
}
```

## 已实现的签名算法

### 1. 上海移动签名算法（CMC Shanghai）

**算法类型**: AES/ECB/PKCS5Padding

**实现逻辑**:
1. 构建JSON payload: `{"time": 当前时间戳(毫秒) + 偏移}`
2. 使用原始API密钥作为AES密钥进行加密
3. Base64编码加密结果
4. 拼接应用ID前缀：`{app_id}.{encrypted_base64}`
5. 缓存有效期：4分钟

**配置参数**:
- `signature_app_id`: 签名应用ID（必填，会拼接在token前缀）
- `signature_time_offset`: 时间偏移（秒，可选，默认0）

## 使用方法

### 1. 在渠道配置中启用签名

在管理后台的"其他设置"中添加以下JSON配置：

```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "your-app-id"
}
```

**配置说明**:
- `signature_type`: 签名算法类型，目前支持 `cmc_sh`，留空则不启用签名（必填）
- `signature_app_id`: 签名应用ID，会包含在签名数据中（可选，部分算法必填）
- `signature_time_offset`: 时间偏移（秒），用于调整时间戳（可选，不填写则默认为0，即无偏移）

### 2. API密钥要求

- **原始密钥**: 在渠道配置的"密钥"字段填写原始API密钥
- **签名后**: 系统会自动使用原始密钥生成签名token，并替换到请求的Authorization头中
- **优势**: 签名后的token可以正常使用后续的请求头覆盖（header mapping）功能

### 3. 示例配置

#### 基础配置（部分算法必须包含signature_app_id）
```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "myapp-12345"
}
```

#### 带时间偏移
```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "myapp-12345",
  "signature_time_offset": 60
}
```

#### ❌ 错误配置（缺少app_id会报错）
```json
{
  "signature_type": "cmc_sh"
}
```
**错误信息**: `signature_app_id is required for cmc_sh algorithm`

#### 带时间偏移（向前偏移60秒）
```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "myapp-12345",
  "signature_time_offset": 60
}
```

## 工作流程

1. **请求到达**: 用户请求到达自定义渠道
2. **检查配置**: 系统检查 `signature_type` 是否配置
3. **生成签名**: 
   - 读取原始API密钥
   - 根据配置生成签名token
   - 检查缓存，如有效则直接使用
4. **替换密钥**: 将生成的签名token替换原始API密钥
5. **发起请求**: 使用签名后的token发起上游请求

## 缓存机制

- **自动缓存**: 生成的签名token会自动缓存
- **有效期**: 根据签名器定义的缓存时长（cmc_sh为4分钟）
- **提前刷新**: 缓存即将过期前1分钟会后台异步刷新
- **缓存键**: 由 `签名类型:API密钥:AppID:时间偏移` 组成

## 性能优化

1. **缓存优先**: 优先使用缓存的token，避免重复计算
2. **异步刷新**: 即将过期的token会后台异步刷新，不影响当前请求
3. **线程安全**: 使用 `sync.Map` 保证并发安全

## 扩展新的签名算法

### 步骤1: 实现Signer接口

创建新文件 `your_signer.go`:

```go
package signature

import "time"

type YourSigner struct{}

func (s *YourSigner) GetName() string {
    return "your_algorithm"
}

func (s *YourSigner) GetCacheDuration() time.Duration {
    return 5 * time.Minute
}

func (s *YourSigner) SupportCache() bool {
    return true
}

func (s *YourSigner) Sign(apiKey string, config SignatureConfig) (string, error) {
    // 实现你的签名逻辑
    return "signed-token", nil
}
```

### 步骤2: 注册签名器

在 `factory.go` 的 `init()` 函数中注册：

```go
func init() {
    RegisterSigner(&CmcShanghaiSigner{})
    RegisterSigner(&YourSigner{})  // 新增
}
```

### 步骤3: 配置使用

```json
{
  "signature_type": "your_algorithm",
  "signature_app_id": "your-app-id"
}
```

## 测试

### 运行单元测试

```bash
cd relay/channel/openai/signature
go test -v
```

### 运行性能测试

```bash
go test -bench=. -benchmem
```

### 测试输出示例

```
=== RUN   TestCmcShanghaiSigner/Basic_Sign
Generated token: k1jRZxJ4e5F...
Decrypted content: {"time":1730188800000,"app_id":"test-app-001"}
--- PASS: TestCmcShanghaiSigner/Basic_Sign (0.00s)

=== RUN   TestCmcShanghaiSigner/Sign_With_Time_Offset
Time offset verification passed, timestamp: 1730188860000
--- PASS: TestCmcShanghaiSigner/Sign_With_Time_Offset (0.00s)
```

## 调试

启用调试模式后，系统会输出签名相关日志：

```
Custom channel signature enabled, type: cmc_sh, app_id: myapp-12345
```

## 安全建议

1. **密钥长度**: AES密钥建议使用16字节（128位）、24字节（192位）或32字节（256位）
2. **密钥保护**: 原始API密钥应妥善保管，不要泄露
3. **时间同步**: 确保服务器时间准确，避免时间偏移导致的问题
4. **缓存刷新**: 缓存过期时间应小于上游API的token有效期

## 故障排查

### 问题1: 签名失败

**错误信息**: `failed to generate signature: signer 'xxx' not found`

**解决方案**: 检查 `signature_type` 配置是否正确，确保签名器已注册

### 问题2: 密钥长度错误

**错误信息**: `failed to create AES cipher: crypto/aes: invalid key size`

**解决方案**: 确保API密钥长度为16、24或32字节

### 问题3: 时间戳不匹配

**解决方案**: 
- 检查服务器时间是否同步
- 调整 `signature_time_offset` 参数补偿时间差

## 版本历史

- **v1.0.0** (2025-10-29)
  - 初始版本
  - 实现上海移动签名算法（cmc_sh）
  - 支持配置驱动的签名策略
  - 完整的缓存机制
