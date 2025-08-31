#!/usr/bin/env node

import { spawn } from "child_process";
import path from "path";
import { fileURLToPath } from "url";
import fs from "fs/promises";
import os from "os";

// Get directory equivalent for ES modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Package info
const pkg = {
  name: "@pc-style/mcf-cli",
  version: "1.0.0",
  description: "MCF (My Claude Flow) CLI - Installation, configuration and setup tool"
};

// ============================================================================
// UTILITY FUNCTIONS (replacing chalk and commander)
// ============================================================================

// Simple color functions (replacing chalk)
const colors = {
  red: (text) => `\x1b[31m${text}\x1b[0m`,
  green: (text) => `\x1b[32m${text}\x1b[0m`,
  yellow: (text) => `\x1b[33m${text}\x1b[0m`,
  blue: (text) => `\x1b[34m${text}\x1b[0m`,
  cyan: (text) => `\x1b[36m${text}\x1b[0m`,
  gray: (text) => `\x1b[90m${text}\x1b[0m`,
  bold: (text) => `\x1b[1m${text}\x1b[0m`
};

// Simple argument parser (replacing commander)
function parseArgs() {
  const args = process.argv.slice(2);
  const result = {
    command: null,
    subcommand: null,
    options: {},
    args: []
  };

  if (args.length === 0) {
    return result;
  }

  // Handle help and version
  if (args.includes('--help') || args.includes('-h')) {
    result.command = 'help';
    return result;
  }

  if (args.includes('--version') || args.includes('-V')) {
    result.command = 'version';
    return result;
  }

  // Get main command (first non-option argument)
  let commandIndex = -1;
  for (let i = 0; i < args.length; i++) {
    if (!args[i].startsWith('-') && (i === 0 || !args[i-1].startsWith('-') || args[i-1] === '--help' || args[i-1] === '-h' || args[i-1] === '--version' || args[i-1] === '-V')) {
      commandIndex = i;
      break;
    }
    // Skip option values
    if (args[i] === '--config' || args[i] === '-c' || args[i] === '--project' || args[i] === '-p' || args[i] === '--working-directory' || args[i] === '-w') {
      i++; // Skip the next argument (option value)
    }
  }
  
  if (commandIndex >= 0) {
    result.command = args[commandIndex];
    
    // For config command, get subcommand
    if (result.command === 'config' && commandIndex + 1 < args.length && !args[commandIndex + 1].startsWith('-')) {
      result.subcommand = args[commandIndex + 1];
      result.args = args.slice(commandIndex + 2);
    } else {
      result.args = args.slice(commandIndex + 1);
    }
  }

  // Parse options
  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg.startsWith('-')) {
      if (arg === '--force' || arg === '-f') {
        result.options.force = true;
      } else if (arg === '--debug' || arg === '-d') {
        result.options.debug = true;
      } else if (arg === '--config' || arg === '-c') {
        result.options.config = args[i + 1];
        i++; // Skip next arg
      } else if (arg === '--project' || arg === '-p') {
        result.options.project = args[i + 1];
        i++; // Skip next arg
      } else if (arg === '--working-directory' || arg === '-w') {
        result.options.workingDirectory = args[i + 1];
        i++; // Skip next arg
      } else if (arg === '--dangerous-skip') {
        result.options.dangerousSkip = true;
      } else if (arg === '--no-interactive') {
        result.options.interactive = false;
      } else if (arg === '--') {
        // Everything after -- goes to pass-through
        result.options.passThroughArgs = args.slice(i + 1);
        break;
      }
    }
  }

  return result;
}

// ============================================================================
// CORE CLASSES
// ============================================================================

/**
 * Base error class for CLI operations
 */
class CLIError extends Error {
  constructor(message, code, details) {
    super(message);
    this.code = code;
    this.details = details;
    this.name = "CLIError";
  }
}

/**
 * Simple logger implementation
 */
class Logger {
  constructor(name) {
    this.name = name;
  }

  info(message, data) {
    console.log(`[INFO] ${message}`, data ? JSON.stringify(data) : "");
  }

  warn(message, data) {
    console.log(colors.yellow(`[WARN] ${message}`), data ? JSON.stringify(data) : "");
  }

  error(message, data) {
    console.log(colors.red(`[ERROR] ${message}`), data ? JSON.stringify(data) : "");
  }

