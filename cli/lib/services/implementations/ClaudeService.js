import { BaseService } from "../../core/base/BaseService.js";
import { ILogger } from "../../core/logging/ILogger.js";
import { LoggerFactory } from "../../core/logging/LoggerFactory.js";
import { CLIError } from "../../types/CLITypes.js";
import { spawn } from "child_process";
import path from "path";
import fs from "fs/promises";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);

/**
 * Claude service implementation for MCF CLI
 * Handles direct integration with Claude Code CLI
 */
export class ClaudeService extends BaseService {
  constructor(config, logger) {
    super();
    this.config = config || {};
    this.logger = logger || LoggerFactory.getLogger("ClaudeService");

    // Default configuration
    this.defaultExecutablePath = config?.defaultExecutablePath || "claude";
    this.defaultTimeout = config?.defaultTimeout || 300000; // 5 minutes
    this.validateOnStartup = config?.validateOnStartup !== false;
    this.supportedModels = config?.supportedModels || [
      "claude-3-5-sonnet-20241022",
      "claude-3-opus-20240229",
      "claude-3-haiku-20240307"
    ];
    this.defaultEnvironment = config?.defaultEnvironment || {};

    // Track running processes
    this.runningProcesses = new Map();

    this.logger.debug("ClaudeService initialized", {
      executablePath: this.defaultExecutablePath,
      timeout: this.defaultTimeout
    });
  }

