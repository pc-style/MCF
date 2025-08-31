import { ILogger, LogLevel } from "./ILogger.js";
import { ConsoleLogger } from "./ConsoleLogger.js";

export class LoggerFactory {
  private static defaultLogger?: ILogger;
  private static customLogger?: ILogger;

  static createLogger(name?: string): ILogger {
    if (LoggerFactory.customLogger) {
      return LoggerFactory.customLogger;
    }

    if (!LoggerFactory.defaultLogger) {
      LoggerFactory.defaultLogger = new ConsoleLogger();
    }

    return LoggerFactory.defaultLogger;
  }

  static setCustomLogger(logger: ILogger): void {
    LoggerFactory.customLogger = logger;
  }

  static setGlobalLogLevel(level: LogLevel): void {
    if (LoggerFactory.defaultLogger) {
      LoggerFactory.defaultLogger.setLogLevel(level);
    }
    if (LoggerFactory.customLogger) {
      LoggerFactory.customLogger.setLogLevel(level);
    }
  }
}
