import React from 'react';
import CmcShanghaiConfig from './algorithms/CmcShanghaiConfig';

/**
 * 签名算法注册表
 * 用于管理所有可用的签名算法及其配置组件
 */
export const SIGNATURE_ALGORITHMS = {
  none: {
    value: '',
    label: 'signature.algorithms.none',
    description: 'signature.algorithms.none_desc',
    component: null, // 无签名算法，不需要额外配置
  },
  cmc_sh: {
    value: 'cmc_sh',
    label: 'signature.algorithms.cmc_sh',
    description: 'signature.algorithms.cmc_sh_desc',
    component: CmcShanghaiConfig,
    requiredFields: ['signature_app_id'], // 必填字段列表
  },
  // 未来可以在此处添加更多签名算法
  // example_algorithm: {
  //   value: 'example_algorithm',
  //   label: 'signature.algorithms.example',
  //   description: 'signature.algorithms.example_desc',
  //   component: ExampleConfig,
  //   requiredFields: ['some_field'],
  // },
};

/**
 * 获取签名算法下拉选项
 * @returns {Array} 签名算法选项数组
 */
export const getSignatureOptions = () => {
  return Object.values(SIGNATURE_ALGORITHMS).map(algo => ({
    value: algo.value,
    label: algo.label,
    description: algo.description,
  }));
};

/**
 * 根据算法名称获取配置组件
 * @param {string} algorithmName - 签名算法名称
 * @returns {React.Component|null} 配置组件或null
 */
export const getSignatureComponent = (algorithmName) => {
  if (!algorithmName || algorithmName === '') {
    return null;
  }
  
  const algorithm = SIGNATURE_ALGORITHMS[algorithmName];
  return algorithm ? algorithm.component : null;
};

/**
 * 获取算法的必填字段列表
 * @param {string} algorithmName - 签名算法名称
 * @returns {Array} 必填字段数组
 */
export const getRequiredFields = (algorithmName) => {
  if (!algorithmName || algorithmName === '') {
    return [];
  }
  
  const algorithm = SIGNATURE_ALGORITHMS[algorithmName];
  return algorithm?.requiredFields || [];
};