  debug(message, data) {
    if (process.env.DEBUG) {
      console.log(colors.gray(`[DEBUG] ${message}`), data ? JSON.stringify(data) : "");
    }
  }
}

/**
 * Configuration service for profile management
 */
class ConfigurationService {
  constructor() {
    this.logger = new Logger("ConfigurationService");
    // Look for profiles in $HOME/.mcf/profiles
    this.profilesDir = path.join(os.homedir(), ".mcf", "profiles");
  }

  async ensureProfilesDirectory() {
    try {
      await fs.mkdir(this.profilesDir, { recursive: true });
    } catch (error) {
      if (error.code !== "EEXIST") {
        throw error;
      }
    }
  }

  getProfilePath(profileId) {
    return path.join(this.profilesDir, `${profileId}.json`);
  }

  async loadProfile(profileId) {
    try {
      const profilePath = this.getProfilePath(profileId);
      const profileData = await fs.readFile(profilePath, "utf-8");
      return JSON.parse(profileData);
    } catch {
      return null;
    }
  }

  async saveProfile(profile) {
    await this.ensureProfilesDirectory();
    const profilePath = this.getProfilePath(profile.id);
    await fs.writeFile(profilePath, JSON.stringify(profile, null, 2), "utf-8");
    this.logger.info(`Profile '${profile.id}' saved to ${profilePath}`);
  }

  async listProfiles() {
    try {
      await this.ensureProfilesDirectory();
      const entries = await fs.readdir(this.profilesDir);
      return entries
        .filter(entry => entry.endsWith(".json"))
        .map(entry => entry.replace(".json", ""));
    } catch {
      return [];
    }
  }

  async deleteProfile(profileId) {
    try {
      const profilePath = this.getProfilePath(profileId);
      await fs.unlink(profilePath);
      return true;
    } catch {
      return false;
    }
  }

  async profileExists(profileId) {
    try {
      const profilePath = this.getProfilePath(profileId);
      await fs.access(profilePath);
      return true;
    } catch {
      return false;
    }
  }

  generateProfileId(name) {
    return name.toLowerCase().replace(/[^a-z0-9]/g, "-").replace(/-+/g, "-").replace(/^-|-$/g, "");
  }

  async createProfile(name, environment = "development") {
    const profileId = this.generateProfileId(name);
    const profile = {
      id: profileId,
      name,
      description: `Profile for ${environment} environment`,
      environment,
      config: {
        timeout: 30000,
        maxRetries: 3,
        logLevel: "info"
      },
      version: "1.0.0",
      lastUpdated: new Date()
    };

    await this.saveProfile(profile);
    return profile;
  }
}

// ============================================================================
// COMMAND IMPLEMENTATIONS
// ============================================================================

/**
 * Show help information
 */
function showHelp() {
  console.log(`Usage: mcf [options] [command]`);
  console.log();
  console.log(pkg.description);
  console.log();
  console.log("Options:");
  console.log("  -V, --version                   output the version number");
  console.log("  -h, --help                      display help for command");
  console.log();
  console.log("Commands:");
  console.log("  config [subcommand] [args...]   Manage MCF configuration profiles");
  console.log("  run [options]                   Execute Claude Code with configuration and flags");
  console.log("  install                         Install MCF CLI to ~/.local/bin/mcf");
  console.log("  status                          Check MCF installation status");
  console.log("  help [command]                  display help for command");
  console.log();
  console.log("Examples:");
  console.log("  mcf config list");
  console.log("  mcf config create myprofile");
  console.log("  mcf run --config myprofile");
  console.log("  mcf install");
}

/**
 * Show version information
 */
function showVersion() {
  console.log(pkg.version);
}

/**
 * Configuration command implementation
 */
