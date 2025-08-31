#!/usr/bin/env node

import { program } from "commander";
import chalk from "chalk";
import path from "path";
import { fileURLToPath } from "url";
import pkg from "../package.json" with { type: "json" };

// Get directory equivalent for ES modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Import ServiceRegistry and core services
import { ServiceRegistry } from "../lib/core/registry/ServiceRegistry.js";
import { LoggerFactory } from "../lib/core/logging/LoggerFactory.js";
import { CommandRegistry } from "../lib/core/registry/CommandRegistry.js";
import { ConfigurationService } from "../lib/services/implementations/ConfigurationService.js";
import { FileSystemService } from "../lib/services/implementations/FileSystemService.js";
import { ProjectService } from "../lib/services/implementations/ProjectService.js";
import { ClaudeService } from "../lib/services/implementations/ClaudeService.js";
import { MCPService } from "../lib/services/implementations/MCPService.js";

// Import command modules (maintaining backward compatibility)
import install from "../lib/install.js";
import setup from "../lib/setup.js";
import run from "../lib/run.js";
import templates from "../lib/templates.js";
import status from "../lib/status.js";

// Initialize ServiceRegistry
const serviceRegistry = ServiceRegistry.getInstance();

// Register core services
const loggerFactory = new LoggerFactory();
serviceRegistry.registerService("LoggerFactory", () => loggerFactory);
serviceRegistry.registerService("ILogger", () => loggerFactory.getLogger("default"));

const commandRegistry = new CommandRegistry();
serviceRegistry.registerService("CommandRegistry", () => commandRegistry);

// Register ConfigurationService
const configService = new ConfigurationService({
  configDirectory: path.join(__dirname, "..", ".mcf"),
  defaultProfileName: "default",
  validateProfiles: true
});
serviceRegistry.registerService("IConfigurationService", () => configService);

// Register FileSystemService
const fileSystemService = new FileSystemService({
  baseDirectory: process.cwd(),
  createParentDirectories: true,
  defaultEncoding: "utf-8",
  maxFileSize: 10 * 1024 * 1024, // 10MB
  permissions: {
    defaultFileMode: 0o644,
    defaultDirectoryMode: 0o755
  }
});
serviceRegistry.registerService("IFileSystemService", () => fileSystemService);

// Register ProjectService
const projectService = new ProjectService({
  defaultProjectPath: process.cwd(),
  workspacePath: path.join(__dirname, "..", ".mcf", "workspaces"),
  maxDiscoveryDepth: 3,
  autoDiscover: true
}, null, fileSystemService);
serviceRegistry.registerService("IProjectService", () => projectService);

// Register ClaudeService
const claudeService = new ClaudeService({
  defaultExecutablePath: "claude",
  defaultTimeout: 300000, // 5 minutes
  validateOnStartup: true,
  supportedModels: [
    "claude-3-5-sonnet-20241022",
    "claude-3-opus-20240229",
    "claude-3-haiku-20240307"
  ]
});
serviceRegistry.registerService("IClaudeService", () => claudeService);

// Register MCPService
const mcpService = new MCPService({
  configDirectory: path.join(__dirname, "..", ".mcf", "mcp"),
  defaultTimeout: 30000,
  healthCheckInterval: 30000,
  enableHealthChecks: true
});
serviceRegistry.registerService("IMCPService", () => mcpService);

// Configure program with enhanced capabilities
program
  .name("mcf")
  .description(
    "MCF (My Claude Flow) CLI - Installation, configuration and setup tool",
  )
  .version(pkg.version);

// Install command (enhanced)
program
  .command("install")
  .description("Install MCF framework")
  .option("-y, --yes", "Skip interactive prompts and proceed automatically")
  .option("-p, --profile <profile>", "Specify installation profile")
  .action(async (options) => {
    const logger = loggerFactory.getLogger("install");
    try {
      await install(options);
    } catch (error) {
      logger.error("Installation failed", error);
      process.exit(1);
    }
  });

// Setup command (enhanced)
program
  .command("setup")
  .description("Configure MCF after installation")
  .option("-p, --profile <profile>", "Specify setup profile")
  .action(async (options) => {
    const logger = loggerFactory.getLogger("setup");
    try {
      await setup(options);
    } catch (error) {
      logger.error("Setup failed", error);
      process.exit(1);
    }
  });

