/**
 * MCP server information
 */
export interface MCPServerInfo {
  id: string;
  name: string;
  description?: string;
  version?: string;
  command: string;
  args?: string[];
  env?: Record<string, string>;
  cwd?: string;
  port?: number;
  status: MCPServerStatus;
  pid?: number;
  lastStarted?: Date;
  lastStopped?: Date;
  autoStart?: boolean;
  dependencies?: string[];
}

/**
 * MCP server status
 */
export type MCPServerStatus = "stopped" | "starting" | "running" | "stopping" | "error";

/**
 * MCP server configuration
 */
export interface MCPServerConfig {
  id: string;
  name: string;
  command: string;
  args?: string[];
  env?: Record<string, string>;
  cwd?: string;
  port?: number;
  autoStart?: boolean;
  restartPolicy?: RestartPolicy;
  healthCheck?: HealthCheckConfig;
  dependencies?: string[];
}

/**
 * Restart policy for MCP servers
 */
export type RestartPolicy = "never" | "always" | "on-failure";

/**
 * Health check configuration
 */
export interface HealthCheckConfig {
  enabled: boolean;
  interval: number; // in milliseconds
  timeout: number; // in milliseconds
  retries: number;
  endpoint?: string;
  command?: string;
}

/**
 * MCP server start options
 */
export interface MCPServerStartOptions {
  waitForReady?: boolean;
  timeout?: number;
  env?: Record<string, string>;
}

/**
 * MCP server health status
 */
export interface MCPHealthStatus {
  serverId: string;
  status: MCPServerStatus;
  healthy: boolean;
  responseTime?: number;
  lastCheck: Date;
  error?: string;
}

/**
 * MCP service interface for MCF CLI
 * Handles MCP server lifecycle management and orchestration
 */
export interface IMCPService {
  /**
   * Register an MCP server configuration
   * @param config Server configuration
   */
  registerServer(config: MCPServerConfig): Promise<void>;

  /**
   * Unregister an MCP server
   * @param serverId Server identifier
   */
  unregisterServer(serverId: string): Promise<boolean>;

  /**
   * Get MCP server information
   * @param serverId Server identifier
   */
  getServer(serverId: string): Promise<MCPServerInfo | null>;

  /**
   * List all registered MCP servers
   */
  listServers(): Promise<MCPServerInfo[]>;

  /**
   * Start an MCP server
   * @param serverId Server identifier
   * @param options Start options
   */
  startServer(serverId: string, options?: MCPServerStartOptions): Promise<boolean>;

  /**
   * Stop an MCP server
   * @param serverId Server identifier
   * @param force Force stop (SIGKILL)
   */
  stopServer(serverId: string, force?: boolean): Promise<boolean>;

  /**
   * Restart an MCP server
   * @param serverId Server identifier
   */
  restartServer(serverId: string): Promise<boolean>;

  /**
   * Check if an MCP server is running
   * @param serverId Server identifier
   */
  isServerRunning(serverId: string): Promise<boolean>;

  /**
   * Get MCP server health status
   * @param serverId Server identifier
   */
  getServerHealth(serverId: string): Promise<MCPHealthStatus>;

  /**
   * Get health status for all servers
   */
  getAllServerHealth(): Promise<MCPHealthStatus[]>;

  /**
   * Start all auto-start MCP servers
   */
  startAutoStartServers(): Promise<string[]>;

  /**
   * Stop all running MCP servers
   * @param force Force stop all servers
   */
  stopAllServers(force?: boolean): Promise<string[]>;

  /**
   * Install an MCP server from a package or repository
   * @param packageName Package name or repository URL
   * @param options Installation options
   */
  installServer(packageName: string, options?: MCPInstallOptions): Promise<MCPServerInfo>;

  /**
   * Update an MCP server
   * @param serverId Server identifier
   */
  updateServer(serverId: string): Promise<boolean>;

  /**
   * Remove an MCP server installation
   * @param serverId Server identifier
   * @param keepConfig Keep configuration files
   */
  removeServer(serverId: string, keepConfig?: boolean): Promise<boolean>;

  /**
   * Get MCP server logs
   * @param serverId Server identifier
   * @param lines Number of lines to retrieve
   */
  getServerLogs(serverId: string, lines?: number): Promise<string[]>;

  /**
   * Validate MCP server configuration
   * @param config Server configuration
   */
  validateServerConfig(config: MCPServerConfig): Promise<MCPValidationResult>;

  /**
   * Get MCP service statistics
   */
  getServiceStats(): Promise<MCPServiceStats>;

  /**
   * Export MCP server configurations
   */
  exportConfigurations(): Promise<string>;

  /**
   * Import MCP server configurations
   * @param configData Configuration data
   */
  importConfigurations(configData: string): Promise<string[]>;
}

/**
 * MCP server installation options
 */
export interface MCPInstallOptions {
  version?: string;
  global?: boolean;
  dependencies?: string[];
  postInstallScript?: string;
}

/**
 * MCP server validation result
 */
export interface MCPValidationResult {
  isValid: boolean;
  errors: string[];
  warnings: string[];
  suggestions: string[];
}

/**
 * MCP service statistics
 */
export interface MCPServiceStats {
  totalServers: number;
  runningServers: number;
  stoppedServers: number;
  erroredServers: number;
  autoStartServers: number;
  totalUptime: number; // in milliseconds
  averageHealthCheckTime: number; // in milliseconds
  serverStats: {
    [serverId: string]: {
      uptime: number;
      restarts: number;
      healthChecks: number;
      successfulHealthChecks: number;
    };
  };
}

/**
 * MCP service configuration
 */
export interface MCPServiceConfig {
  configDirectory?: string;
  serversDirectory?: string;
  logsDirectory?: string;
  defaultTimeout?: number;
  healthCheckInterval?: number;
  maxConcurrentStarts?: number;
  autoStartDelay?: number; // delay between auto-starting servers
  enableHealthChecks?: boolean;
  logRotation?: {
    enabled: boolean;
    maxSize: number; // in bytes
    maxFiles: number;
  };
}