async function configCommand(subcommand, args, options) {
  const configService = new ConfigurationService();

  if (!subcommand) {
    console.log(colors.blue("MCF Config Command Help"));
    console.log("─".repeat(50));
    console.log();
    console.log("Manage MCF configuration profiles");
    console.log();
    console.log("Available subcommands:");
    console.log("  list, ls                    List all profiles");
    console.log("  show, get <id>             Show profile details");
    console.log("  create, new <name> [env]   Create new profile");
    console.log("  delete, del, rm <id>       Delete profile");
    console.log();
    return;
  }

  switch (subcommand.toLowerCase()) {
    case "list":
    case "ls":
      const profileIds = await configService.listProfiles();

      if (profileIds.length === 0) {
        console.log(colors.yellow("No configuration profiles found."));
        console.log("Create one with: mcf config create <name> [environment]");
        return;
      }

      console.log(colors.blue("MCF Configuration Profiles:"));
      console.log();

      for (const profileId of profileIds) {
        console.log(`  ${colors.cyan(profileId)}`);
      }

      console.log();
      console.log(`Total: ${profileIds.length} profile(s)`);
      break;

    case "show":
    case "get":
      const profileId = args[0];
      if (!profileId) {
        console.log(colors.red("Profile ID is required"));
        console.log("Usage: mcf config show <profile-id>");
        return;
      }

      const profile = await configService.loadProfile(profileId);

      if (!profile) {
        console.log(colors.red(`Profile '${profileId}' not found`));
        return;
      }

      console.log(colors.blue(`Profile: ${profileId}`));
      console.log("─".repeat(50));
      console.log(`Name: ${colors.cyan(profile.name)}`);
      console.log(`Environment: ${colors.cyan(profile.environment || "undefined")}`);

      if (profile.description) {
        console.log(`Description: ${profile.description}`);
      }

      if (profile.config?.claude?.configDirectory) {
        console.log();
        console.log(colors.blue("Claude Configuration:"));
        console.log(`  Config Directory: ${colors.green(profile.config.claude.configDirectory)}`);
      }
      break;

    case "create":
    case "new":
      const [name, environment = "development"] = args;

      if (!name) {
        console.log(colors.red("Profile name is required"));
        console.log("Usage: mcf config create <name> [environment]");
        return;
      }

      const validEnvs = ["development", "production", "staging", "test"];
      if (!validEnvs.includes(environment)) {
        console.log(colors.red(`Invalid environment: ${environment}`));
        console.log(`Valid environments: ${validEnvs.join(", ")}`);
        return;
      }

      const newProfileId = configService.generateProfileId(name);
      if (await configService.profileExists(newProfileId)) {
        console.log(colors.red(`Profile '${newProfileId}' already exists`));
        return;
      }

      const newProfile = await configService.createProfile(name, environment);
      console.log(colors.green(`Profile '${newProfile.id}' created successfully`));
      break;

    case "delete":
    case "del":
    case "remove":
    case "rm":
      const deleteProfileId = args[0];
      if (!deleteProfileId) {
        console.log(colors.red("Profile ID is required"));
        return;
      }

      if (!(await configService.profileExists(deleteProfileId))) {
        console.log(colors.red(`Profile '${deleteProfileId}' not found`));
        return;
      }

      const success = await configService.deleteProfile(deleteProfileId);
      if (success) {
        console.log(colors.green(`Profile '${deleteProfileId}' deleted successfully`));
      }
      break;

    default:
      console.log(colors.red(`Unknown config subcommand: ${subcommand}`));
      break;
  }
}

/**
 * Run command implementation
 */