// Run command (enhanced with Claude integration)
program
  .command("run")
  .description("Execute Claude Code with configuration and flags")
  .option("-d, --debug", "Enable debug mode")
  .option("-c, --config <profile>", "Use specific configuration profile")
  .option("-p, --project <name>", "Set Claude project name")
  .option("-w, --working-directory <path>", "Set working directory")
  .option("-t, --timeout <ms>", "Set execution timeout (milliseconds)")
  .option("--no-interactive", "Run in non-interactive mode")
  .option("--dangerous-skip", "Skip permission checks (dangerous)")
  .allowUnknownOption() // Support -- pass-through for additional args
  .action(async (options, command) => {
    try {
      const { RunCommand } = await import("../lib/commands/run/RunCommand.js");
      const runCommand = new RunCommand(serviceRegistry);

      // Get all arguments after "run"
      const runArgsIndex = process.argv.indexOf("run") + 1;
      const allArgs = process.argv.slice(runArgsIndex);

      await runCommand.execute(allArgs);
    } catch (error) {
      console.error(chalk.red(`Run command failed: ${error.message}`));
      process.exit(1);
    }
  });

// Templates command (enhanced)
program
  .command("templates")
  .alias("t")
  .description("Manage MCF templates")
  .option("-p, --profile <profile>", "Specify template profile")
  .argument("[action]", "Action to perform (list, init, info)")
  .argument("[name]", "Template name")
  .action(async (action, name, options) => {
    const logger = loggerFactory.getLogger("templates");
    try {
      await templates(action, name, options);
    } catch (error) {
      logger.error("Templates operation failed", error);
      process.exit(1);
    }
  });

// Status command (enhanced)
program
  .command("status")
  .description("Check MCF installation status")
  .option("-p, --profile <profile>", "Specify status profile")
  .action(async (options) => {
    const logger = loggerFactory.getLogger("status");
    try {
      await status(options);
    } catch (error) {
      logger.error("Status check failed", error);
      process.exit(1);
    }
  });

// Config command (new)
program
  .command("config")
  .description("Manage MCF configuration profiles")
  .argument("[subcommand]", "Config subcommand")
  .argument("[args...]", "Subcommand arguments")
  .action(async (subcommand, args, options) => {
    try {
      const { ConfigCommand } = await import("../lib/commands/config/ConfigCommand.js");
      const configCommand = new ConfigCommand(serviceRegistry);
      const allArgs = subcommand ? [subcommand, ...args] : [];
      await configCommand.execute(allArgs);
    } catch (error) {
      console.error(chalk.red(`Config command failed: ${error.message}`));
      process.exit(1);
    }
  });

// Project command (new)
program
  .command("project")
  .description("Manage MCF projects and workspaces")
  .argument("[subcommand]", "Project subcommand")
  .argument("[args...]", "Subcommand arguments")
  .action(async (subcommand, args, options) => {
    try {
      const { ProjectCommand } = await import("../lib/commands/project/ProjectCommand.js");
      const projectCommand = new ProjectCommand(serviceRegistry);
      const allArgs = subcommand ? [subcommand, ...args] : [];
      await projectCommand.execute(allArgs);
    } catch (error) {
      console.error(chalk.red(`Project command failed: ${error.message}`));
      process.exit(1);
    }
  });

// MCP command (new)
program
  .command("mcp")
  .description("Manage MCP servers and their lifecycle")
  .argument("[subcommand]", "MCP subcommand")
  .argument("[args...]", "Subcommand arguments")
  .action(async (subcommand, args, options) => {
    try {
      const { MCPCommand } = await import("../lib/commands/mcp/MCPCommand.js");
      const mcpCommand = new MCPCommand(serviceRegistry);
      const allArgs = subcommand ? [subcommand, ...args] : [];
      await mcpCommand.execute(allArgs);
    } catch (error) {
      console.error(chalk.red(`MCP command failed: ${error.message}`));
      process.exit(1);
    }
  });

// Enhanced error handling
program.configureOutput({
  writeErr: (str) => process.stderr.write(chalk.red(str)),
});

// Parse arguments
program.parse();

// Show help if no command provided
if (!process.argv.slice(2).length) {
  program.outputHelp();
}

// Export ServiceRegistry for potential programmatic use
export { serviceRegistry };
