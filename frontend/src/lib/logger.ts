/**
 * Logger Utility
 *
 * Provides structured logging with environment-aware output.
 * In development: logs to console with colors and timestamps.
 * In production: sends to error tracking service (TODO).
 */

export enum LogLevel {
  DEBUG = 0,
  INFO = 1,
  WARN = 2,
  ERROR = 3,
}

interface LogEntry {
  level: LogLevel
  timestamp: string
  message: string
  context?: Record<string, unknown>
  error?: Error
}

class Logger {
  private minLevel: LogLevel
  private isDev: boolean

  constructor() {
    this.isDev = import.meta.env.DEV
    this.minLevel = this.isDev ? LogLevel.DEBUG : LogLevel.WARN
  }

  private shouldLog(level: LogLevel): boolean {
    return level >= this.minLevel
  }

  private formatMessage(entry: LogEntry): string {
    const { level, timestamp, message, context, error } = entry

    const levelNames = ['DEBUG', 'INFO', 'WARN', 'ERROR']
    const levelName = levelNames[level]

    let formatted = `[${timestamp}] [${levelName}] ${message}`

    if (context) {
      formatted += ` ${JSON.stringify(context)}`
    }

    if (error) {
      formatted += `\n  Error: ${error.message}`
      if (error.stack && this.isDev) {
        formatted += `\n  Stack: ${error.stack}`
      }
    }

    return formatted
  }

  private log(entry: LogEntry) {
    if (!this.shouldLog(entry.level)) {
      return
    }

    const formatted = this.formatMessage(entry)

    // Console output based on level
    switch (entry.level) {
      case LogLevel.DEBUG:
        console.debug(formatted)
        break
      case LogLevel.INFO:
        console.info(formatted)
        break
      case LogLevel.WARN:
        console.warn(formatted)
        break
      case LogLevel.ERROR:
        console.error(formatted)
        break
    }

    // In production, send to error tracking service
    if (!this.isDev && entry.level >= LogLevel.ERROR) {
      this.sendToErrorTracking(entry)
    }
  }

  private sendToErrorTracking(entry: LogEntry) {
    // TODO: Integrate with error tracking service (e.g., Sentry, LogRocket)
    // For now, we'll just keep the console output
    // Example integration:
    // Sentry.captureException(entry.error, {
    //   level: 'error',
    //   extra: entry.context,
    //   tags: { component: entry.context?.component },
    // })
  }

  debug(message: string, context?: Record<string, unknown>) {
    this.log({
      level: LogLevel.DEBUG,
      timestamp: new Date().toISOString(),
      message,
      context,
    })
  }

  info(message: string, context?: Record<string, unknown>) {
    this.log({
      level: LogLevel.INFO,
      timestamp: new Date().toISOString(),
      message,
      context,
    })
  }

  warn(message: string, context?: Record<string, unknown>) {
    this.log({
      level: LogLevel.WARN,
      timestamp: new Date().toISOString(),
      message,
      context,
    })
  }

  error(
    message: string,
    error?: Error | unknown,
    context?: Record<string, unknown>
  ) {
    const errorObj = error instanceof Error ? error : new Error(String(error))

    this.log({
      level: LogLevel.ERROR,
      timestamp: new Date().toISOString(),
      message,
      error: errorObj,
      context,
    })
  }
}

// Singleton instance
export const logger = new Logger()

/**
 * Convenience function for logging in auth context
 */
export function logAuthError(message: string, error: unknown) {
  logger.error(message, error, { module: 'auth' })
}

/**
 * Convenience function for logging in API context
 */
export function logApiError(endpoint: string, error: unknown) {
  logger.error(`API call failed: ${endpoint}`, error, {
    module: 'api',
    endpoint,
  })
}
