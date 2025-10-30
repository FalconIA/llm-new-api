# 签名配置组件

## 概述
此目录包含用于自定义渠道（Type 8）的签名算法配置组件。签名功能允许渠道在发送请求时自动生成并添加签名令牌到请求头。

## 目录结构
```
signature/
├── README.md                      # 本文档
├── signatureRegistry.js           # 签名算法注册表
├── SignatureConfigPanel.jsx       # 签名配置面板主组件
└── algorithms/                    # 各签名算法的配置组件
    └── CmcShanghaiConfig.jsx      # CMC Shanghai 算法配置组件
```

## 核心文件说明

### 1. signatureRegistry.js
**功能**：管理所有可用的签名算法及其配置组件。

**主要导出**：
- `SIGNATURE_ALGORITHMS`：算法注册表对象
- `getSignatureOptions()`：获取下拉选项列表
- `getSignatureComponent(algorithmName)`：根据算法名称获取配置组件
- `getRequiredFields(algorithmName)`：获取算法的必填字段列表

**添加新算法**：
```javascript
export const SIGNATURE_ALGORITHMS = {
  // ... 现有算法
  
  // 添加新算法
  your_algorithm: {
    value: 'your_algorithm',
    label: 'signature.algorithms.your_algorithm',
    description: 'signature.algorithms.your_algorithm_desc',
    component: YourAlgorithmConfig,  // 导入你的配置组件
    requiredFields: ['field1', 'field2'],  // 必填字段列表
  },
};
```

### 2. SignatureConfigPanel.jsx
**功能**：签名配置的主面板组件，用于在渠道编辑模态框中显示。

**Props**：
- `channelOtherSettings` (Object)：渠道其他设置对象
- `handleChannelOtherSettingsChange` (Function)：设置变更处理函数

**使用示例**：
```jsx
<SignatureConfigPanel
  channelOtherSettings={inputs.other || {}}
  handleChannelOtherSettingsChange={handleChannelOtherSettingsChange}
/>
```

### 3. algorithms/CmcShanghaiConfig.jsx
**功能**：上海移动签名算法的配置组件。

**配置字段**：
- `signature_app_id`：签名应用ID（必填）
- `signature_time_offset`：时间偏移量（可选，默认0）

## 数据流

### 1. 加载渠道数据
```
API 响应 → loadChannel() → 解析 settings 字段 → 设置到 inputs 状态
```

在 `EditChannelModal.jsx` 的 `loadChannel()` 函数中：
```javascript
const parsedSettings = JSON.parse(data.settings);
data.signature_type = parsedSettings.signature_type || '';
data.signature_app_id = parsedSettings.signature_app_id || '';
data.signature_time_offset = parsedSettings.signature_time_offset ?? 0;
```

### 2. 用户修改配置
```
用户输入 → SignatureConfigPanel → handleChannelOtherSettingsChange() → 更新 inputs 和 settings
```

`handleChannelOtherSettingsChange()` 函数会：
1. 更新组件内部状态
2. 同步到表单字段
3. 更新 `settings` JSON 字符串

### 3. 提交保存
```
submit() → 收集表单数据 → 构建 settings JSON → 清理临时字段 → API 请求
```

在 `submit()` 函数中：
```javascript
// type === 8 时保存签名配置
if (localInputs.type === 8) {
  if (localInputs.signature_type) {
    settings.signature_type = localInputs.signature_type;
  }
  // ... 其他字段
}

// 清理临时字段
delete localInputs.signature_type;
delete localInputs.signature_app_id;
delete localInputs.signature_time_offset;
```

## 国际化

签名相关的翻译键在 `web/src/i18n/locales/` 目录下的 `zh.json` 和 `en.json` 中：

```json
{
  "signature.title": "签名配置",
  "signature.fields.signature_type": "签名算法",
  "signature.algorithms.none": "无",
  "signature.algorithms.cmc_sh": "上海移动签名算法",
  // ... 更多键
}
```

添加新算法时，请确保在两个语言文件中都添加相应的翻译。

## 添加新签名算法

### 步骤 1：创建配置组件
在 `algorithms/` 目录下创建新文件，例如 `YourAlgorithmConfig.jsx`：

