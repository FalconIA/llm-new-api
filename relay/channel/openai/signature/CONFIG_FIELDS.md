# 签名配置字段说明

## 配置字段详解

在渠道的"其他设置"中，签名相关的配置字段如下：

| 字段名 | 类型 | 必填 | 默认值 | 说明 |
|-------|------|------|--------|------|
| `signature_type` | string | **必填** | - | 签名算法类型。留空则不启用签名 |
| `signature_app_id` | string | **条件必填** | `""` | 签名应用ID，会包含在签名数据中，部分算法必填 |
| `signature_time_offset` | int64 | 选填 | `0` | 时间偏移量（秒），用于调整时间戳 |

## 配置行为说明

### 必填字段

- **signature_type**: 这是唯一全局必填的字段。当此字段为空或不存在时，签名功能不会启用。

### 可选字段

- **signature_app_id**: 
  - 部分算法必填，不提供会报错
  - 具体要求取决于所使用的签名算法

- **signature_time_offset**:
  - 可以不配置，不配置时默认值为 `0`（无时间偏移）
  - 配置为 `0` 时，等同于不配置
  - 由于使用了 `omitempty` 标签，JSON序列化时值为0会自动省略

## 配置示例

### 最简配置（上海移动签名算法）
```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "myapp-12345"
}
```
**效果**: 
- 启用签名
- 包含应用ID（必填）
- 无时间偏移

**生成的payload**:
```json
{"time": 1730188800000, "app_id": "myapp-12345"}
```

---

### ❌ 错误配置（缺少必填字段）
```json
{
  "signature_type": "cmc_sh"
}
```
**错误**: `signature_app_id is required for cmc_sh algorithm`

---

### 带时间偏移配置
```json
{
  "signature_type": "cmc_sh",
  "signature_app_id": "myapp-12345",
  "signature_time_offset": 300
}
```
**效果**:
- 启用签名
- 包含应用ID
- 时间向前偏移300秒

**生成的payload**:
```json
{"time": 1730189100000, "app_id": "myapp-12345"}
```

## 零值处理规则

### Go语言零值
- `string` 类型的零值: `""`
- `int64` 类型的零值: `0`

### JSON omitempty 标签
配置字段都使用了 `omitempty` 标签：
```go
SignatureType       string `json:"signature_type,omitempty"`
SignatureAppId      string `json:"signature_app_id,omitempty"`
SignatureTimeOffset int64  `json:"signature_time_offset,omitempty"`
```

**作用**:
- 当字段值为零值时，JSON序列化时会自动省略该字段
- 反序列化时，如果JSON中不存在该字段，会自动设置为零值

### 代码中的零值判断

```go
// 时间偏移处理
if config.TimeOffset != 0 {
    timeMillis += config.TimeOffset * 1000
}

// 应用ID处理
if config.AppID != "" {
    payloadMap["app_id"] = config.AppID
}
```

## 常见问题

### Q: 不配置 signature_app_id 会怎样？
**A**: 取决于所使用的签名算法。部分算法要求必填，会报错：`signature_app_id is required for cmc_sh algorithm`。

### Q: 不配置 signature_time_offset 会怎样？
**A**: 系统会使用默认值 `0`，即不进行时间偏移，使用当前系统时间。

### Q: signature_time_offset 设置为 0 和不设置有区别吗？
**A**: 没有区别，两者效果完全相同。

### Q: 哪些字段是真正必须配置的？
**A**: 
- `signature_type` - 全局必填
- `signature_app_id` - 取决于算法，部分算法必填（如 cmc_sh）
- `signature_time_offset` - 完全可选

## 推荐配置

根据使用场景选择合适的配置：

| 场景 | 推荐配置 |
|------|---------|
| 标准配置 | 配置 `signature_type` |
| 部分算法，需要标识应用 | 配置 `signature_type` + `signature_app_id` |
| 时间同步有问题 | 配置 `signature_type` + `signature_time_offset` |
| 完整配置 | 配置所有字段 |

**建议**: 查看所使用签名算法的具体要求，按需配置必要字段。
