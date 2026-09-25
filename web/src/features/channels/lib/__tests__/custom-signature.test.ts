/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { describe, expect, test } from 'vitest'

import { channelSchema } from '../../types'
import {
  CHANNEL_FORM_DEFAULT_VALUES,
  channelFormSchema,
  transformChannelToFormDefaults,
  transformFormDataToCreatePayload,
} from '../channel-form'

const customForm = {
  ...CHANNEL_FORM_DEFAULT_VALUES,
  name: 'Signed upstream',
  type: 8,
  base_url: 'https://upstream.example',
  key: '1234567890abcdef',
  models: 'gpt-4o',
  signature_type: 'cmc_sh' as const,
  signature_app_id: 'app_1',
  signature_time_offset: 2,
}

describe('custom channel signature settings', () => {
  test('requires an application ID when signing is enabled', () => {
    expect(channelFormSchema.safeParse(customForm).success).toBe(true)
    expect(
      channelFormSchema.safeParse({ ...customForm, signature_app_id: ' ' })
        .success
    ).toBe(false)
  })

  test('serializes and reloads signed settings', () => {
    const payload = transformFormDataToCreatePayload(customForm).channel
    const saved = JSON.parse(payload.settings || '{}')
    expect(saved).toMatchObject({
      signature_type: 'cmc_sh',
      signature_app_id: 'app_1',
      signature_time_offset: 2,
    })
    const channel = channelSchema.parse({
      ...payload,
      id: 1,
      key: '',
      other: '',
      remark: '',
      status: 1,
      created_time: 0,
      test_time: 0,
      response_time: 0,
      balance_updated_time: 0,
    })
    const reloaded = transformChannelToFormDefaults(channel)
    expect(reloaded.signature_type).toBe('cmc_sh')
    expect(reloaded.signature_app_id).toBe('app_1')
    expect(reloaded.signature_time_offset).toBe(2)
  })

  test('does not retain signing when the channel type changes', () => {
    const payload = transformFormDataToCreatePayload({
      ...customForm,
      type: 1,
    }).channel
    const saved = JSON.parse(payload.settings || '{}')
    expect(saved).not.toHaveProperty('signature_type')
    expect(saved).not.toHaveProperty('signature_app_id')
  })
})
