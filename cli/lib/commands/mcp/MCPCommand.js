import { BaseCommand } from "../../core/interfaces/BaseCommand.js";
import { ServiceRegistry } from "../../core/registry/ServiceRegistry.js";
import { ILogger } from "../../core/logging/ILogger.js";
import { LoggerFactory } from "../../core/logging/LoggerFactory.js";
import { CLIError } from "../../types/CLITypes.js";
import chalk from "chalk";

/**
 * MCPCommand - MCF CLI MCP Server Management
 * Provides subcommands for managing MCP servers
 */
export class MCPCommand extends BaseCommand {
  constructor(serviceRegistry) {
    super();
    this.serviceRegistry = serviceRegistry;
    this.logger = LoggerFactory.getLogger("MCPCommand");
    this.mcpService = null;
  }

  static get metadata() {
    return {
      name: "MCPCommand",
      description: "Manage MCP servers and their lifecycle",
      category: "mcp",
      version: "1.0.0",
      dependencies: {
        services: ["IMCPService"],
        commands: [],
        external: []
      }
    };
  }

  async initialize() {
    try {
      this.mcpService = this.serviceRegistry.getService("IMCPService");
      this.logger.debug("MCPCommand initialized with MCP service");
    } catch (error) {
      this.logger.error("Failed to initialize MCPCommand", error);
      throw new CLIError(
        "Failed to initialize MCP services",
        "MCP_COMMAND_INIT_FAILED"
      );
    }
  }

  async execute(args = []) {
    await this.initialize();

    if (args.length === 0) {
      return this.showHelp();
    }

    const [subcommand, ...subArgs] = args;

    switch (subcommand.toLowerCase()) {
      case "list":
      case "ls":
        return await this.listServers();
      case "show":
      case "info":
        return await this.showServer(subArgs[0]);
      case "start":
        return await this.startServer(subArgs[0]);
      case "stop":
        return await this.stopServer(subArgs);
      case "restart":
        return await this.restartServer(subArgs[0]);
      case "status":
        return await this.showStatus();
      case "health":
        return await this.showHealth(subArgs[0]);
      case "install":
        return await this.installServer(subArgs[0]);
      case "remove":
      case "uninstall":
        return await this.removeServer(subArgs);
      case "logs":
        return await this.showLogs(subArgs);
      case "config":
        return await this.handleConfigCommand(subArgs);
      default:
        console.log(chalk.red(`Unknown MCP subcommand: ${subcommand}`));
        return this.showHelp();
    }
  }

  async listServers() {
    try {
      const servers = await this.mcpService.listServers();

      if (servers.length === 0) {
        console.log(chalk.yellow("No MCP servers registered."));
        console.log("Install one with: mcf mcp install <package-name>");
        return;
      }

      console.log(chalk.blue("MCP Servers:"));
      console.log();

      for (const server of servers) {
        const statusColor = this.getStatusColor(server.status);
        const autoStart = server.autoStart ? chalk.green(" (auto)") : "";
        console.log(`  ${chalk.cyan(server.id)}${autoStart}`);
        console.log(`    Name: ${server.name}`);
        console.log(`    Status: ${statusColor(server.status)}`);
        console.log(`    Command: ${chalk.gray(server.command)}`);
        if (server.description) {
          console.log(`    Description: ${server.description}`);
        }
        if (server.pid) {
          console.log(`    PID: ${server.pid}`);
        }
        console.log();
      }

      console.log(`Total: ${servers.length} server(s)`);
    } catch (error) {
      console.error(chalk.red(`Failed to list MCP servers: ${error.message}`));
      throw error;
    }
  }

