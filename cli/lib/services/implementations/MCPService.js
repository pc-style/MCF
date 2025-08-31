import { BaseService } from "../../core/base/BaseService.js";
import { ILogger } from "../../core/logging/ILogger.js";
import { LoggerFactory } from "../../core/logging/LoggerFactory.js";
import { CLIError } from "../../types/CLITypes.js";
import { spawn } from "child_process";
import fs from "fs/promises";
import path from "path";
import http from "http";
import https from "https";

/**
 * MCP service implementation for MCF CLI
 * Handles MCP server lifecycle management and orchestration
 */
export class MCPService extends BaseService {
  constructor(config, logger) {
    super();
    this.config = config || {};
    this.logger = logger || LoggerFactory.getLogger("MCPService");

    // Default configuration
    this.configDirectory = config?.configDirectory || path.join(process.cwd(), ".mcf", "mcp");
    this.serversDirectory = config?.serversDirectory || path.join(this.configDirectory, "servers");
    this.logsDirectory = config?.logsDirectory || path.join(this.configDirectory, "logs");
    this.defaultTimeout = config?.defaultTimeout || 30000;
    this.healthCheckInterval = config?.healthCheckInterval || 30000;
    this.maxConcurrentStarts = config?.maxConcurrentStarts || 3;
    this.autoStartDelay = config?.autoStartDelay || 2000;
    this.enableHealthChecks = config?.enableHealthChecks !== false;

    // Track running servers and health checks
    this.runningServers = new Map();
    this.serverConfigs = new Map();
    this.healthCheckTimers = new Map();

    this.logger.debug("MCPService initialized", {
      configDirectory: this.configDirectory,
      serversDirectory: this.serversDirectory
    });
  }