async function runCommand(subcommand, args, options) {
  const configService = new ConfigurationService();

  try {
    // Load profile configuration if specified
    let profileConfig = null;
    if (options.config) {
      profileConfig = await configService.loadProfile(options.config);
      if (!profileConfig) {
        console.log(colors.yellow(`⚠️  Profile '${options.config}' not found, using defaults`));
      } else {
        console.log(colors.blue(`🔧 Using profile: ${options.config}`));
      }
    }

    // Determine working directory
    let workingDirectory = options.workingDirectory || process.cwd();

    // Configure environment
    const env = { ...process.env };

    // Apply profile configuration for CLAUDE_CONFIG_DIR
    if (profileConfig?.config?.claude?.configDirectory) {
      env.CLAUDE_CONFIG_DIR = profileConfig.config.claude.configDirectory;
    }

    // Show execution details
    console.log(colors.blue(colors.bold("🚀 Starting Claude Code")));
    console.log();

    if (options.config) {
      console.log(colors.gray(`Profile: ${options.config}`));
    }

    if (options.debug) {
      console.log(colors.gray("Debug mode: enabled"));
    }

    if (options.project) {
      console.log(colors.gray(`Project: ${options.project}`));
    }

    if (options.passThroughArgs && options.passThroughArgs.length > 0) {
      console.log(colors.gray(`Arguments: ${options.passThroughArgs.join(" ")}`));
    }

    console.log(colors.gray(`Directory: ${workingDirectory}`));

    if (env.CLAUDE_CONFIG_DIR) {
      console.log(colors.gray(`🗂️  CLAUDE_CONFIG_DIR: ${env.CLAUDE_CONFIG_DIR}`));
    }

    console.log();

    // Build Claude arguments
    const claudeArgs = [];
    if (options.dangerousSkip) {
      claudeArgs.push("--dangerously-skip-permissions");
    }
    if (options.debug) {
      claudeArgs.push("--debug");
    }
    if (options.project) {
      claudeArgs.push("--project", options.project);
    }
    if (options.passThroughArgs) {
      claudeArgs.push(...options.passThroughArgs);
    }

    console.log(colors.blue("🚀 Launching Claude Code..."));
    console.log(colors.gray(`Command: claude ${claudeArgs.join(" ")}`));

    // Execute Claude directly
    const child = spawn("claude", claudeArgs, {
      stdio: "inherit",
      env,
      cwd: workingDirectory
    });

    return new Promise((resolve) => {
      child.on("exit", (code, signal) => {
        const result = {
          exitCode: code || 0,
          success: (code || 0) === 0,
          signal
        };

        if (result.success) {
          console.log();
          console.log(colors.green(`✅ Claude Code completed successfully`));
        } else if (signal) {
          console.log();
          console.log(colors.yellow(`⚠️  Claude Code terminated by signal: ${signal}`));
        } else {
          console.log();
          console.log(colors.red(`❌ Claude Code exited with code: ${code}`));
        }

        resolve(result);
      });

      child.on("error", (error) => {
        console.log();
        console.error(colors.red(`❌ Failed to start Claude: ${error.message}`));
        resolve({
          exitCode: 1,
          success: false,
          error: error.message
        });
      });
    });

  } catch (error) {
    console.error(colors.red(`❌ Failed to run Claude Code: ${error.message}`));
    throw error;
  }
}

/**
 * Install command implementation
 */
async function installCommand(subcommand, args, options) {
  console.log(colors.blue(colors.bold("🚀 MCF CLI Self-Installer")));
  console.log();
  
  try {
    const homeDir = os.homedir();
    const localBinDir = path.join(homeDir, ".local", "bin");
    const targetPath = path.join(localBinDir, "mcf");
    
    // Ensure ~/.local/bin exists
    console.log(colors.blue("📁 Ensuring ~/.local/bin directory exists..."));
    await fs.mkdir(localBinDir, { recursive: true });
    
    // Check if already installed
    if (!options.force) {
      try {
        await fs.access(targetPath);
        console.log(colors.yellow("⚠️  MCF CLI is already installed at ~/.local/bin/mcf"));
        console.log("Use --force to overwrite or run 'mcf --version' to check current version");
        return;
      } catch {
        // File doesn't exist, proceed with installation
      }
    }
    
    // Copy this file to ~/.local/bin/mcf
    console.log(colors.blue("📋 Installing MCF CLI to ~/.local/bin/mcf..."));
    const currentScript = __filename;
    
    // Copy the script directly (no modification needed since profiles are in $HOME/.mcf)
    await fs.copyFile(currentScript, targetPath);
    
    // Make executable
    await fs.chmod(targetPath, 0o755);
    
    console.log(colors.green("✅ MCF CLI installed successfully!"));
    console.log();
    console.log(colors.blue("Installation complete:"));
    console.log(`📁 Location: ${colors.cyan(targetPath)}`);
    console.log(`🔧 Version: ${colors.cyan(pkg.version)}`);
    console.log();
    console.log(colors.blue("Next steps:"));
    console.log("1. Add ~/.local/bin to your PATH if not already done:");
    console.log(colors.cyan('   export PATH="$HOME/.local/bin:$PATH"'));
    console.log("2. Reload your shell or run:");
    console.log(colors.cyan("   source ~/.zshrc"));
    console.log("3. Test the installation:");
    console.log(colors.cyan("   mcf --version"));
    console.log(colors.cyan("   mcf config list"));
    console.log(colors.cyan("   mcf run --config <profile-name>"));
    console.log();
    console.log(colors.blue("Profile management:"));
    console.log("• Create profiles with: mcf config create <name>");
    console.log("• Use profiles with: mcf run --config <profile-name>");
    console.log("• Profiles control CLAUDE_CONFIG_DIR for different Claude configurations");
    
  } catch (error) {
    console.log(colors.red("❌ Installation failed"));
    console.log(colors.red(`Error: ${error.message}`));
    console.log();
    console.log(colors.blue("Manual installation:"));
    console.log(`1. Copy ${__filename} to ~/.local/bin/mcf`);
    console.log("2. Make it executable: chmod +x ~/.local/bin/mcf");
    console.log("3. Add ~/.local/bin to your PATH");
    process.exit(1);
  }
}