  async showServer(serverId) {
    if (!serverId) {
      console.log(chalk.red("Server ID is required"));
      console.log("Usage: mcf mcp show <server-id>");
      return;
    }

    try {
      const server = await this.mcpService.getServer(serverId);

      if (!server) {
        console.log(chalk.red(`MCP server '${serverId}' not found`));
        return;
      }

      console.log(chalk.blue(`MCP Server: ${server.id}`));
      console.log(chalk.gray("─".repeat(50)));
      console.log(`Name: ${chalk.cyan(server.name)}`);
      console.log(`Status: ${this.getStatusColor(server.status)(server.status)}`);
      console.log(`Command: ${chalk.gray(server.command)}`);

      if (server.args && server.args.length > 0) {
        console.log(`Arguments: ${chalk.gray(server.args.join(" "))}`);
      }

      if (server.cwd) {
        console.log(`Working Directory: ${chalk.gray(server.cwd)}`);
      }

      if (server.port) {
        console.log(`Port: ${server.port}`);
      }

      if (server.description) {
        console.log(`Description: ${server.description}`);
      }

      if (server.version) {
        console.log(`Version: ${server.version}`);
      }

      if (server.pid) {
        console.log(`Process ID: ${server.pid}`);
      }

      if (server.lastStarted) {
        console.log(`Last Started: ${new Date(server.lastStarted).toLocaleString()}`);
      }

      if (server.lastStopped) {
        console.log(`Last Stopped: ${new Date(server.lastStopped).toLocaleString()}`);
      }

      console.log();
      console.log(chalk.blue("Configuration:"));

      if (server.autoStart) {
        console.log(`  Auto Start: ${chalk.green("enabled")}`);
      }

      if (server.restartPolicy) {
        console.log(`  Restart Policy: ${server.restartPolicy}`);
      }

      if (server.healthCheck?.enabled) {
        console.log(`  Health Check: ${chalk.green("enabled")}`);
        console.log(`    Interval: ${server.healthCheck.interval}ms`);
        console.log(`    Timeout: ${server.healthCheck.timeout}ms`);
        console.log(`    Retries: ${server.healthCheck.retries}`);
      }

      if (server.env && Object.keys(server.env).length > 0) {
        console.log();
        console.log(chalk.blue("Environment Variables:"));
        Object.entries(server.env).forEach(([key, value]) => {
          console.log(`  ${key}=${value}`);
        });
      }
    } catch (error) {
      console.error(chalk.red(`Failed to show MCP server: ${error.message}`));
      throw error;
    }
  }

  async startServer(serverId) {
    if (!serverId) {
      console.log(chalk.red("Server ID is required"));
      console.log("Usage: mcf mcp start <server-id>");
      return;
    }

    try {
      console.log(chalk.blue(`Starting MCP server '${serverId}'...`));
      const success = await this.mcpService.startServer(serverId);

      if (success) {
        console.log(chalk.green(`MCP server '${serverId}' started successfully`));

        // Wait a moment and show status
        await new Promise(resolve => setTimeout(resolve, 1000));
        await this.showServer(serverId);
      } else {
        console.log(chalk.red(`Failed to start MCP server '${serverId}'`));
      }
    } catch (error) {
      console.error(chalk.red(`Failed to start MCP server: ${error.message}`));
      throw error;
    }
  }

  async stopServer(args) {
    const [serverId, ...flags] = args;
    const force = flags.includes("--force") || flags.includes("-f");

    if (!serverId) {
      console.log(chalk.red("Server ID is required"));
      console.log("Usage: mcf mcp stop <server-id> [--force]");
      return;
    }

    try {
      console.log(chalk.blue(`Stopping MCP server '${serverId}'${force ? " (force)" : ""}...`));
      const success = await this.mcpService.stopServer(serverId, force);

      if (success) {
        console.log(chalk.green(`MCP server '${serverId}' stopped successfully`));
      } else {
        console.log(chalk.yellow(`MCP server '${serverId}' was not running`));
      }
    } catch (error) {
      console.error(chalk.red(`Failed to stop MCP server: ${error.message}`));
      throw error;
    }
  }

  async restartServer(serverId) {
    if (!serverId) {
      console.log(chalk.red("Server ID is required"));
      console.log("Usage: mcf mcp restart <server-id>");
      return;
    }

    try {
      console.log(chalk.blue(`Restarting MCP server '${serverId}'...`));
      const success = await this.mcpService.restartServer(serverId);

      if (success) {
        console.log(chalk.green(`MCP server '${serverId}' restarted successfully`));

        // Wait a moment and show status
        await new Promise(resolve => setTimeout(resolve, 1000));
        await this.showServer(serverId);
      } else {
        console.log(chalk.red(`Failed to restart MCP server '${serverId}'`));
      }
    } catch (error) {
      console.error(chalk.red(`Failed to restart MCP server: ${error.message}`));
      throw error;
    }
  }