  /**
   * Execute Claude Code with specified options
   */
  async runClaude(options = {}) {
    try {
      const startTime = Date.now();
      const mergedOptions = { ...this.getDefaultOptions(), ...options };

      // Validate configuration
      const validation = await this.validateConfiguration(mergedOptions);
      if (!validation.isValid) {
        throw new CLIError(
          `Claude configuration validation failed: ${validation.errors.join(", ")}`,
          "CLAUDE_CONFIG_INVALID",
          { errors: validation.errors, warnings: validation.warnings }
        );
      }

      if (validation.warnings.length > 0) {
        this.logger.warn(`Claude configuration warnings: ${validation.warnings.join(", ")}`);
      }

      // Skip installation check for now to avoid hanging
      // if (!(await this.isInstalled())) {
      //   throw new CLIError(
      //     "Claude Code is not installed or not accessible",
      //     "CLAUDE_NOT_INSTALLED",
      //     { executablePath: this.defaultExecutablePath }
      //   );
      // }

      // Build command arguments
      const args = this.buildClaudeArguments(mergedOptions);

      // Configure environment
      const environment = await this.configureEnvironment(mergedOptions);

      this.logger.info("Starting Claude Code", {
        args: args.join(" "),
        workingDirectory: mergedOptions.workingDirectory || process.cwd(),
        timeout: mergedOptions.timeout
      });

      return new Promise((resolve, reject) => {
        const child = spawn(this.defaultExecutablePath, args, {
          stdio: mergedOptions.interactive !== false ? "inherit" : "pipe",
          env: {
            ...process.env,
            ...environment
          },
          cwd: mergedOptions.workingDirectory || process.cwd(),
          shell: process.platform === "win32"
        });

        // Track the process
        this.runningProcesses.set(child.pid, {
          process: child,
          startTime: new Date(),
          args: args,
          cwd: mergedOptions.workingDirectory || process.cwd()
        });

        let output = "";
        let errorOutput = "";

        // Handle stdout if not in interactive mode
        if (mergedOptions.interactive === false && child.stdout) {
          child.stdout.on("data", (data) => {
            output += data.toString();
          });
        }

        // Handle stderr if not in interactive mode
        if (mergedOptions.interactive === false && child.stderr) {
          child.stderr.on("data", (data) => {
            errorOutput += data.toString();
          });
        }

        // Handle process exit
        child.on("exit", (code, signal) => {
          const executionTime = Date.now() - startTime;
          const exitCode = code || 0;
          const success = exitCode === 0 && !signal;

          // Remove from tracking
          this.runningProcesses.delete(child.pid);

          this.logger.info("Claude Code exited", {
            exitCode,
            signal,
            executionTime,
            success
          });

          resolve({
            exitCode,
            executionTime,
            signal,
            success,
            output: output || undefined,
            errorOutput: errorOutput || undefined
          });
        });

        // Handle process error
        child.on("error", (error) => {
          const executionTime = Date.now() - startTime;

          // Remove from tracking
          this.runningProcesses.delete(child.pid);

          this.logger.error("Failed to start Claude Code", {
            error: error.message,
            command: this.defaultExecutablePath,
            args
          });

          reject(
            new CLIError(
              `Failed to start Claude: ${error.message}`,
              "CLAUDE_START_FAILED",
              { originalError: error, args, executionTime }
            )
          );
        });

        // Handle timeout
        if (mergedOptions.timeout) {
          setTimeout(() => {
            if (!child.killed) {
              this.logger.warn(`Claude execution timeout after ${mergedOptions.timeout}ms`);
              child.kill("SIGTERM");

              // Give it 5 seconds to terminate gracefully
              setTimeout(() => {
                if (!child.killed) {
                  child.kill("SIGKILL");
                }
              }, 5000);
            }
          }, mergedOptions.timeout);
        }

        // Handle process signals for graceful shutdown
        if (mergedOptions.interactive !== false) {
          const handleSignal = (signal) => {
            this.logger.info(`Received ${signal}, terminating Claude Code`);
            if (!child.killed) {
              child.kill(signal);
            }
          };

          process.on("SIGINT", handleSignal);
          process.on("SIGTERM", handleSignal);
        }
      });
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to run Claude: ${message}`);
      throw new CLIError(
        `Failed to run Claude: ${message}`,
        "CLAUDE_EXECUTION_FAILED"
      );
    }
  }

  /**
   * Get Claude Code version information
   */
  async getVersion() {
    try {
      const result = await this.runClaude({
        additionalArgs: ["--version"],
        interactive: false,
        timeout: 10000
      });

      if (!result.success) {
        return {
          version: "unknown",
          installed: false
        };
      }

      const version = result.output?.trim() || "unknown";

      return {
        version,
        installed: true,
        executablePath: await this.getExecutablePath(),
        supportedModels: this.supportedModels
      };
    } catch (error) {
      this.logger.debug(`Failed to get Claude version: ${error.message}`);
      return {
        version: "unknown",
        installed: false
      };
    }
  }

  /**
   * Check if Claude Code is installed and accessible
   */
  async isInstalled() {
    try {
      const executablePath = await this.getExecutablePath();
      if (!executablePath) {
        return false;
      }

      // Try to run --help to check if it's working
      const result = await this.runClaude({
        additionalArgs: ["--help"],
        interactive: false,
        timeout: 5000
      });

      return result.success || result.exitCode !== 127; // 127 = command not found
    } catch {
      return false;
    }
  }

  /**
   * Get the path to Claude executable
   */
  async getExecutablePath() {
    // If we have a configured path, use it
    if (this.config?.defaultExecutablePath) {
      try {
        await fs.access(this.config.defaultExecutablePath);
        return this.config.defaultExecutablePath;
      } catch {
        // Fall through to search
      }
    }

    // Search for claude in PATH
    const paths = process.env.PATH?.split(path.delimiter) || [];
    for (const searchPath of paths) {
      const executablePath = path.join(searchPath, process.platform === "win32" ? "claude.cmd" : "claude");
      try {
        await fs.access(executablePath);
        return executablePath;
      } catch {
        // Continue searching
      }
    }

    return null;
  }

  /**
   * Validate Claude configuration
   */
  async validateConfiguration(config) {
    const errors = [];
    const warnings = [];
    const suggestions = [];

    // Validate working directory
    if (config.workingDirectory) {
      try {
        await fs.access(config.workingDirectory);
      } catch {
        errors.push(`Working directory does not exist: ${config.workingDirectory}`);
      }
    }

    // Validate model names
    if (config.model && !this.supportedModels.includes(config.model)) {
      warnings.push(`Unknown model: ${config.model}`);
      suggestions.push(`Supported models: ${this.supportedModels.join(", ")}`);
    }

    if (config.smallFastModel && !this.supportedModels.includes(config.smallFastModel)) {
      warnings.push(`Unknown small fast model: ${config.smallFastModel}`);
      suggestions.push(`Supported models: ${this.supportedModels.join(", ")}`);
    }

    // Validate timeout
    if (config.timeout && (config.timeout < 1000 || config.timeout > 3600000)) {
      warnings.push("Timeout should be between 1 second and 1 hour");
      suggestions.push("Consider timeout between 30000ms (30s) and 300000ms (5min)");
    }

    return {
      isValid: errors.length === 0,
      errors,
      warnings,
      suggestions
    };
  }

  /**
   * Build Claude command arguments from options
   */
  buildClaudeArguments(options) {
    const args = [];

    // Add flags
    if (options.dangerousSkip) {
      args.push("--dangerously-skip-permissions");
    }

    if (options.debug) {
      args.push("--debug");
    }

    // Add project
    if (options.project) {
      args.push("--project", options.project);
    }

    // Add additional arguments
    if (options.additionalArgs) {
      args.push(...options.additionalArgs);
    }

    return args;
  }

  /**
   * Configure Claude environment variables
   */
  async configureEnvironment(options) {
    const env = { ...this.defaultEnvironment };

    // Set base URL
    if (options.anthropicBaseUrl) {
      env.ANTHROPIC_BASE_URL = options.anthropicBaseUrl;
    }

    // Set auth token
    if (options.authToken) {
      env.ANTHROPIC_AUTH_TOKEN = options.authToken;
    }

    // Set config directory
    if (options.configDirectory) {
      env.CLAUDE_CONFIG_DIR = options.configDirectory;
    }

    // Set model
    if (options.model) {
      env.ANTHROPIC_MODEL = options.model;
    }

    // Set small fast model
    if (options.smallFastModel) {
      env.ANTHROPIC_SMALL_FAST_MODEL = options.smallFastModel;
    }

    // Add custom environment variables
    if (options.environment) {
      Object.assign(env, options.environment);
    }

    return env;
  }

  /**
   * Get default Claude execution options
   */
  getDefaultOptions() {
    return {
      workingDirectory: process.cwd(),
      dangerousSkip: false,
      debug: false,
      interactive: true,
      timeout: this.defaultTimeout,
      additionalArgs: []
    };
  }

  /**
   * Kill a running Claude process
   */
  async killClaude(signal = "SIGTERM") {
    try {
      if (this.runningProcesses.size === 0) {
        this.logger.debug("No running Claude processes to kill");
        return false;
      }

      let killed = false;
      for (const [pid, processInfo] of this.runningProcesses) {
        try {
          processInfo.process.kill(signal);
          killed = true;
          this.logger.info(`Sent ${signal} to Claude process ${pid}`);
        } catch (error) {
          this.logger.warn(`Failed to kill Claude process ${pid}: ${error.message}`);
        }
      }

      // Clear tracking
      this.runningProcesses.clear();

      return killed;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to kill Claude processes: ${message}`);
      throw new CLIError(
        `Failed to kill Claude processes: ${message}`,
        "CLAUDE_KILL_FAILED"
      );
    }
  }

  /**
   * Check if Claude process is currently running
   */
  async isRunning() {
    return this.runningProcesses.size > 0;
  }

  /**
   * Get Claude process information
   */
  async getProcessInfo() {
    if (this.runningProcesses.size === 0) {
      return null;
    }

    // Return info for the first (most recent) process
    const [pid, processInfo] = this.runningProcesses.entries().next().value;

    return {
      pid,
      running: !processInfo.process.killed,
      startTime: processInfo.startTime,
      args: processInfo.args,
      cwd: processInfo.cwd
    };
  }

  /**
   * Initialize the service
   */
  async onInit() {
    if (this.validateOnStartup) {
      const installed = await this.isInstalled();
      if (!installed) {
        this.logger.warn("Claude Code is not installed or not accessible");
      } else {
        const versionInfo = await this.getVersion();
        this.logger.info("Claude Code integration ready", {
          installed: true,
          version: versionInfo.version,
          executablePath: versionInfo.executablePath
        });
      }
    }
  }
}
