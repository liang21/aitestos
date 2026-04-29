import { Tag } from '@arco-design/web-react'

/**
 * Status category determines which color mapping to use
 */
export type StatusCategory =
  | 'case_status'
  | 'plan_status'
  | 'task_status'
  | 'draft_status'
  | 'priority'
  | 'confidence'
  | 'case_type'
  | 'document_type'
  | 'document_status'
  | 'result_status'
  | 'scene_type'
  | 'user_role'

interface ColorMapping {
  color: string
  text: string
  bg: string
}

/**
 * Complete color mapping for all status categories
 * Based on UX design spec v1.0
 */
const COLOR_MAP: Record<StatusCategory, Record<string, ColorMapping>> = {
  case_status: {
    unexecuted: {
      color: '#86909C',
      text: '未执行',
      bg: 'rgba(134,144,156,0.10)',
    },
    pass: { color: '#00B42A', text: '通过', bg: 'rgba(0,180,42,0.10)' },
    block: { color: '#FF7D00', text: '阻塞', bg: 'rgba(255,125,0,0.10)' },
    fail: { color: '#F53F3F', text: '失败', bg: 'rgba(245,63,63,0.10)' },
  },
  plan_status: {
    draft: { color: '#86909C', text: '草稿', bg: 'rgba(134,144,156,0.10)' },
    active: { color: '#165DFF', text: '进行中', bg: 'rgba(22,93,255,0.10)' },
    completed: { color: '#00B42A', text: '已完成', bg: 'rgba(0,180,42,0.10)' },
    archived: {
      color: '#C9CDD4',
      text: '已归档',
      bg: 'rgba(201,205,212,0.10)',
    },
  },
  task_status: {
    pending: { color: '#86909C', text: '待处理', bg: 'rgba(134,144,156,0.10)' },
    processing: {
      color: '#165DFF',
      text: '处理中',
      bg: 'rgba(22,93,255,0.10)',
    },
    completed: { color: '#00B42A', text: '已完成', bg: 'rgba(0,180,42,0.10)' },
    failed: { color: '#F53F3F', text: '失败', bg: 'rgba(245,63,63,0.10)' },
  },
  draft_status: {
    pending: { color: '#FF7D00', text: '待确认', bg: 'rgba(255,125,0,0.10)' },
    confirmed: { color: '#00B42A', text: '已确认', bg: 'rgba(0,180,42,0.10)' },
    rejected: { color: '#F53F3F', text: '已拒绝', bg: 'rgba(245,63,63,0.10)' },
  },
  priority: {
    P0: { color: '#F53F3F', text: 'P0 紧急', bg: 'rgba(245,63,63,0.10)' },
    P1: { color: '#FF7D00', text: 'P1 高', bg: 'rgba(255,125,0,0.10)' },
    P2: { color: '#7B61FF', text: 'P2 中', bg: 'rgba(123,97,255,0.10)' },
    P3: { color: '#86909C', text: 'P3 低', bg: 'rgba(134,144,156,0.10)' },
  },
  confidence: {
    high: { color: '#00B42A', text: '高置信度', bg: 'rgba(0,180,42,0.10)' },
    medium: { color: '#FF7D00', text: '中置信度', bg: 'rgba(255,125,0,0.10)' },
    low: { color: '#F53F3F', text: '低置信度', bg: 'rgba(245,63,63,0.10)' },
  },
  case_type: {
    functionality: {
      color: '#165DFF',
      text: '功能测试',
      bg: 'rgba(22,93,255,0.10)',
    },
    performance: {
      color: '#7B61FF',
      text: '性能测试',
      bg: 'rgba(123,97,255,0.10)',
    },
    api: { color: '#00B42A', text: 'API 测试', bg: 'rgba(0,180,42,0.10)' },
    ui: { color: '#FF7D00', text: 'UI 测试', bg: 'rgba(255,125,0,0.10)' },
    security: {
      color: '#F53F3F',
      text: '安全测试',
      bg: 'rgba(245,63,63,0.10)',
    },
  },
  document_type: {
    prd: { color: '#165DFF', text: 'PRD', bg: 'rgba(22,93,255,0.10)' },
    figma: { color: '#F53F3F', text: 'Figma', bg: 'rgba(245,63,63,0.10)' },
    api_spec: { color: '#00B42A', text: 'API Spec', bg: 'rgba(0,180,42,0.10)' },
    swagger: { color: '#7B61FF', text: 'Swagger', bg: 'rgba(123,97,255,0.10)' },
    markdown: {
      color: '#86909C',
      text: 'Markdown',
      bg: 'rgba(134,144,156,0.10)',
    },
  },
  document_status: {
    pending: { color: '#86909C', text: '待处理', bg: 'rgba(134,144,156,0.10)' },
    processing: {
      color: '#165DFF',
      text: '解析中',
      bg: 'rgba(22,93,255,0.10)',
    },
    completed: { color: '#00B42A', text: '已完成', bg: 'rgba(0,180,42,0.10)' },
    failed: { color: '#F53F3F', text: '失败', bg: 'rgba(245,63,63,0.10)' },
  },
  result_status: {
    pass: { color: '#00B42A', text: '通过', bg: 'rgba(0,180,42,0.10)' },
    fail: { color: '#F53F3F', text: '失败', bg: 'rgba(245,63,63,0.10)' },
    block: { color: '#FF7D00', text: '阻塞', bg: 'rgba(255,125,0,0.10)' },
    skip: { color: '#86909C', text: '跳过', bg: 'rgba(134,144,156,0.10)' },
  },
  scene_type: {
    positive: { color: '#00B42A', text: '正向场景', bg: 'rgba(0,180,42,0.10)' },
    negative: {
      color: '#F53F3F',
      text: '负向场景',
      bg: 'rgba(245,63,63,0.10)',
    },
    boundary: {
      color: '#FF7D00',
      text: '边界场景',
      bg: 'rgba(255,125,0,0.10)',
    },
  },
  user_role: {
    super_admin: {
      color: '#F53F3F',
      text: '超级管理员',
      bg: 'rgba(245,63,63,0.10)',
    },
    admin: { color: '#FF7D00', text: '管理员', bg: 'rgba(255,125,0,0.10)' },
    normal: {
      color: '#86909C',
      text: '普通用户',
      bg: 'rgba(134,144,156,0.10)',
    },
  },
}

export interface StatusTagProps {
  /** Status enum value, e.g., 'pass' | 'fail' | 'pending' */
  status: string
  /** Status category determines color mapping */
  category: StatusCategory
  /** Custom label text (overrides default mapping) */
  label?: string
  /** Tag size */
  size?: 'small' | 'default' | 'large'
  /** Additional CSS class name */
  className?: string
}

/**
 * Unified status tag component with consistent color mapping
 * Used across all status displays in the application
 */
export function StatusTag({
  status,
  category,
  label,
  size = 'small',
  className,
}: StatusTagProps) {
  const mapping = COLOR_MAP[category]?.[status]
  if (!mapping) return null

  return (
    <Tag
      size={size}
      className={className}
      style={{
        color: mapping.color,
        backgroundColor: mapping.bg,
        borderColor: 'transparent',
        fontWeight: category === 'confidence' ? 500 : undefined,
      }}
    >
      {label ?? mapping.text}
    </Tag>
  )
}