  async showStatus() {
    try {
      const stats = await this.mcpService.getServiceStats();
      const healthStatuses = await this.mcpService.getAllServerHealth();

      console.log(chalk.blue("MCP Service Status"));
      console.log(chalk.gray("─".repeat(50)));
      console.log(`Total Servers: ${chalk.cyan(stats.totalServers)}`);
      console.log(`Running: ${chalk.green(stats.runningServers)}`);
      console.log(`Stopped: ${chalk.gray(stats.stoppedServers)}`);
      console.log(`Errors: ${chalk.red(stats.erroredServers)}`);
      console.log(`Auto-start: ${chalk.blue(stats.autoStartServers)}`);
      console.log();

      if (stats.totalServers > 0) {
        console.log(chalk.blue("Server Health:"));
        healthStatuses.forEach(health => {
          const statusColor = health.healthy ? chalk.green : chalk.red;
          const statusText = health.healthy ? "healthy" : "unhealthy";
          console.log(`  ${chalk.cyan(health.serverId)}: ${statusColor(statusText)}`);
        });
        console.log();

        console.log(chalk.blue("Uptime Statistics:"));
        console.log(`  Total Uptime: ${this.formatUptime(stats.totalUptime)}`);
        console.log(`  Average Health Check: ${Math.round(stats.averageHealthCheckTime)}ms`);
      }
    } catch (error) {
      console.error(chalk.red(`Failed to get MCP status: ${error.message}`));
      throw error;
    }
  }

  async showHealth(serverId) {
    if (!serverId) {
      console.log(chalk.red("Server ID is required"));
      console.log("Usage: mcf mcp health <server-id>");
      return;
    }

    try {
      const health = await this.mcpService.getServerHealth(serverId);

      console.log(chalk.blue(`MCP Server Health: ${serverId}`));
      console.log(chalk.gray("─".repeat(50)));

      const statusColor = health.healthy ? chalk.green : chalk.red;
      const statusText = health.healthy ? "HEALTHY" : "UNHEALTHY";
      console.log(`Status: ${statusColor(statusText)}`);

      if (health.responseTime !== undefined) {
        console.log(`Response Time: ${health.responseTime}ms`);
      }

      console.log(`Last Check: ${health.lastCheck.toLocaleString()}`);

      if (health.error) {
        console.log(`Error: ${chalk.red(health.error)}`);
      }

      // Show server details
      console.log();
      const server = await this.mcpService.getServer(serverId);
      if (server) {
        console.log(`Server Status: ${this.getStatusColor(server.status)(server.status)}`);
        if (server.pid) {
          console.log(`Process ID: ${server.pid}`);
        }
      }
    } catch (error) {
      console.error(chalk.red(`Failed to get MCP server health: ${error.message}`));
      throw error;
    }
  }

  async installServer(packageName) {
    if (!packageName) {
      console.log(chalk.red("Package name is required"));
      console.log("Usage: mcf mcp install <package-name>");
      return;
    }

    try {
      console.log(chalk.blue(`Installing MCP server: ${packageName}...`));
      const server = await this.mcpService.installServer(packageName);

      console.log(chalk.green(`MCP server '${server.id}' installed successfully`));
      console.log(`Name: ${server.name}`);
      console.log(`Command: ${server.command}`);
      console.log();
      console.log("You can now:");
      console.log(`  • Start the server: mcf mcp start ${server.id}`);
      console.log(`  • Show details: mcf mcp show ${server.id}`);
      console.log(`  • Check health: mcf mcp health ${server.id}`);
    } catch (error) {
      console.error(chalk.red(`Failed to install MCP server: ${error.message}`));
      throw error;
    }
  }

  async removeServer(args) {
    const [serverId, ...flags] = args;
    const keepConfig = flags.includes("--keep-config");

    if (!serverId) {
      console.log(chalk.red("Server ID is required"));
      console.log("Usage: mcf mcp remove <server-id> [--keep-config]");
      return;
    }

    try {
      console.log(chalk.blue(`Removing MCP server '${serverId}'...`));
      const success = await this.mcpService.removeServer(serverId, keepConfig);

      if (success) {
        if (keepConfig) {
          console.log(chalk.green(`MCP server '${serverId}' removed (config preserved)`));
        } else {
          console.log(chalk.green(`MCP server '${serverId}' and config removed successfully`));
        }
      } else {
        console.log(chalk.red(`Failed to remove MCP server '${serverId}'`));
      }
    } catch (error) {
      console.error(chalk.red(`Failed to remove MCP server: ${error.message}`));
      throw error;
    }
  }

