// Generic type constructor for runtime type checking
export type Type<T> = new (...args: any[]) => T;

// Utility types for CLI operations
export type CLIOutput = {
  success: boolean;
  message: string;
  data?: unknown;
};

export type CLIEnvironment = "development" | "production" | "staging" | "test";

export interface CLIConfig {
  environment: CLIEnvironment;
  verbose?: boolean;
  debug?: boolean;
}

export type CLIPermission =
  | "read"
  | "write"
  | "execute"
  | "admin"
  | "network"
  | "system";

// Generic error type for CLI operations
export class CLIError extends Error {
  constructor(
    message: string,
    public code?: string,
    public details?: Record<string, unknown>,
  ) {
    super(message);
    this.name = "CLIError";
  }
}
