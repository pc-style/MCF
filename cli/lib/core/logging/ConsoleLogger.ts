import { ILogger, LogLevel } from "./ILogger.js";

export class ConsoleLogger implements ILogger {
  private currentLogLevel: LogLevel = LogLevel.INFO;

  debug(message: string, ...args: unknown[]): void {
    if (this.isLevelEnabled(LogLevel.DEBUG)) {
      console.debug(`[DEBUG] ${message}`, ...args);
    }
  }

  info(message: string, ...args: unknown[]): void {
    if (this.isLevelEnabled(LogLevel.INFO)) {
      console.info(`[INFO] ${message}`, ...args);
    }
  }

  warn(message: string, ...args: unknown[]): void {
    if (this.isLevelEnabled(LogLevel.WARN)) {
      console.warn(`[WARN] ${message}`, ...args);
    }
  }

  error(message: string, ...args: unknown[]): void {
    if (this.isLevelEnabled(LogLevel.ERROR)) {
      console.error(`[ERROR] ${message}`, ...args);
    }
  }

  setLogLevel(level: LogLevel): void {
    this.currentLogLevel = level;
  }

  getLogLevel(): LogLevel {
    return this.currentLogLevel;
  }

  private isLevelEnabled(level: LogLevel): boolean {
    const logLevelPriority = {
      [LogLevel.DEBUG]: 0,
      [LogLevel.INFO]: 1,
      [LogLevel.WARN]: 2,
      [LogLevel.ERROR]: 3,
    };

    return logLevelPriority[level] >= logLevelPriority[this.currentLogLevel];
  }
}
