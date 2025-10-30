import React from 'react';
import { Banner, Form } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';

/**
 * 上海移动签名算法配置组件
 * @param {Object} props - 组件属性
 * @param {Object} props.channelOtherSettings - 渠道其他设置对象
 * @param {Function} props.handleChannelOtherSettingsChange - 设置变更处理函数
 * @returns {React.Component}
 */
const CmcShanghaiConfig = ({ channelOtherSettings, handleChannelOtherSettingsChange }) => {
  const { t } = useTranslation();

  return (
    <div style={{ marginTop: '16px' }}>
      <Banner
        type="info"
        description={t('signature.algorithms.cmc_sh_hint')}
        style={{ marginBottom: '16px' }}
      />
      
      <Form.Input
        label={t('signature.fields.signature_app_id')}
        placeholder={t('signature.fields.signature_app_id_placeholder')}
        field="signature_app_id"
        value={channelOtherSettings.signature_app_id || ''}
        onChange={(value) => handleChannelOtherSettingsChange('signature_app_id', value || null)}
        rules={[
          { required: true, message: t('signature.validation.app_id_required') }
        ]}
      />
      
      <Form.InputNumber
        label={t('signature.fields.signature_time_offset')}
        placeholder={t('signature.fields.signature_time_offset_placeholder')}
        field="signature_time_offset"
        value={channelOtherSettings.signature_time_offset ?? 0}
        onChange={(value) => handleChannelOtherSettingsChange('signature_time_offset', value ?? null)}
        min={-3600}
        max={3600}
        step={1}
        style={{ width: '100%' }}
        suffix={t('signature.fields.seconds')}
      />
      
      <div style={{ marginTop: '8px', fontSize: '12px', color: '#999' }}>
        {t('signature.fields.signature_time_offset_hint')}
      </div>
    </div>
  );
};

export default CmcShanghaiConfig;
