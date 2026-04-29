import { z } from 'zod'

/**
 * Configuration form schema
 * Validates config key, value, and optional description
 */
export const configSchema = z.object({
  key: z
    .string()
    .min(1, '配置键不能为空')
    .max(100, '配置键不能超过100个字符')
    .regex(/^[a-zA-Z0-9_]+$/, '配置键只能包含字母、数字和下划线'),
  value: z.string().min(1, '配置值不能为空'),
  description: z.string().max(200, '描述不能超过200个字符').optional(),
})

export type ConfigInput = z.infer<typeof configSchema>