  /**
   * Register an MCP server configuration
   */
  async registerServer(config) {
    try {
      // Validate configuration
      const validation = await this.validateServerConfig(config);
      if (!validation.isValid) {
        throw new CLIError(
          `MCP server configuration validation failed: ${validation.errors.join(", ")}`,
          "MCP_CONFIG_INVALID",
          { errors: validation.errors, warnings: validation.warnings }
        );
      }

      if (validation.warnings.length > 0) {
        this.logger.warn(`MCP server configuration warnings: ${validation.warnings.join(", ")}`);
      }

      // Ensure directories exist
      await this.ensureDirectories();

      // Save configuration
      const configPath = path.join(this.configDirectory, "servers", `${config.id}.json`);
      await fs.writeFile(configPath, JSON.stringify(config, null, 2), "utf-8");

      // Register in memory
      this.serverConfigs.set(config.id, config);

      // Set up health monitoring if enabled
      if (this.enableHealthChecks && config.healthCheck?.enabled) {
        this.setupHealthCheck(config.id);
      }

      this.logger.info(`MCP server '${config.id}' registered successfully`);
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to register MCP server '${config.id}': ${message}`);
      throw new CLIError(
        `Failed to register MCP server: ${message}`,
        "MCP_SERVER_REGISTRATION_FAILED",
        { serverId: config.id }
      );
    }
  }

  /**
   * Unregister an MCP server
   */
  async unregisterServer(serverId) {
    try {
      // Stop server if running
      if (await this.isServerRunning(serverId)) {
        await this.stopServer(serverId, true);
      }

      // Remove health check
      this.removeHealthCheck(serverId);

      // Remove configuration file
      const configPath = path.join(this.configDirectory, "servers", `${serverId}.json`);
      try {
        await fs.unlink(configPath);
      } catch {
        // Ignore if file doesn't exist
      }

      // Remove from memory
      this.serverConfigs.delete(serverId);

      this.logger.info(`MCP server '${serverId}' unregistered successfully`);
      return true;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to unregister MCP server '${serverId}': ${message}`);
      throw new CLIError(
        `Failed to unregister MCP server: ${message}`,
        "MCP_SERVER_UNREGISTRATION_FAILED",
        { serverId }
      );
    }
  }

  /**
   * Get MCP server information
   */
  async getServer(serverId) {
    try {
      // Try memory cache first
      let config = this.serverConfigs.get(serverId);

      if (!config) {
        // Load from file
        const configPath = path.join(this.configDirectory, "servers", `${serverId}.json`);
        try {
          const configData = await fs.readFile(configPath, "utf-8");
          config = JSON.parse(configData);
          this.serverConfigs.set(serverId, config);
        } catch {
          return null;
        }
      }

      // Get current status
      const isRunning = await this.isServerRunning(serverId);
      const serverInfo = this.runningServers.get(serverId);

      return {
        ...config,
        status: isRunning ? "running" : "stopped",
        pid: serverInfo?.process?.pid,
        lastStarted: serverInfo?.startTime,
        lastStopped: serverInfo?.stopTime
      };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get MCP server '${serverId}': ${message}`);
      throw new CLIError(
        `Failed to get MCP server: ${message}`,
        "MCP_SERVER_GET_FAILED",
        { serverId }
      );
    }
  }

  /**
   * List all registered MCP servers
   */
  async listServers() {
    try {
      await this.ensureDirectories();

      const servers = [];
      const serverFiles = await fs.readdir(path.join(this.configDirectory, "servers"));

      for (const file of serverFiles) {
        if (file.endsWith(".json")) {
          const serverId = file.replace(".json", "");
          const server = await this.getServer(serverId);
          if (server) {
            servers.push(server);
          }
        }
      }

      return servers;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to list MCP servers: ${message}`);
      throw new CLIError(
        `Failed to list MCP servers: ${message}`,
        "MCP_SERVER_LIST_FAILED"
      );
    }
  }

  /**
   * Start an MCP server
   */
  async startServer(serverId, options = {}) {
    try {
      const config = await this.getServer(serverId);
      if (!config) {
        throw new CLIError(
          `MCP server '${serverId}' not found`,
          "MCP_SERVER_NOT_FOUND"
        );
      }

      // Check if already running
      if (await this.isServerRunning(serverId)) {
        this.logger.warn(`MCP server '${serverId}' is already running`);
        return true;
      }

      this.logger.info(`Starting MCP server '${serverId}'`);

      // Prepare environment
      const env = {
        ...process.env,
        ...config.env,
        ...options.env
      };

      // Prepare arguments
      const args = config.args || [];

      // Start the process
      const child = spawn(config.command, args, {
        cwd: config.cwd || process.cwd(),
        env,
        stdio: ["pipe", "pipe", "pipe"],
        shell: process.platform === "win32"
      });

      // Set up process tracking
      const serverInfo = {
        process: child,
        startTime: new Date(),
        config,
        healthCheckFailures: 0
      };

      this.runningServers.set(serverId, serverInfo);

      // Set up log capture
      this.setupLogCapture(serverId, child);

      // Handle process events
      child.on("exit", (code, signal) => {
        const exitCode = code || 0;
        const success = exitCode === 0;

        this.logger.info(`MCP server '${serverId}' exited`, {
          exitCode,
          signal,
          success
        });

        // Update server info
        const info = this.runningServers.get(serverId);
        if (info) {
          info.stopTime = new Date();
          info.exitCode = exitCode;
          info.signal = signal;
        }

        // Clean up after a delay
        setTimeout(() => {
          this.runningServers.delete(serverId);
        }, 5000);

        // Handle restart policy
        if (!success && config.restartPolicy === "on-failure") {
          this.logger.info(`Restarting MCP server '${serverId}' due to failure`);
          setTimeout(() => {
            this.startServer(serverId, options);
          }, 5000);
        } else if (config.restartPolicy === "always") {
          this.logger.info(`Restarting MCP server '${serverId}'`);
          setTimeout(() => {
            this.startServer(serverId, options);
          }, 5000);
        }
      });

      child.on("error", (error) => {
        this.logger.error(`MCP server '${serverId}' error: ${error.message}`);

        // Update status
        const info = this.runningServers.get(serverId);
        if (info) {
          info.error = error.message;
        }
      });

      // Wait for server to be ready if requested
      if (options.waitForReady) {
        await this.waitForServerReady(serverId, options.timeout || this.defaultTimeout);
      }

      this.logger.info(`MCP server '${serverId}' started successfully`);
      return true;
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to start MCP server '${serverId}': ${message}`);
      throw new CLIError(
        `Failed to start MCP server: ${message}`,
        "MCP_SERVER_START_FAILED",
        { serverId }
      );
    }
  }

  /**
   * Stop an MCP server
   */
  async stopServer(serverId, force = false) {
    try {
      const serverInfo = this.runningServers.get(serverId);
      if (!serverInfo || !serverInfo.process) {
        this.logger.warn(`MCP server '${serverId}' is not running`);
        return false;
      }

      this.logger.info(`Stopping MCP server '${serverId}'${force ? " (force)" : ""}`);

      // Send termination signal
      const signal = force ? "SIGKILL" : "SIGTERM";
      serverInfo.process.kill(signal);

      // Wait for process to exit
      await new Promise((resolve) => {
        const timeout = setTimeout(() => {
          if (force) {
            serverInfo.process.kill("SIGKILL");
          }
          resolve();
        }, force ? 1000 : 5000);

        serverInfo.process.on("exit", () => {
          clearTimeout(timeout);
          resolve();
        });
      });

      this.logger.info(`MCP server '${serverId}' stopped successfully`);
      return true;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to stop MCP server '${serverId}': ${message}`);
      throw new CLIError(
        `Failed to stop MCP server: ${message}`,
        "MCP_SERVER_STOP_FAILED",
        { serverId }
      );
    }
  }

  /**
   * Restart an MCP server
   */
  async restartServer(serverId) {
    try {
      // Stop the server
      await this.stopServer(serverId);

      // Wait a moment
      await new Promise(resolve => setTimeout(resolve, 2000));

      // Start the server
      return await this.startServer(serverId);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to restart MCP server '${serverId}': ${message}`);
      throw new CLIError(
        `Failed to restart MCP server: ${message}`,
        "MCP_SERVER_RESTART_FAILED",
        { serverId }
      );
    }
  }

  /**
   * Check if an MCP server is running
   */
  async isServerRunning(serverId) {
    const serverInfo = this.runningServers.get(serverId);
    return serverInfo && serverInfo.process && !serverInfo.process.killed;
  }

  /**
   * Get MCP server health status
   */
  async getServerHealth(serverId) {
    try {
      const server = await this.getServer(serverId);
      if (!server) {
        throw new CLIError(
          `MCP server '${serverId}' not found`,
          "MCP_SERVER_NOT_FOUND"
        );
      }

      const isRunning = await this.isServerRunning(serverId);
      const healthCheckResult = await this.performHealthCheck(serverId);

      return {
        serverId,
        status: server.status,
        healthy: healthCheckResult.healthy,
        responseTime: healthCheckResult.responseTime,
        lastCheck: new Date(),
        error: healthCheckResult.error
      };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get MCP server health '${serverId}': ${message}`);

      return {
        serverId,
        status: "error",
        healthy: false,
        lastCheck: new Date(),
        error: message
      };
    }
  }

  /**
   * Get health status for all servers
   */
  async getAllServerHealth() {
    const servers = await this.listServers();
    const healthPromises = servers.map(server => this.getServerHealth(server.id));
    return await Promise.all(healthPromises);
  }

  /**
   * Start all auto-start MCP servers
   */
  async startAutoStartServers() {
    try {
      const servers = await this.listServers();
      const autoStartServers = servers.filter(server => server.autoStart);

      if (autoStartServers.length === 0) {
        this.logger.info("No auto-start MCP servers configured");
        return [];
      }

      this.logger.info(`Starting ${autoStartServers.length} auto-start MCP servers`);

      const startedServers = [];

      // Start servers with delay between each
      for (let i = 0; i < autoStartServers.length; i++) {
        const server = autoStartServers[i];

        try {
          await this.startServer(server.id);

          // Wait before starting next server
          if (i < autoStartServers.length - 1) {
            await new Promise(resolve => setTimeout(resolve, this.autoStartDelay));
          }

          startedServers.push(server.id);
        } catch (error) {
          this.logger.error(`Failed to auto-start MCP server '${server.id}': ${error.message}`);
        }
      }

      this.logger.info(`Auto-started ${startedServers.length} MCP servers`);
      return startedServers;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to start auto-start servers: ${message}`);
      throw new CLIError(
        `Failed to start auto-start servers: ${message}`,
        "MCP_AUTO_START_FAILED"
      );
    }
  }

  /**
   * Stop all running MCP servers
   */
  async stopAllServers(force = false) {
    try {
      const runningServerIds = Array.from(this.runningServers.keys());

      if (runningServerIds.length === 0) {
        this.logger.info("No running MCP servers to stop");
        return [];
      }

      this.logger.info(`Stopping ${runningServerIds.length} MCP servers${force ? " (force)" : ""}`);

      const stopPromises = runningServerIds.map(serverId =>
        this.stopServer(serverId, force).catch(error => {
          this.logger.error(`Failed to stop MCP server '${serverId}': ${error.message}`);
          return false;
        })
      );

      const results = await Promise.all(stopPromises);
      const stoppedServers = runningServerIds.filter((_, index) => results[index]);

      this.logger.info(`Stopped ${stoppedServers.length} MCP servers`);
      return stoppedServers;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to stop all servers: ${message}`);
      throw new CLIError(
        `Failed to stop all servers: ${message}`,
        "MCP_STOP_ALL_FAILED"
      );
    }
  }

  /**
   * Install an MCP server from a package or repository
   */
  async installServer(packageName, options = {}) {
    try {
      this.logger.info(`Installing MCP server: ${packageName}`);

      // For now, create a basic server configuration
      // In a full implementation, this would handle npm installs, git clones, etc.
      const serverId = this.generateServerId(packageName);

      const serverConfig = {
        id: serverId,
        name: packageName,
        description: `MCP server for ${packageName}`,
        command: packageName,
        args: [],
        env: {},
        autoStart: options.autoStart || false,
        restartPolicy: "never",
        healthCheck: {
          enabled: true,
          interval: 30000,
          timeout: 5000,
          retries: 3
        }
      };

      await this.registerServer(serverConfig);

      this.logger.info(`MCP server '${serverId}' installed successfully`);
      return await this.getServer(serverId);
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to install MCP server '${packageName}': ${message}`);
      throw new CLIError(
        `Failed to install MCP server: ${message}`,
        "MCP_SERVER_INSTALL_FAILED",
        { packageName }
      );
    }
  }

  /**
   * Update an MCP server
   */
  async updateServer(serverId) {
    try {
      const server = await this.getServer(serverId);
      if (!server) {
        throw new CLIError(
          `MCP server '${serverId}' not found`,
          "MCP_SERVER_NOT_FOUND"
        );
      }

      // For now, just restart the server
      // In a full implementation, this would handle version updates
      if (await this.isServerRunning(serverId)) {
        await this.restartServer(serverId);
      }

      this.logger.info(`MCP server '${serverId}' updated successfully`);
      return true;
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to update MCP server '${serverId}': ${message}`);
      throw new CLIError(
        `Failed to update MCP server: ${message}`,
        "MCP_SERVER_UPDATE_FAILED",
        { serverId }
      );
    }
  }

  /**
   * Remove an MCP server installation
   */
  async removeServer(serverId, keepConfig = false) {
    try {
      // Stop server if running
      if (await this.isServerRunning(serverId)) {
        await this.stopServer(serverId, true);
      }

      // Remove server files/directory
      const server = await this.getServer(serverId);
      if (server && server.cwd) {
        try {
          await fs.rm(server.cwd, { recursive: true, force: true });
        } catch (error) {
          this.logger.warn(`Failed to remove server directory: ${error.message}`);
        }
      }

      // Remove configuration unless keeping it
      if (!keepConfig) {
        await this.unregisterServer(serverId);
      }

      this.logger.info(`MCP server '${serverId}' removed successfully`);
      return true;
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to remove MCP server '${serverId}': ${message}`);
      throw new CLIError(
        `Failed to remove MCP server: ${message}`,
        "MCP_SERVER_REMOVE_FAILED",
        { serverId }
      );
    }
  }

  /**
   * Get MCP server logs
   */
  async getServerLogs(serverId, lines = 100) {
    try {
      const logPath = path.join(this.logsDirectory, `${serverId}.log`);

      try {
        const logData = await fs.readFile(logPath, "utf-8");
        const logLines = logData.split("\n").filter(line => line.trim());

        // Return last N lines
        return logLines.slice(-lines);
      } catch {
        return ["No logs available"];
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get MCP server logs '${serverId}': ${message}`);
      throw new CLIError(
        `Failed to get server logs: ${message}`,
        "MCP_SERVER_LOGS_FAILED",
        { serverId }
      );
    }
  }

  /**
   * Validate MCP server configuration
   */
  async validateServerConfig(config) {
    const errors = [];
    const warnings = [];
    const suggestions = [];

    // Validate required fields
    if (!config.id || typeof config.id !== "string" || config.id.trim() === "") {
      errors.push("Server ID is required and must be a non-empty string");
    }

    if (!config.name || typeof config.name !== "string" || config.name.trim() === "") {
      errors.push("Server name is required and must be a non-empty string");
    }

    if (!config.command || typeof config.command !== "string" || config.command.trim() === "") {
      errors.push("Server command is required and must be a non-empty string");
    }

    // Validate restart policy
    if (config.restartPolicy && !["never", "always", "on-failure"].includes(config.restartPolicy)) {
      errors.push("Restart policy must be one of: never, always, on-failure");
    }

    // Validate health check configuration
    if (config.healthCheck) {
      if (config.healthCheck.interval && config.healthCheck.interval < 1000) {
        warnings.push("Health check interval should be at least 1000ms");
        suggestions.push("Consider health check interval between 5000ms and 60000ms");
      }

      if (config.healthCheck.timeout && config.healthCheck.timeout < 500) {
        warnings.push("Health check timeout should be at least 500ms");
      }

      if (config.healthCheck.retries && (config.healthCheck.retries < 1 || config.healthCheck.retries > 10)) {
        warnings.push("Health check retries should be between 1 and 10");
      }
    }

    return {
      isValid: errors.length === 0,
      errors,
      warnings,
      suggestions
    };
  }

  /**
   * Get MCP service statistics
   */
  async getServiceStats() {
    try {
      const servers = await this.listServers();
      const healthStatuses = await this.getAllServerHealth();

      const runningServers = servers.filter(server => server.status === "running").length;
      const stoppedServers = servers.filter(server => server.status === "stopped").length;
      const erroredServers = servers.filter(server => server.status === "error").length;
      const autoStartServers = servers.filter(server => server.autoStart).length;

      // Calculate total uptime
      let totalUptime = 0;
      const serverStats = {};

      for (const server of servers) {
        const serverId = server.id;
        const healthStatus = healthStatuses.find(h => h.serverId === serverId);

        if (server.lastStarted && server.status === "running") {
          totalUptime += Date.now() - server.lastStarted.getTime();
        }

        serverStats[serverId] = {
          uptime: server.lastStarted && server.status === "running"
            ? Date.now() - server.lastStarted.getTime()
            : 0,
          restarts: 0, // Would track actual restarts in full implementation
          healthChecks: healthStatus ? 1 : 0,
          successfulHealthChecks: healthStatus?.healthy ? 1 : 0
        };
      }

      // Calculate average health check time
      const healthCheckTimes = healthStatuses
        .filter(h => h.responseTime)
        .map(h => h.responseTime);

      const averageHealthCheckTime = healthCheckTimes.length > 0
        ? healthCheckTimes.reduce((sum, time) => sum + time, 0) / healthCheckTimes.length
        : 0;

      return {
        totalServers: servers.length,
        runningServers,
        stoppedServers,
        erroredServers,
        autoStartServers,
        totalUptime,
        averageHealthCheckTime,
        serverStats
      };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get MCP service stats: ${message}`);
      throw new CLIError(
        `Failed to get service stats: ${message}`,
        "MCP_SERVICE_STATS_FAILED"
      );
    }
  }

  /**
   * Export MCP server configurations
   */
  async exportConfigurations() {
    try {
      const servers = await this.listServers();
      const exportData = {
        version: "1.0",
        exportDate: new Date().toISOString(),
        servers: servers.map(server => ({
          ...server,
          // Remove runtime fields
          status: undefined,
          pid: undefined,
          lastStarted: undefined,
          lastStopped: undefined
        }))
      };

      return JSON.stringify(exportData, null, 2);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to export configurations: ${message}`);
      throw new CLIError(
        `Failed to export configurations: ${message}`,
        "MCP_EXPORT_FAILED"
      );
    }
  }

  /**
   * Import MCP server configurations
   */
  async importConfigurations(configData) {
    try {
      const importData = JSON.parse(configData);
      const importedServers = [];

      if (!importData.servers || !Array.isArray(importData.servers)) {
        throw new CLIError("Invalid configuration data format", "MCP_INVALID_CONFIG_DATA");
      }

      for (const serverConfig of importData.servers) {
        try {
          await this.registerServer(serverConfig);
          importedServers.push(serverConfig.id);
        } catch (error) {
          this.logger.warn(`Failed to import server '${serverConfig.id}': ${error.message}`);
        }
      }

      this.logger.info(`Imported ${importedServers.length} MCP server configurations`);
      return importedServers;
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to import configurations: ${message}`);
      throw new CLIError(
        `Failed to import configurations: ${message}`,
        "MCP_IMPORT_FAILED"
      );
    }
  }

  /**
   * Ensure required directories exist
   */
  async ensureDirectories() {
    const dirs = [
      this.configDirectory,
      this.serversDirectory,
      this.logsDirectory
    ];

    for (const dir of dirs) {
      try {
        await fs.mkdir(dir, { recursive: true });
      } catch (error) {
        if (error.code !== "EEXIST") {
          throw error;
        }
      }
    }
  }

  /**
   * Wait for server to be ready
   */
  async waitForServerReady(serverId, timeout) {
    const startTime = Date.now();

    while (Date.now() - startTime < timeout) {
      try {
        const health = await this.getServerHealth(serverId);
        if (health.healthy) {
          return;
        }
      } catch {
        // Continue waiting
      }

      // Wait 1 second before checking again
      await new Promise(resolve => setTimeout(resolve, 1000));
    }

    throw new CLIError(
      `MCP server '${serverId}' did not become ready within ${timeout}ms`,
      "MCP_SERVER_READY_TIMEOUT"
    );
  }

  /**
   * Perform health check on a server
   */
  async performHealthCheck(serverId) {
    try {
      const server = await this.getServer(serverId);
      if (!server) {
        return { healthy: false, error: "Server not found" };
      }

      if (server.status !== "running") {
        return { healthy: false, error: "Server not running" };
      }

      const healthCheck = server.healthCheck;
      if (!healthCheck?.enabled) {
        return { healthy: true };
      }

      const startTime = Date.now();

      // Perform health check based on configuration
      if (healthCheck.endpoint) {
        // HTTP health check
        const url = new URL(healthCheck.endpoint);
        const client = url.protocol === "https:" ? https : http;

        return new Promise((resolve) => {
          const req = client.request({
            hostname: url.hostname,
            port: url.port || (url.protocol === "https:" ? 443 : 80),
            path: url.pathname,
            method: "GET",
            timeout: healthCheck.timeout || 5000
          });

          req.on("response", (res) => {
            const responseTime = Date.now() - startTime;
            resolve({
              healthy: res.statusCode >= 200 && res.statusCode < 300,
              responseTime
            });
          });

          req.on("error", (error) => {
            const responseTime = Date.now() - startTime;
            resolve({
              healthy: false,
              responseTime,
              error: error.message
            });
          });

          req.end();
        });
      } else if (healthCheck.command) {
        // Command-based health check
        return new Promise((resolve) => {
          const child = spawn(healthCheck.command, [], {
            cwd: server.cwd || process.cwd(),
            timeout: healthCheck.timeout || 5000
          });

          const startTime = Date.now();

          child.on("exit", (code) => {
            const responseTime = Date.now() - startTime;
            resolve({
              healthy: code === 0,
              responseTime
            });
          });

          child.on("error", (error) => {
            const responseTime = Date.now() - startTime;
            resolve({
              healthy: false,
              responseTime,
              error: error.message
            });
          });
        });
      } else {
        // Basic process check
        const isRunning = await this.isServerRunning(serverId);
        return {
          healthy: isRunning,
          responseTime: 0
        };
      }
    } catch (error) {
      return {
        healthy: false,
        error: error instanceof Error ? error.message : "Unknown error"
      };
    }
  }

  /**
   * Set up health check monitoring for a server
   */
  setupHealthCheck(serverId) {
    const timer = setInterval(async () => {
      try {
        const result = await this.performHealthCheck(serverId);

        if (!result.healthy) {
          const serverInfo = this.runningServers.get(serverId);
          if (serverInfo) {
            serverInfo.healthCheckFailures++;

            this.logger.warn(`MCP server '${serverId}' health check failed (${serverInfo.healthCheckFailures}/${serverInfo.config.healthCheck?.retries || 3})`);

            // Stop server if too many failures
            if (serverInfo.healthCheckFailures >= (serverInfo.config.healthCheck?.retries || 3)) {
              this.logger.error(`MCP server '${serverId}' health check failed too many times, stopping server`);
              await this.stopServer(serverId);
            }
          }
        } else {
          // Reset failure count on successful health check
          const serverInfo = this.runningServers.get(serverId);
          if (serverInfo) {
            serverInfo.healthCheckFailures = 0;
          }
        }
      } catch (error) {
        this.logger.error(`Health check error for '${serverId}': ${error.message}`);
      }
    }, this.healthCheckInterval);

    this.healthCheckTimers.set(serverId, timer);
  }

  /**
   * Remove health check monitoring for a server
   */
  removeHealthCheck(serverId) {
    const timer = this.healthCheckTimers.get(serverId);
    if (timer) {
      clearInterval(timer);
      this.healthCheckTimers.delete(serverId);
    }
  }

  /**
   * Set up log capture for a server process
   */
  setupLogCapture(serverId, child) {
    const logPath = path.join(this.logsDirectory, `${serverId}.log`);
    const logStream = fs.createWriteStream(logPath, { flags: "a" });

    if (child.stdout) {
      child.stdout.on("data", (data) => {
        const logLine = `[${new Date().toISOString()}] [STDOUT] ${data.toString().trim()}\n`;
        logStream.write(logLine);
      });
    }

    if (child.stderr) {
      child.stderr.on("data", (data) => {
        const logLine = `[${new Date().toISOString()}] [STDERR] ${data.toString().trim()}\n`;
        logStream.write(logLine);
      });
    }

    // Close log stream when process exits
    child.on("exit", () => {
      setTimeout(() => {
        logStream.end();
      }, 1000);
    });
  }

  /**
   * Generate a server ID from package name
   */
  generateServerId(packageName) {
    return packageName.toLowerCase().replace(/[^a-z0-9]/g, "-").replace(/-+/g, "-").replace(/^-|-$/g, "");
  }

  /**
   * Initialize the service
   */
  async onInit() {
    await this.ensureDirectories();

    // Load existing server configurations
    try {
      const serverFiles = await fs.readdir(this.serversDirectory);
      for (const file of serverFiles) {
        if (file.endsWith(".json")) {
          const serverId = file.replace(".json", "");
          await this.getServer(serverId); // This loads into memory cache
        }
      }

      this.logger.info(`MCPService initialized with ${this.serverConfigs.size} servers`);
    } catch (error) {
      this.logger.warn(`Failed to load existing server configurations: ${error.message}`);
    }

    // Start auto-start servers
    if (this.serverConfigs.size > 0) {
      setTimeout(async () => {
        try {
          await this.startAutoStartServers();
        } catch (error) {
          this.logger.error(`Failed to start auto-start servers: ${error.message}`);
        }
      }, 1000);
    }
  }

  /**
   * Cleanup when service is destroyed
   */
  async onDestroy() {
    // Stop all running servers
    await this.stopAllServers(true);

    // Clear all timers
    for (const timer of this.healthCheckTimers.values()) {
      clearInterval(timer);
    }
    this.healthCheckTimers.clear();

    this.logger.info("MCPService destroyed");
  }
}


