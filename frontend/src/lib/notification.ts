/**
 * Notification Utility
 *
 * Simple notification implementation using native browser APIs.
 * Compatible with React 19 and all modern browsers.
 */

type MessageConfigType = 'success' | 'error' | 'warning' | 'info'

interface MessageConfig {
  content?: string
  duration?: number
  closable?: boolean
  [key: string]: unknown
}

// Store for message elements
let messageContainer: HTMLElement | null = null

/**
 * Ensure message container exists
 */
function ensureContainer() {
  if (!messageContainer) {
    messageContainer = document.createElement('div')
    messageContainer.id = 'arco-message-container'
    messageContainer.style.cssText = `
      position: fixed;
      top: 16px;
      left: 50%;
      transform: translateX(-50%);
      z-index: 10000;
      display: flex;
      flex-direction: column;
      gap: 8px;
      pointer-events: none;
    `
    document.body.appendChild(messageContainer)
  }
  return messageContainer
}

/**
 * Create a message element
 */
function createMessageElement(type: MessageConfigType, content: string, duration: number): HTMLElement {
  const element = document.createElement('div')

  const colors = {
    success: '#00b42a',
    error: '#f53f3f',
    warning: '#ff7d00',
    info: '#165dff',
  }

  const icons = {
    success: '✓',
    error: '✕',
    warning: '!',
    info: 'i',
  }

  element.style.cssText = `
    background: white;
    border-left: 4px solid ${colors[type]};
    border-radius: 4px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    padding: 12px 16px;
    display: flex;
    align-items: center;
    gap: 12px;
    min-width: 300px;
    max-width: 480px;
    pointer-events: auto;
    animation: slideIn 0.3s ease-out;
    font-size: 14px;
    color: #1d2129;
  `

  element.innerHTML = `
    <span style="color: ${colors[type]}; font-weight: bold; font-size: 16px;">${icons[type]}</span>
    <span style="flex: 1;">${content}</span>
    <span style="cursor: pointer; color: #86909c; margin-left: 8px;">✕</span>
  `

  // Add close button functionality
  const closeBtn = element.querySelector('span:last-child') as HTMLElement
  closeBtn.style.cursor = 'pointer'
  closeBtn.onclick = () => {
    element.remove()
  }

  // Add animation keyframes
  if (!document.getElementById('arco-message-animations')) {
    const style = document.createElement('style')
    style.id = 'arco-message-animations'
    style.textContent = `
      @keyframes slideIn {
        from { opacity: 0; transform: translateY(-20px); }
        to { opacity: 1; transform: translateY(0); }
      }
      @keyframes fadeOut {
        from { opacity: 1; }
        to { opacity: 0; }
      }
    `
    document.head.appendChild(style)
  }

  // Auto remove after duration
  setTimeout(() => {
    element.style.animation = 'fadeOut 0.3s ease-out forwards'
    setTimeout(() => element.remove(), 300)
  }, duration)

  return element
}

/**
 * Show a success message
 */
export function messageSuccess(content: string, duration = 2000) {
  const container = ensureContainer()
  const element = createMessageElement('success', content, duration)
  container.appendChild(element)
}

/**
 * Show an error message
 */
export function messageError(content: string, duration = 3000) {
  const container = ensureContainer()
  const element = createMessageElement('error', content, duration)
  container.appendChild(element)
}

/**
 * Show a warning message
 */
export function messageWarning(content: string, duration = 3000) {
  const container = ensureContainer()
  const element = createMessageElement('warning', content, duration)
  container.appendChild(element)
}

/**
 * Show an info message
 */
export function messageInfo(content: string, duration = 3000) {
  const container = ensureContainer()
  const element = createMessageElement('info', content, duration)
  container.appendChild(element)
}

/**
 * Show a message with custom config
 */
export function message(type: MessageConfigType, config: MessageConfig) {
  const content = config.content || ''
  const duration = config.duration || 3000

  switch (type) {
    case 'success':
      messageSuccess(content, duration)
      break
    case 'error':
      messageError(content, duration)
      break
    case 'warning':
      messageWarning(content, duration)
      break
    case 'info':
      messageInfo(content, duration)
      break
  }
}
