# 自定义渠道签名功能使用指南

## 快速开始

### 场景说明
当你的API服务商要求使用签名算法生成动态token时，可以使用本功能。

### 配置步骤

#### 1. 创建自定义渠道
在管理后台创建一个自定义渠道，填写基本信息：
- **渠道类型**: Custom（自定义）
- **名称**: 例如 "CMC上海AI"
- **密钥**: 填写原始的API密钥（例如：`1234567890123456`）
- **Base URL**: 填写上游API地址

#### 2. 配置签名参数
在"其他设置"字段中填写JSON配置：

```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "your-app-id"
}
```

**参数说明**：
- `signature_type`: **必填**，签名算法类型，当前支持 `cmc_sh`
- `signature_app_id`: **可选**，部分算法必填，签名应用ID，会包含在签名数据中
- `signature_time_offset`: **可选**，时间偏移量（秒），不填写则默认为0（无偏移）

#### 3. 保存并测试
保存渠道配置后，系统会自动：
1. 使用原始密钥生成签名token
2. 将签名后的token填充到 `Authorization: Bearer {token}` 头中
3. 发起上游请求

## 配置示例

### 示例1：基础配置（必填字段）
```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "myapp-12345"
}
```
生成的payload示例：
```json
{"time": 1730188800000, "app_id": "myapp-12345"}
```

### 示例2：带时间偏移（高级用法）
```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "myapp-12345",
  "signature_time_offset": 300
}
```
说明：时间戳会向前偏移300秒（5分钟），用于服务器时间不同步的场景

**注意**：大多数情况下不需要配置 `signature_time_offset`，只在时间同步存在问题时才使用。

### ❌ 示例3：错误配置
```json
{
  "signature_type": "cmc_sh"
}
```
**错误**：缺少 `signature_app_id` 字段，会报错：`signature_app_id is required for cmc_sh algorithm`

## 工作原理

### 上海移动签名算法（CMC Shanghai）流程

1. **准备数据**
   ```
   当前时间戳 = 当前时间(毫秒) + 偏移量(秒) * 1000
   ```

2. **构建JSON**
   ```json
   {
     "time": 当前时间戳   // 毫秒
   }
   ```

3. **AES加密**
   - 算法：AES/ECB/PKCS5Padding
   - 密钥：原始API密钥
   - 输入：JSON字符串

4. **Base64编码**
   ```
   加密Token = Base64(加密后的字节)
   ```

4. **拼接应用ID**
   ```
   最终Token = {拼接应用ID}.{加密Token}
   ```

5. **使用token**
   ```
   Authorization: {最终Token}
   ```

## 调试技巧

### 1. 查看生成的token
启用调试模式后，日志中会显示：
```
Custom channel signature enabled, type: cmc_sh, app_id: myapp-12345
```

### 2. 验证签名正确性
可以使用测试工具解密验证：
1. 分离token的两部分：`{app_id}.{encrypted_token}`
2. 对加密部分（`encrypted_token`）进行Base64解码
3. 使用原始密钥进行AES解密
4. 验证解密后的JSON格式：应该只包含 `{"time": 时间戳}` 字段
5. 验证app_id在token前缀中，而不在加密的JSON payload中

### 3. 常见问题排查

**问题：签名失败 - 签名器未找到**
```
错误：failed to generate signature: signer 'cmc_sh' not found
```
解决：检查 `signature_type` 拼写是否正确

**问题：缺少必填字段**
```
错误：signature_app_id is required for cmc_sh algorithm
```
解决：部分算法要求必须配置 `signature_app_id`，请添加该字段

**问题：密钥错误**
```
错误：failed to create AES cipher: invalid key size
```
解决：确保API密钥长度为16、24或32字节

**问题：时间不同步**
解决：使用 `signature_time_offset` 调整时间偏移

## 性能说明

- **缓存机制**: 生成的token会缓存4分钟
- **并发安全**: 支持多并发请求
- **自动刷新**: 缓存即将过期时自动后台刷新

## 完整配置示例
## 完整配置示例

### 推荐配置（最常用）
```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "production-app-001"
}
```

### 完整配置（包含可选参数）
```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "production-app-001",
  "signature_time_offset": 0
}
```
**说明**：`signature_time_offset` 为可选参数，不填写时默认为0。

渠道其他信息：
- **密钥**: `1234567890123456` (16字节AES密钥)
- **Base URL**: `https://api.example.com/v1`
- **模型映射**: 根据需要配置

## 注意事项

1. ⚠️ **部分算法要求签名应用ID必填**: 部分签名算法要求必须配置 `signature_app_id`，否则会报错
2. ⚠️ **密钥长度**: 必须是16、24或32字节（对应AES-128/192/256）
3. ⚠️ **时间同步**: 大多数情况下服务器时间会自动同步，无需配置 `signature_time_offset`
4. ⚠️ **缓存时间**: token缓存4分钟，请确保上游API支持
5. ✅ **请求头覆盖**: 签名后可以正常使用请求头覆盖功能
6. ✅ **自动重试**: 缓存失效会自动重新生成token
7. ✅ **时间偏移可选**: `signature_time_offset` 是可选的，可根据实际需求配置

## 相关文档

- [完整技术文档](./README.md)
- [单元测试](./cmc_sh_signer_test.go)
- [签名器实现](./cmc_sh_signer.go)