  async showLogs(args) {
    const [serverId, lines = "50"] = args;
    const linesNum = parseInt(lines);

    if (!serverId) {
      console.log(chalk.red("Server ID is required"));
      console.log("Usage: mcf mcp logs <server-id> [lines]");
      return;
    }

    if (isNaN(linesNum) || linesNum < 1) {
      console.log(chalk.red("Lines must be a positive number"));
      return;
    }

    try {
      const logs = await this.mcpService.getServerLogs(serverId, linesNum);

      if (logs.length === 0 || logs[0] === "No logs available") {
        console.log(chalk.yellow(`No logs available for MCP server '${serverId}'`));
        return;
      }

      console.log(chalk.blue(`MCP Server Logs: ${serverId}`));
      console.log(chalk.gray("─".repeat(50)));
      console.log();

      logs.forEach(logLine => {
        console.log(logLine);
      });

      console.log();
      console.log(chalk.gray(`Showing last ${logs.length} log lines`));
    } catch (error) {
      console.error(chalk.red(`Failed to get MCP server logs: ${error.message}`));
      throw error;
    }
  }

  async handleConfigCommand(args) {
    const [configSubcommand, ...configArgs] = args;

    switch (configSubcommand?.toLowerCase()) {
      case "export":
        return await this.exportConfig();
      case "import":
        return await this.importConfig(configArgs[0]);
      default:
        console.log(chalk.red("Config subcommand required"));
        console.log("Available config commands:");
        console.log("  export                 Export server configurations");
        console.log("  import <file>          Import server configurations");
        return;
    }
  }

  async exportConfig() {
    try {
      const configData = await this.mcpService.exportConfigurations();
      console.log(configData);
      console.log();
      console.log(chalk.gray("MCP server configurations exported to stdout"));
    } catch (error) {
      console.error(chalk.red(`Failed to export configurations: ${error.message}`));
      throw error;
    }
  }

  async importConfig(filePath) {
    if (!filePath) {
      console.log(chalk.red("Configuration file path is required"));
      console.log("Usage: mcf mcp config import <file-path>");
      return;
    }

    try {
      // Read config file
      const fs = await import("fs/promises");
      const configData = await fs.readFile(filePath, "utf-8");

      console.log(chalk.blue("Importing MCP server configurations..."));
      const importedServers = await this.mcpService.importConfigurations(configData);

      console.log(chalk.green(`Successfully imported ${importedServers.length} server configurations:`));
      importedServers.forEach(serverId => {
        console.log(`  • ${chalk.cyan(serverId)}`);
      });
    } catch (error) {
      console.error(chalk.red(`Failed to import configurations: ${error.message}`));
      throw error;
    }
  }

  showHelp() {
    console.log(chalk.blue("MCF MCP Command Help"));
    console.log(chalk.gray("─".repeat(50)));
    console.log();
    console.log("Manage MCP (Model Context Protocol) servers");
    console.log();
    console.log(chalk.blue("Server Management:"));
    console.log();
    console.log("  list, ls                    List all MCP servers");
    console.log("  show, info <id>             Show server details");
    console.log("  start <id>                  Start a server");
    console.log("  stop <id> [--force]         Stop a server");
    console.log("  restart <id>                Restart a server");
    console.log();
    console.log(chalk.blue("Monitoring:"));
    console.log();
    console.log("  status                      Show overall service status");
    console.log("  health <id>                 Check server health");
    console.log("  logs <id> [lines]           Show server logs");
    console.log();
    console.log(chalk.blue("Installation:"));
    console.log();
    console.log("  install <package>           Install MCP server package");
    console.log("  remove, uninstall <id>      Remove MCP server");
    console.log();
    console.log(chalk.blue("Configuration:"));
    console.log();
    console.log("  config export               Export server configurations");
    console.log("  config import <file>        Import server configurations");
    console.log();
    console.log(chalk.blue("Examples:"));
    console.log();
    console.log("  mcf mcp list");
    console.log("  mcf mcp install @modelcontextprotocol/server-filesystem");
    console.log("  mcf mcp start filesystem-server");
    console.log("  mcf mcp health filesystem-server");
    console.log("  mcf mcp logs filesystem-server 100");
    console.log("  mcf mcp config export > servers.json");
    console.log();
    return { success: true };
  }

  getStatusColor(status) {
    switch (status) {
      case "running":
        return chalk.green;
      case "stopped":
        return chalk.gray;
      case "starting":
        return chalk.blue;
      case "stopping":
        return chalk.yellow;
      case "error":
        return chalk.red;
      default:
        return chalk.gray;
    }
  }

  formatUptime(milliseconds) {
    const seconds = Math.floor(milliseconds / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (days > 0) {
      return `${days}d ${hours % 24}h ${minutes % 60}m`;
    } else if (hours > 0) {
      return `${hours}h ${minutes % 60}m ${seconds % 60}s`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds % 60}s`;
    } else {
      return `${seconds}s`;
    }
  }

  getMetadata() {
    return MCPCommand.metadata;
  }
}
