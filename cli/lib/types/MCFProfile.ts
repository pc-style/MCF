import { CLIEnvironment } from "./CLITypes.js";

export interface MCFProfile {
  // Basic profile identification
  id: string;
  name: string;
  description?: string;

  // Environment configuration
  environment: CLIEnvironment;

  // Configuration options
  config: {
    // Runtime settings
    timeout?: number;
    maxRetries?: number;

    // Logging and tracing
    logLevel?: "debug" | "info" | "warn" | "error";
    traceId?: string;

    // Network and connection settings
    proxy?: {
      host: string;
      port: number;
      protocol?: "http" | "https";
    };

    // Authentication
    auth?: {
      type: "none" | "basic" | "oauth" | "token";
      credentials?: Record<string, string>;
    };
  };

  // Permissions and security
  permissions?: {
    allowedServices?: string[];
    blockedServices?: string[];
    roleBasedAccess?: string[];
  };

  // Metadata and versioning
  version?: string;
  lastUpdated?: Date;
}