/**
 * Status command implementation
 */
async function statusCommand(subcommand, args, options) {
  console.log(colors.blue(colors.bold("📊 MCF Status Check")));
  console.log();
  
  // Check Claude installation
  try {
    const child = spawn("claude", ["--version"], { stdio: "pipe" });
    await new Promise((resolve) => {
      child.on("exit", (code) => {
        if (code === 0) {
          console.log(colors.green("✅ Claude Code is installed"));
        } else {
          console.log(colors.red("❌ Claude Code not found"));
        }
        resolve();
      });
      child.on("error", () => {
        console.log(colors.red("❌ Claude Code not found"));
        resolve();
      });
    });
  } catch {
    console.log(colors.red("❌ Claude Code not found"));
  }

  // Check profiles
  const configService = new ConfigurationService();
  const profiles = await configService.listProfiles();
  console.log(`📁 Configuration profiles: ${profiles.length}`);
  
  if (profiles.length > 0) {
    console.log(colors.blue("Available profiles:"));
    profiles.forEach(profile => {
      console.log(`  • ${colors.cyan(profile)}`);
    });
  }
  
  console.log();
  console.log(colors.blue("MCF CLI is ready to use!"));
  console.log();
  console.log("Usage examples:");
  console.log(colors.cyan("  mcf config create myprofile"));
  console.log(colors.cyan("  mcf run --config myprofile"));
}

// ============================================================================
// MAIN CLI LOGIC
// ============================================================================

async function main() {
  const args = process.argv.slice(2);
  
  // Handle help and version directly (but not if they're after --)
  const hasDoubleDash = args.includes('--');
  const beforeDoubleDash = hasDoubleDash ? args.slice(0, args.indexOf('--')) : args;
  
  if (args.length === 0 || beforeDoubleDash.includes('--help') || beforeDoubleDash.includes('-h')) {
    showHelp();
    return;
  }

  if (beforeDoubleDash.includes('--version') || beforeDoubleDash.includes('-V')) {
    showVersion();
    return;
  }

  const command = args[0];

  try {
    switch (command) {
      case "config":
        const configSubcommand = args[1];
        const configArgs = args.slice(2);
        await configCommand(configSubcommand, configArgs, {});
        break;

      case "run":
        // Parse run command options
        const runOptions = {
          config: null,
          debug: false,
          project: null,
          workingDirectory: null,
          dangerousSkip: false,
          interactive: true,
          passThroughArgs: []
        };

        // Find -- separator
        const separatorIndex = args.indexOf("--");
        const runArgs = separatorIndex >= 0 ? args.slice(1, separatorIndex) : args.slice(1);
        const passThrough = separatorIndex >= 0 ? args.slice(separatorIndex + 1) : [];

        // Parse run options
        for (let i = 0; i < runArgs.length; i++) {
          const arg = runArgs[i];
          switch (arg) {
            case "--config":
            case "-c":
              runOptions.config = runArgs[++i];
              break;
            case "--debug":
            case "-d":
              runOptions.debug = true;
              break;
            case "--project":
            case "-p":
              runOptions.project = runArgs[++i];
              break;
            case "--working-directory":
            case "-w":
              runOptions.workingDirectory = runArgs[++i];
              break;
            case "--dangerous-skip":
              runOptions.dangerousSkip = true;
              break;
            case "--no-interactive":
              runOptions.interactive = false;
              break;
          }
        }

        runOptions.passThroughArgs = passThrough;
        await runCommand(null, [], runOptions);
        break;

      case "install":
        const installOptions = { force: args.includes("--force") || args.includes("-f") };
        await installCommand(null, [], installOptions);
        break;

      case "status":
        await statusCommand(null, [], {});
        break;

      default:
        console.log(colors.red(`Unknown command: ${command}`));
        console.log("Run 'mcf --help' for available commands");
        process.exit(1);
    }
  } catch (error) {
    console.error(colors.red(`❌ Command failed: ${error.message}`));
    process.exit(1);
  }
}

// Run the CLI
main();
