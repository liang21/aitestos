import { describe, it, expect, beforeEach, vi } from 'vitest'

describe('logger', () => {
  let consoleDebugSpy: ReturnType<typeof vi.spyOn>
  let consoleInfoSpy: ReturnType<typeof vi.spyOf>
  let consoleWarnSpy: ReturnType<typeof vi.spyOn>
  let consoleErrorSpy: ReturnType<typeof vi.spyOn>

  beforeEach(() => {
    consoleDebugSpy = vi.spyOn(console, 'debug').mockImplementation(() => {})
    consoleInfoSpy = vi.spyOn(console, 'info').mockImplementation(() => {})
    consoleWarnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})
    consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
  })

  // Import logger after setting up mocks
  const { logger, logAuthError, logApiError } = require('./logger')

  it('should log debug messages in development', () => {
    logger.debug('Test debug message', { key: 'value' })

    expect(consoleDebugSpy).toHaveBeenCalled()
    const loggedMessage = consoleDebugSpy.mock.calls[0][0] as string
    expect(loggedMessage).toContain('DEBUG')
    expect(loggedMessage).toContain('Test debug message')
    expect(loggedMessage).toContain('key')
  })

  it('should log info messages', () => {
    logger.info('Test info message')

    expect(consoleInfoSpy).toHaveBeenCalled()
    const loggedMessage = consoleInfoSpy.mock.calls[0][0] as string
    expect(loggedMessage).toContain('INFO')
    expect(loggedMessage).toContain('Test info message')
  })

  it('should log warning messages', () => {
    logger.warn('Test warning message')

    expect(consoleWarnSpy).toHaveBeenCalled()
    const loggedMessage = consoleWarnSpy.mock.calls[0][0] as string
    expect(loggedMessage).toContain('WARN')
    expect(loggedMessage).toContain('Test warning message')
  })

  it('should log error messages with Error object', () => {
    const error = new Error('Test error')
    logger.error('Test error message', error)

    expect(consoleErrorSpy).toHaveBeenCalled()
    const loggedMessage = consoleErrorSpy.mock.calls[0][0] as string
    expect(loggedMessage).toContain('ERROR')
    expect(loggedMessage).toContain('Test error message')
    expect(loggedMessage).toContain('Error: Test error')
  })

  it('should log error messages with non-Error objects', () => {
    logger.error('Test error message', 'string error')

    expect(consoleErrorSpy).toHaveBeenCalled()
    const loggedMessage = consoleErrorSpy.mock.calls[0][0] as string
    expect(loggedMessage).toContain('ERROR')
    expect(loggedMessage).toContain('Error: string error')
  })

  it('should include context in log messages', () => {
    logger.info('Test message', { userId: '123', action: 'login' })

    expect(consoleInfoSpy).toHaveBeenCalled()
    const loggedMessage = consoleInfoSpy.mock.calls[0][0] as string
    expect(loggedMessage).toContain('userId')
    expect(loggedMessage).toContain('123')
  })

  it('should include stack trace for errors in development', () => {
    const error = new Error('Test error')
    logger.error('Test message', error)

    expect(consoleErrorSpy).toHaveBeenCalled()
    const loggedMessage = consoleErrorSpy.mock.calls[0][0] as string
    expect(loggedMessage).toContain('Stack:')
  })
})

describe('logAuthError', () => {
  it('should log auth errors with correct context', () => {
    const consoleErrorSpy = vi
      .spyOn(console, 'error')
      .mockImplementation(() => {})

    const { logAuthError } = require('./logger')
    const error = new Error('Auth failed')
    logAuthError('Login failed', error)

    expect(consoleErrorSpy).toHaveBeenCalled()
    const loggedMessage = consoleErrorSpy.mock.calls[0][0] as string
    expect(loggedMessage).toContain('auth')

    consoleErrorSpy.mockRestore()
  })
})

describe('logApiError', () => {
  it('should log API errors with correct context', () => {
    const consoleErrorSpy = vi
      .spyOn(console, 'error')
      .mockImplementation(() => {})

    const { logApiError } = require('./logger')
    const error = new Error('Network error')
    logApiError('/api/v1/auth/login', error)

    expect(consoleErrorSpy).toHaveBeenCalled()
    const loggedMessage = consoleErrorSpy.mock.calls[0][0] as string
    expect(loggedMessage).toContain('api')
    expect(loggedMessage).toContain('/api/v1/auth/login')

    consoleErrorSpy.mockRestore()
  })
})
