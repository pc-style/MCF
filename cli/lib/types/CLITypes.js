// Generic type constructor for runtime type checking
export const Type = (T) => T;

// Utility types for CLI operations
export const CLIOutput = {
  success: false,
  message: "",
  data: undefined
};

export const CLIEnvironment = {
  DEVELOPMENT: "development",
  PRODUCTION: "production",
  STAGING: "staging",
  TEST: "test"
};

export const CLIConfig = {
  environment: CLIEnvironment.DEVELOPMENT,
  verbose: false,
  debug: false
};

export const CLIPermission = {
  READ: "read",
  WRITE: "write",
  EXECUTE: "execute",
  ADMIN: "admin",
  NETWORK: "network",
  SYSTEM: "system"
};

// Generic error type for CLI operations
export class CLIError extends Error {
  constructor(message, code, details) {
    super(message);
    this.code = code;
    this.details = details;
    this.name = "CLIError";
  }
}



