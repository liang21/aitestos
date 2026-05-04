/**
 * Notification Utility
 *
 * Provides a compatibility layer for Arco Design Message API with React 19.
 * Uses Modal.message which is compatible with React 19's rendering system.
 */

import { Modal } from '@arco-design/web-react'

type MessageConfigType = 'success' | 'error' | 'warning' | 'info'

interface MessageConfig {
  content?: string
  duration?: number
  closable?: boolean
  [key: string]: unknown
}

/**
 * Show a success message
 */
export function messageSuccess(content: string, duration = 2000) {
  Modal.message?.success({ content, duration })
}

/**
 * Show an error message
 */
export function messageError(content: string, duration = 3000) {
  Modal.message?.error({ content, duration })
}

/**
 * Show a warning message
 */
export function messageWarning(content: string, duration = 3000) {
  Modal.message?.warning({ content, duration })
}

/**
 * Show an info message
 */
export function messageInfo(content: string, duration = 3000) {
  Modal.message?.info({ content, duration })
}

/**
 * Show a message with custom config
 */
export function message(type: MessageConfigType, config: MessageConfig) {
  Modal.message?.[type]?.(config)
}
