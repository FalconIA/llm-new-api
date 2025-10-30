# 前端签名配置功能实现总结

## 实现概述
为 New API 项目的自定义渠道（Type 8）添加了签名算法配置功能的前端实现。采用动态表单配置方案，支持灵活扩展多种签名算法。

## 文件清单

### 新增文件
```
web/src/components/table/channels/modals/signature/
├── README.md                           # 组件文档
├── signatureRegistry.js                # 签名算法注册表
├── SignatureConfigPanel.jsx            # 签名配置面板主组件
└── algorithms/
    └── CmcShanghaiConfig.jsx           # 上海移动签名算法配置组件
```

### 修改文件
1. `web/src/components/table/channels/modals/EditChannelModal.jsx`
   - 导入 SignatureConfigPanel 组件
   - 在 type===8 的渠道配置区域添加签名配置面板
   - 在 originInputs 中添加签名字段默认值
   - 在 loadChannel() 中添加签名字段解析逻辑
   - 在 submit() 中添加签名字段保存和清理逻辑

2. `web/src/i18n/locales/zh.json`
   - 添加签名配置相关的中文翻译

3. `web/src/i18n/locales/en.json`
   - 添加签名配置相关的英文翻译

## 核心特性

### 1. 模块化设计
- **注册表模式**：`signatureRegistry.js` 集中管理所有签名算法
- **组件解耦**：每个签名算法对应独立的配置组件
- **动态渲染**：根据选择的算法动态加载对应的配置组件

### 2. 用户体验
- **默认值**：签名算法默认为"无"，不影响现有功能
- **表单验证**：对必填字段进行前端验证
- **提示信息**：使用 Banner 组件显示算法说明和使用提示
- **样式一致**：使用 Semi Design UI 组件，保持界面风格统一

### 3. 数据处理
- **序列化/反序列化**：签名配置存储在 `settings` 字段（JSON 格式）
- **字段清理**：提交时自动清理临时字段，避免污染数据库
- **向后兼容**：老版本渠道没有签名配置时使用默认值

## 数据流图

```
┌─────────────────┐
│ API 响应数据    │
│ (settings 字段) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ loadChannel()   │
│ 解析 JSON       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ inputs 状态     │
│ (分离字段)      │
└────────┬────────┘
         │
         ▼
┌──────────────────────┐
│ SignatureConfigPanel │
│ 用户交互界面         │
└────────┬─────────────┘
         │
         ▼
┌──────────────────────────────┐
│ handleChannelOtherSettings   │
│ Change()                     │
│ 更新状态和 settings JSON     │
└────────┬─────────────────────┘
         │
         ▼
┌─────────────────┐
│ submit()        │
│ 构建 settings   │
│ 清理临时字段    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ API 请求        │
│ (settings 字段) │
└─────────────────┘
```

## 配置字段映射

| 前端字段名              | 后端字段名              | 数据类型 | 必填 | 说明                    |
|------------------------|------------------------|---------|------|------------------------|
| signature_type         | signature_type         | string  | 否   | 签名算法名称            |
| signature_app_id       | signature_app_id       | string  | 条件 | 签名应用ID（cmc_sh必填）|
| signature_time_offset  | signature_time_offset  | int     | 否   | 时间偏移量（秒）        |

## 支持的签名算法

### 1. 无（none）
- **值**：空字符串 `""`
- **说明**：不使用签名功能
- **配置**：无需额外配置

### 2. 上海移动签名算法（cmc_sh）
- **值**：`"cmc_sh"`
- **说明**：使用 AES/ECB/PKCS5Padding 加密算法生成签名
- **必填字段**：
  - `signature_app_id`：签名应用ID
- **可选字段**：
  - `signature_time_offset`：时间偏移量，范围 -3600 到 3600 秒

## 国际化翻译键

| 翻译键                                      | 中文                                    | 英文                                    |
|--------------------------------------------|----------------------------------------|----------------------------------------|
| signature.title                            | 签名配置                                | Signature Configuration                |
| signature.fields.signature_type            | 签名算法                                | Signature Algorithm                    |
| signature.fields.signature_app_id          | 签名应用ID                              | Signature App ID                       |
| signature.fields.signature_time_offset     | 时间偏移量                              | Time Offset                            |
| signature.algorithms.none                  | 无                                      | None                                   |
| signature.algorithms.cmc_sh                | 上海移动签名算法                         | CMC Shanghai signature algorithm        |
| signature.validation.app_id_required       | 签名应用ID为必填项                      | Signature App ID is required           |
| signature.no_algorithm_selected            | 未选择签名算法，默认不对请求进行签名处理 | No signature algorithm selected        |

## 扩展指南

### 添加新签名算法的步骤