```jsx
import React from 'react';
import { Banner, Form } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';

const YourAlgorithmConfig = ({ channelOtherSettings, handleChannelOtherSettingsChange }) => {
  const { t } = useTranslation();

  return (
    <div style={{ marginTop: '16px' }}>
      <Banner
        type="info"
        description={t('signature.algorithms.your_algorithm_hint')}
        style={{ marginBottom: '16px' }}
      />
      
      <Form.Input
        label={t('signature.fields.your_field')}
        field="your_field"
        value={channelOtherSettings.your_field || ''}
        onChange={(value) => handleChannelOtherSettingsChange('your_field', value)}
        rules={[{ required: true, message: t('signature.validation.your_field_required') }]}
      />
    </div>
  );
};

export default YourAlgorithmConfig;
```

### 步骤 2：在注册表中注册
编辑 `signatureRegistry.js`：

```javascript
import YourAlgorithmConfig from './algorithms/YourAlgorithmConfig';

export const SIGNATURE_ALGORITHMS = {
  // ... 现有算法
  your_algorithm: {
    value: 'your_algorithm',
    label: 'signature.algorithms.your_algorithm',
    description: 'signature.algorithms.your_algorithm_desc',
    component: YourAlgorithmConfig,
    requiredFields: ['your_field'],
  },
};
```

### 步骤 3：添加翻译
在 `zh.json` 和 `en.json` 中添加：

```json
{
  "signature.algorithms.your_algorithm": "您的算法",
  "signature.algorithms.your_algorithm_desc": "您的算法描述",
  "signature.algorithms.your_algorithm_hint": "使用提示",
  "signature.fields.your_field": "字段名称",
  "signature.validation.your_field_required": "该字段为必填项"
}
```

### 步骤 4：实现后端逻辑
在 `relay/channel/openai/signature/` 目录下创建对应的签名实现：

```go
// your_algorithm_signer.go
package signature

type YourAlgorithmSigner struct{}

func (s *YourAlgorithmSigner) Sign(apiKey string, config SignatureConfig) (string, error) {
    // 实现签名逻辑
}

func (s *YourAlgorithmSigner) GetName() string {
    return "your_algorithm"
}

// 在 factory.go 中注册
func init() {
    RegisterSigner(&YourAlgorithmSigner{})
}
```

## 样式一致性

所有配置组件应遵循以下样式规范：

1. **间距**：使用 `marginTop: '16px'` 作为顶部间距
2. **Banner**：使用 Semi Design 的 Banner 组件显示提示信息
3. **表单字段**：使用 Semi Design 的 Form 组件（Form.Input、Form.InputNumber 等）
4. **验证规则**：为必填字段添加 `rules` 属性
5. **国际化**：所有文本使用 `t()` 函数进行翻译

## 调试技巧

### 1. 查看表单值
在浏览器控制台中：
```javascript
// 获取当前表单所有值
formApiRef.current?.getValues()

// 获取特定字段的值
formApiRef.current?.getValue('signature_type')
```

### 2. 查看 settings JSON
在提交前在 `submit()` 函数中添加断点或日志：
```javascript
console.log('Settings JSON:', localInputs.settings);
```

### 3. 检查后端接收的数据
查看后端日志或使用浏览器开发者工具的 Network 标签。

## 常见问题

**Q: 为什么我的配置没有保存？**
A: 确保字段名称在以下位置保持一致：
- `originInputs` 的默认值定义
- `loadChannel()` 的解析逻辑
- `submit()` 的保存逻辑
- 后端 DTO 的字段定义

**Q: 如何测试签名功能？**
A: 
1. 创建或编辑自定义渠道（Type 8）
2. 选择签名算法并配置必要参数
3. 保存渠道
4. 发送测试请求，检查请求头中的 Authorization 字段

**Q: 签名算法默认值是什么？**
A: 默认值为"无"（空字符串），表示不使用签名功能。

## 相关文档

- [后端签名实现文档](../../../../../relay/channel/openai/signature/README.md)
- [后端使用指南](../../../../../relay/channel/openai/signature/USAGE_GUIDE.md)
- [配置字段参考](../../../../../relay/channel/openai/signature/CONFIG_FIELDS.md)
