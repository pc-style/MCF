export const LogLevel = {
  DEBUG: "DEBUG",
  INFO: "INFO",
  WARN: "WARN",
  ERROR: "ERROR",
};

// ILogger interface for JavaScript runtime
export class ILogger {
  debug(message, ...args) {
    throw new Error("Method 'debug' must be implemented");
  }

  info(message, ...args) {
    throw new Error("Method 'info' must be implemented");
  }

  warn(message, ...args) {
    throw new Error("Method 'warn' must be implemented");
  }

  error(message, ...args) {
    throw new Error("Method 'error' must be implemented");
  }

  setLogLevel(level) {
    throw new Error("Method 'setLogLevel' must be implemented");
  }

  getLogLevel() {
    throw new Error("Method 'getLogLevel' must be implemented");
  }
}
