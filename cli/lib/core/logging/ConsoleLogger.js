import { LogLevel } from "./ILogger.js";

export class ConsoleLogger {
  constructor() {
    this.currentLogLevel = LogLevel.INFO;
  }

  debug(message, ...args) {
    if (this.isLevelEnabled(LogLevel.DEBUG)) {
      console.debug(`[DEBUG] ${message}`, ...args);
    }
  }

  info(message, ...args) {
    if (this.isLevelEnabled(LogLevel.INFO)) {
      console.info(`[INFO] ${message}`, ...args);
    }
  }

  warn(message, ...args) {
    if (this.isLevelEnabled(LogLevel.WARN)) {
      console.warn(`[WARN] ${message}`, ...args);
    }
  }

  error(message, ...args) {
    if (this.isLevelEnabled(LogLevel.ERROR)) {
      console.error(`[ERROR] ${message}`, ...args);
    }
  }

  setLogLevel(level) {
    this.currentLogLevel = level;
  }

  getLogLevel() {
    return this.currentLogLevel;
  }

  isLevelEnabled(level) {
    const logLevelPriority = {
      [LogLevel.DEBUG]: 0,
      [LogLevel.INFO]: 1,
      [LogLevel.WARN]: 2,
      [LogLevel.ERROR]: 3,
    };

    return logLevelPriority[level] >= logLevelPriority[this.currentLogLevel];
  }
}