#### 1. 创建配置组件
```jsx
// algorithms/NewAlgorithmConfig.jsx
import React from 'react';
import { Form } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';

const NewAlgorithmConfig = ({ channelOtherSettings, handleChannelOtherSettingsChange }) => {
  const { t } = useTranslation();
  
  return (
    <div style={{ marginTop: '16px' }}>
      <Form.Input
        label={t('signature.fields.new_field')}
        field="new_field"
        value={channelOtherSettings.new_field || ''}
        onChange={(value) => handleChannelOtherSettingsChange('new_field', value)}
      />
    </div>
  );
};

export default NewAlgorithmConfig;
```

#### 2. 在注册表中注册
```javascript
// signatureRegistry.js
import NewAlgorithmConfig from './algorithms/NewAlgorithmConfig';

export const SIGNATURE_ALGORITHMS = {
  // ... 现有算法
  new_algorithm: {
    value: 'new_algorithm',
    label: 'signature.algorithms.new_algorithm',
    description: 'signature.algorithms.new_algorithm_desc',
    component: NewAlgorithmConfig,
    requiredFields: ['new_field'],
  },
};
```

#### 3. 添加翻译
在 `zh.json` 和 `en.json` 中添加相应的翻译键。

#### 4. 实现后端逻辑
在后端实现对应的签名算法（参考后端文档）。

## 技术要点

### 1. React Hooks 使用
- `useState`：管理组件内部状态
- `useTranslation`：国际化翻译

### 2. Semi Design UI 组件
- `Card`：配置面板容器
- `Form.Select`：算法选择下拉框
- `Form.Input`：文本输入框
- `Form.InputNumber`：数字输入框
- `Banner`：提示信息横幅

### 3. 条件渲染
```jsx
{SignatureConfigComponent && (
  <SignatureConfigComponent
    channelOtherSettings={channelOtherSettings}
    handleChannelOtherSettingsChange={handleChannelOtherSettingsChange}
  />
)}
```

### 4. 动态组件加载
```javascript
const SignatureConfigComponent = getSignatureComponent(currentAlgorithm);
```

## 测试要点

### 功能测试
1. ✅ 默认显示"无"签名算法
2. ✅ 切换签名算法时动态显示对应配置项
3. ✅ 上海移动签名算法必填字段验证
4. ✅ 保存后重新加载，配置正确显示
5. ✅ 切换算法时清空之前的配置

### 兼容性测试
1. ✅ 老版本渠道（无签名配置）正常显示
2. ✅ 非自定义渠道（type !== 8）不显示签名配置
3. ✅ 中英文翻译正确切换

### 数据完整性测试
1. ✅ settings JSON 格式正确
2. ✅ 临时字段正确清理
3. ✅ 后端接收的数据结构正确

## 与后端集成

### API 数据格式

#### 保存时（前端 → 后端）
```json
{
  "settings": "{\"signature_type\":\"cmc_sh\",\"signature_app_id\":\"your-app-id\",\"signature_time_offset\":0}"
}
```

#### 加载时（后端 → 前端）
```json
{
  "settings": "{\"signature_type\":\"cmc_sh\",\"signature_app_id\":\"your-app-id\",\"signature_time_offset\":0}"
}
```

### 后端处理流程
1. 接收 `settings` 字段（JSON 字符串）
2. 解析为 `ChannelOtherSettings` 结构体
3. 根据 `signature_type` 选择签名器
4. 使用 `signature_app_id` 和 `signature_time_offset` 生成签名
5. 将签名添加到请求头的 `Authorization` 字段

## 已知限制

1. **算法切换**：切换算法时会清空所有签名配置字段
2. **验证时机**：前端验证仅在提交时触发，实时验证需额外实现
3. **错误提示**：后端签名生成失败的错误提示需要后端返回详细信息

## 未来优化方向

1. **实时验证**：添加字段级别的实时验证
2. **配置预设**：支持保存和加载常用配置模板
3. **测试功能**：在配置界面直接测试签名生成
4. **批量配置**：支持批量创建渠道时配置签名
5. **配置导入导出**：支持 JSON 格式的配置导入导出

## 相关文档

- [后端签名实现文档](../../../../relay/channel/openai/signature/README.md)
- [后端使用指南](../../../../relay/channel/openai/signature/USAGE_GUIDE.md)
- [组件使用文档](./README.md)

## 维护说明

### 代码规范
- 遵循项目现有的代码风格
- 使用 ESLint 进行代码检查
- 组件命名使用 PascalCase
- 函数命名使用 camelCase

### 提交规范
- 前端相关提交使用 `feat(frontend):` 前缀
- Bug 修复使用 `fix(frontend):` 前缀
- 文档更新使用 `docs(frontend):` 前缀

### 版本兼容
- 确保向后兼容老版本渠道数据
- 新增字段使用默认值，避免破坏性变更
- API 变更需同步更新前后端

## 贡献者
- 功能设计：基于用户需求和后端架构
- 前端实现：遵循方案二（动态表单配置）
- 文档编写：包含用户文档和开发者文档

## 版本历史
- v1.0 (2025)：初始实现，支持上海移动签名算法（CMC Shanghai）
