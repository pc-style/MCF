import { ConsoleLogger } from "./ConsoleLogger.js";

export class LoggerFactory {
  /** @type {Object} */
  static defaultLogger;

  /** @type {Object} */
  static customLogger;

  /**
   * Create a logger instance
   * @param {string} [name]
   * @returns {Object}
   */
  static createLogger(name) {
    if (LoggerFactory.customLogger) {
      return LoggerFactory.customLogger;
    }

    if (!LoggerFactory.defaultLogger) {
      LoggerFactory.defaultLogger = new ConsoleLogger();
    }

    return LoggerFactory.defaultLogger;
  }

  /**
   * Set a custom logger
   * @param {Object} logger
   */
  static setCustomLogger(logger) {
    LoggerFactory.customLogger = logger;
  }

  /**
   * Set global log level
   * @param {string} level
   */
  static setGlobalLogLevel(level) {
    if (LoggerFactory.defaultLogger) {
      LoggerFactory.defaultLogger.setLogLevel(level);
    }
    if (LoggerFactory.customLogger) {
      LoggerFactory.customLogger.setLogLevel(level);
    }
  }

  /**
   * Get logger for a specific name (alias for createLogger)
   * @param {string} name
   * @returns {Object}
   */
  static getLogger(name) {
    return this.createLogger(name);
  }
}
