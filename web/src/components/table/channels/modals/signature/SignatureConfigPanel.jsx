import React from 'react';
import { Card, Form, Typography } from '@douyinfe/semi-ui';
import { IconBolt } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import { getSignatureOptions, getSignatureComponent } from './signatureRegistry';

const { Text } = Typography;

/**
 * 签名配置面板组件
 * 用于在渠道编辑模态框中显示签名算法配置
 * @param {Object} props - 组件属性
 * @param {Object} props.channelOtherSettings - 渠道其他设置对象
 * @param {Function} props.handleChannelOtherSettingsChange - 设置变更处理函数
 * @returns {React.Component}
 */
const SignatureConfigPanel = ({ channelOtherSettings, handleChannelOtherSettingsChange }) => {
  const { t } = useTranslation();
  
  const currentAlgorithm = channelOtherSettings.signature_type || '';
  const SignatureConfigComponent = getSignatureComponent(currentAlgorithm);

  const handleAlgorithmChange = (value) => {
    // 切换算法时清空相关配置字段
    // 如果选择"无"，则设置为 null
    const newValue = value === '' ? null : value;
    handleChannelOtherSettingsChange('signature_type', newValue);
    handleChannelOtherSettingsChange('signature_app_id', null);
    handleChannelOtherSettingsChange('signature_time_offset', null);
  };

  return (
    <Card
      title={
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <IconBolt size="large" />
          <Text strong>{t('signature.title')}</Text>
        </div>
      }
      bordered
      style={{ marginTop: '16px' }}
      headerStyle={{ borderBottom: '1px solid var(--semi-color-border)' }}
    >
      <Form.Select
        label={t('signature.fields.signature_type')}
        placeholder={t('signature.fields.signature_type_placeholder')}
        field="signature_type"
        value={currentAlgorithm}
        onChange={handleAlgorithmChange}
        optionList={getSignatureOptions().map(opt => ({
          value: opt.value,
          label: t(opt.label),
          // 可选：添加描述信息
          // description: t(opt.description),
        }))}
        style={{ width: '100%' }}
      />
      
      {SignatureConfigComponent && (
        <SignatureConfigComponent
          channelOtherSettings={channelOtherSettings}
          handleChannelOtherSettingsChange={handleChannelOtherSettingsChange}
        />
      )}
      
      {!currentAlgorithm && (
        <div style={{ marginTop: '12px', padding: '12px', backgroundColor: '#f7f7f7', borderRadius: '4px' }}>
          <Text type="tertiary" size="small">
            {t('signature.no_algorithm_selected')}
          </Text>
        </div>
      )}
    </Card>
  );
};

export default SignatureConfigPanel;
