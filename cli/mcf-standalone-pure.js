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
  description:
    "MCF (My Claude Flow) CLI - Installation, configuration and setup tool",
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
  bold: (text) => `\x1b[1m${text}\x1b[0m`,
};

// Simple argument parser (replacing commander)
function parseArgs() {
  const args = process.argv.slice(2);
  const result = {
    command: null,
    subcommand: null,
    options: {},
    args: [],
  };

  if (args.length === 0) {
    return result;
  }

  // Handle help and version
  if (args.includes("--help") || args.includes("-h")) {
    result.showHelp = true;
    return result;
  }

  if (args.includes("--version") || args.includes("-V")) {
    result.showVersion = true;
    return result;
  }

  // Parse command and options
  result.command = args[0];

  // Simple option parsing
  for (let i = 1; i < args.length; i++) {
    const arg = args[i];
    if (arg.startsWith("--")) {
      const key = arg.substring(2);
      if (i + 1 < args.length && !args[i + 1].startsWith("--")) {
        result.options[key] = args[i + 1];
        i++; // skip next arg
      } else {
        result.options[key] = true;
      }
    } else if (arg.startsWith("-")) {
      result.options[arg.substring(1)] = true;
    } else {
      result.args.push(arg);
    }
  }

  if (result.args.length > 0) {
    result.subcommand = result.args[0];
    result.args = result.args.slice(1);
  }

  return result;
}

// ============================================================================
// CONFIGURATION SERVICE (embedded to avoid external dependencies)
// ============================================================================

class ConfigurationService {
  constructor() {
    this.configDir = path.join(os.homedir(), ".mcf");
    this.profilesDir = path.join(this.configDir, "profiles");
  }

  async ensureConfigDir() {
    await fs.mkdir(this.profilesDir, { recursive: true });
  }

  async listProfiles() {
    try {
      await this.ensureConfigDir();
      const files = await fs.readdir(this.profilesDir);
      return files
        .filter((file) => file.endsWith(".json"))
        .map((file) => path.basename(file, ".json"));
    } catch (error) {
      return [];
    }
  }

  async getProfile(profileName) {
    const profilePath = path.join(this.profilesDir, `${profileName}.json`);
    try {
      const content = await fs.readFile(profilePath, "utf-8");
      return JSON.parse(content);
    } catch (error) {
      return null;
    }
  }

  async setProfile(profileName, profileData) {
    await this.ensureConfigDir();
    const profilePath = path.join(this.profilesDir, `${profileName}.json`);
    await fs.writeFile(profilePath, JSON.stringify(profileData, null, 2));
  }

  async deleteProfile(profileName) {
    const profilePath = path.join(this.profilesDir, `${profileName}.json`);
    try {
      await fs.unlink(profilePath);
      return true;
    } catch (error) {
      return false;
    }
  }

  createDefaultProfile(name, environment = "development") {
    return {
      id: this.generateProfileId(name),
      name: name,
      environment: environment,
      version: "1.0.0",
      created: new Date().toISOString(),
      lastUpdated: new Date().toISOString(),
      config: {
        timeout: 60000,
        logLevel: "info",
        maxRetries: 3,
        claude: {
          configDirectory: path.join(
            this.configDir,
            "claude-configs",
            this.generateProfileId(name),
          ),
        },
      },
      permissions: {
        allowedServices: ["claude", "filesystem"],
        blockedServices: [],
      },
    };
  }

  generateProfileId(name) {
    return name
      .toLowerCase()
      .replace(/[^a-z0-9]/g, "-")
      .replace(/-+/g, "-")
      .replace(/^-|-$/g, "");
  }
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

/**
 * Simple spinner implementation
 */
function createSpinner(text) {
  const frames = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];
  let frameIndex = 0;
  let intervalId = null;
  let currentText = text;

  return {
    start() {
      process.stdout.write("\x1B[?25l"); // Hide cursor
      intervalId = setInterval(() => {
        process.stdout.write(`\r${frames[frameIndex]} ${currentText}`);
        frameIndex = (frameIndex + 1) % frames.length;
      }, 100);
    },

    succeed(text) {
      if (intervalId) {
        clearInterval(intervalId);
        intervalId = null;
      }
      process.stdout.write(`\r✅ ${text || currentText}\n`);
      process.stdout.write("\x1B[?25h"); // Show cursor
    },

    fail(text) {
      if (intervalId) {
        clearInterval(intervalId);
        intervalId = null;
      }
      process.stdout.write(`\r❌ ${text || currentText}\n`);
      process.stdout.write("\x1B[?25h"); // Show cursor
    },

    text(newText) {
      currentText = newText;
    },
  };
}

/**
 * Check if git is available
 */
async function checkGitAvailability() {
  return new Promise((resolve, reject) => {
    const git = spawn("git", ["--version"], { stdio: "pipe" });

    git.on("close", (code) => {
      if (code === 0) {
        resolve();
      } else {
        reject(
          new Error(
            "Git is not available. Please install git to use the init command.",
          ),
        );
      }
    });

    git.on("error", (error) => {
      reject(
        new Error(
          `Git is not available: ${error.message}. Please install git to use the init command.`,
        ),
      );
    });
  });
}

/**
 * Check if target directory exists and handle accordingly
 */
async function checkTargetDirectory(targetDirectory, force) {
  const absolutePath = path.resolve(process.cwd(), targetDirectory);

  try {
    const stats = await fs.stat(absolutePath);

    if (stats.isDirectory()) {
      // Check if directory is empty
      const files = await fs.readdir(absolutePath);
      if (files.length > 0) {
        if (!force) {
          throw new Error(
            `Directory '${targetDirectory}' already exists and is not empty. Use --force to overwrite.`,
          );
        } else {
          console.log(
            colors.yellow(
              `⚠️  Directory '${targetDirectory}' exists. Overwriting due to --force flag.`,
            ),
          );
        }
      }
    } else {
      throw new Error(`'${targetDirectory}' exists but is not a directory.`);
    }
  } catch (error) {
    if (error.code === "ENOENT") {
      // Directory doesn't exist, which is fine
    } else {
      throw error;
    }
  }
}

/**
 * Clone the MCF repository
 */
async function cloneRepository(repositoryUrl, targetDirectory, options) {
  const spinner = createSpinner("Cloning MCF repository...");
  spinner.start();

  return new Promise((resolve, reject) => {
    // Build git clone arguments
    const gitArgs = ["clone"];

    if (options.shallow !== false) {
      gitArgs.push("--depth", "1");
    }

    if (options.branch) {
      gitArgs.push("--branch", options.branch);
    }

    gitArgs.push(repositoryUrl, targetDirectory);

    // Execute git clone
    const git = spawn("git", gitArgs, {
      stdio: ["inherit", "pipe", "pipe"],
      cwd: process.cwd(),
    });

    let stdout = "";
    let stderr = "";

    git.stdout.on("data", (data) => {
      stdout += data.toString();
      // Update spinner with progress if available
      const output = data.toString().trim();
      if (output) {
        spinner.text(`Cloning MCF repository... ${output.split("\n").pop()}`);
      }
    });

    git.stderr.on("data", (data) => {
      stderr += data.toString();
      // Git often sends progress to stderr
      const output = data.toString().trim();
      if (output && !output.toLowerCase().includes("error")) {
        spinner.text(`Cloning MCF repository... ${output.split("\n").pop()}`);
      }
    });

    git.on("close", (code) => {
      if (code === 0) {
        spinner.succeed("MCF repository cloned successfully!");
        resolve();
      } else {
        spinner.fail("Failed to clone MCF repository");
        const errorMessage = stderr || stdout || "Unknown git error";
        reject(new Error(`Git clone failed: ${errorMessage}`));
      }
    });

    git.on("error", (error) => {
      spinner.fail("Failed to clone MCF repository");
      reject(new Error(`Git clone failed: ${error.message}`));
    });
  });
}

// ============================================================================
// HELP AND VERSION
// ============================================================================

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
  console.log(
    "  init [directory] [options]      Initialize a new MCF project by cloning the repository",
  );
  console.log(
    "  config [subcommand] [args...]   Manage MCF configuration profiles",
  );
  console.log(
    "  run [options]                   Execute Claude Code with configuration and flags",
  );
  console.log(
    "  install                         Install MCF CLI to ~/.local/bin/mcf",
  );
  console.log(
    "  status                          Check MCF installation status",
  );
  console.log("  help [command]                  display help for command");
  console.log();
  console.log("Examples:");
  console.log("  mcf init                        # Clone into 'MCF' directory");
  console.log(
    "  mcf init my-project             # Clone into 'my-project' directory",
  );
  console.log("  mcf init --branch develop       # Clone 'develop' branch");
  console.log("  mcf config list");
  console.log("  mcf config create myprofile");
  console.log("  mcf run --config myprofile");
  console.log("  mcf install");
}

/**
 * Show version information
 */
function showVersion() {
  console.log(`${pkg.name} ${pkg.version}`);
}

// ============================================================================
// COMMAND IMPLEMENTATIONS
// ============================================================================

/**
 * Init command implementation
 */
async function initCommand(subcommand, args, options) {
  console.log(colors.blue(colors.bold("🚀 MCF Project Initialization")));
  console.log();

  const targetDirectory = args.length > 0 ? args[0] : "MCF";
  const repositoryUrl = "https://github.com/pc-style/MCF.git";

  try {
    // Check if git is available
    await checkGitAvailability();

    // Check if target directory already exists
    await checkTargetDirectory(targetDirectory, options.force);

    // Clone the repository
    await cloneRepository(repositoryUrl, targetDirectory, {
      branch: options.branch || options.b,
      shallow: !options["no-shallow"],
      force: options.force || options.f,
    });

    // Display success message and next steps
    console.log();
    console.log(
      colors.green(colors.bold("✅ MCF project initialized successfully!")),
    );
    console.log();
    console.log(colors.blue("Next steps:"));
    console.log(`  1. ${colors.yellow(`cd ${targetDirectory}`)}`);
    console.log(`  2. ${colors.yellow("mcf install")} - Install MCF framework`);
    console.log(
      `  3. ${colors.yellow("mcf config create <profile>")} - Create a configuration profile`,
    );
    console.log(
      `  4. ${colors.yellow("mcf run --config <profile>")} - Start using MCF with Claude Code`,
    );
    console.log();
    console.log(colors.blue("Documentation:"));
    console.log(
      `  • README: ${colors.cyan(path.join(targetDirectory, "README.md"))}`,
    );
    console.log(
      `  • Documentation: ${colors.cyan(path.join(targetDirectory, "docs/"))}`,
    );
    console.log();
    console.log(colors.gray("Happy coding with MCF! 🚀"));
  } catch (error) {
    console.error(
      colors.red(`❌ Failed to initialize MCF project: ${error.message}`),
    );
    throw error;
  }
}

/**
 * Config command implementation
 */
async function configCommand(subcommand, args, options) {
  const configService = new ConfigurationService();

  if (!subcommand) {
    console.log(colors.blue(colors.bold("MCF Configuration Management")));
    console.log();
    console.log("Available subcommands:");
    console.log("  list, ls           List all profiles");
    console.log("  create <name>      Create new profile");
    console.log("  show <name>        Show profile details");
    console.log("  delete <name>      Delete profile");
    console.log();
    console.log("Usage: mcf config <subcommand> [args...]");
    return;
  }

  switch (subcommand) {
    case "list":
    case "ls":
      const profiles = await configService.listProfiles();
      if (profiles.length === 0) {
        console.log(colors.yellow("No configuration profiles found."));
        console.log("Create one with: mcf config create <name>");
      } else {
        console.log(colors.blue("MCF Configuration Profiles:"));
        console.log();
        profiles.forEach((profile) => {
          console.log(`  ${colors.cyan(profile)}`);
        });
        console.log();
        console.log(`Total: ${profiles.length} profile(s)`);
      }
      break;

    case "create":
      if (args.length === 0) {
        console.log(colors.red("Profile name is required"));
        console.log("Usage: mcf config create <name> [environment]");
        return;
      }

      const profileName = args[0];
      const environment = args[1] || "development";

      // Check if profile already exists
      const existingProfile = await configService.getProfile(profileName);
      if (existingProfile) {
        console.log(colors.red(`Profile '${profileName}' already exists`));
        return;
      }

      const newProfile = configService.createDefaultProfile(
        profileName,
        environment,
      );
      await configService.setProfile(profileName, newProfile);

      console.log(
        colors.green(`Profile '${profileName}' created successfully`),
      );
      console.log(`Environment: ${environment}`);
      console.log(
        `Claude config directory: ${newProfile.config.claude.configDirectory}`,
      );
      console.log();
      console.log("Use with: mcf run --config " + profileName);
      break;

    case "show":
      if (args.length === 0) {
        console.log(colors.red("Profile name is required"));
        console.log("Usage: mcf config show <name>");
        return;
      }

      const profile = await configService.getProfile(args[0]);
      if (!profile) {
        console.log(colors.red(`Profile '${args[0]}' not found`));
        return;
      }

      console.log(colors.blue(`Profile: ${args[0]}`));
      console.log(colors.gray("─".repeat(50)));
      console.log(`Name: ${colors.cyan(profile.name)}`);
      console.log(`Environment: ${colors.cyan(profile.environment)}`);
      console.log(`Version: ${profile.version}`);
      console.log(
        `Claude Config Dir: ${colors.green(profile.config.claude.configDirectory)}`,
      );

      if (profile.created) {
        console.log(`Created: ${new Date(profile.created).toLocaleString()}`);
      }
      break;

    case "delete":
      if (args.length === 0) {
        console.log(colors.red("Profile name is required"));
        console.log("Usage: mcf config delete <name>");
        return;
      }

      const deleteResult = await configService.deleteProfile(args[0]);
      if (deleteResult) {
        console.log(colors.green(`Profile '${args[0]}' deleted successfully`));
      } else {
        console.log(colors.red(`Profile '${args[0]}' not found`));
      }
      break;

    default:
      console.log(colors.red(`Unknown config subcommand: ${subcommand}`));
      console.log("Available: list, create, show, delete");
  }
}

/**
 * Run command implementation
 */
async function runCommand(subcommand, args, options) {
  console.log(colors.blue(colors.bold("🚀 Starting Claude Code")));
  console.log();

  // Load configuration if profile specified
  let profile = null;
  let claudeConfigDir = null;

  if (options.config) {
    const configService = new ConfigurationService();
    profile = await configService.getProfile(options.config);

    if (!profile) {
      console.log(
        colors.yellow(
          `⚠️  Profile '${options.config}' not found, using defaults`,
        ),
      );
    } else {
      console.log(colors.blue(`🔧 Using profile: ${options.config}`));
      claudeConfigDir = profile.config.claude.configDirectory;
    }
  }

  // Show execution details
  if (options.debug) {
    console.log(colors.gray("Debug mode: enabled"));
  }

  if (options.project) {
    console.log(colors.gray(`Project: ${options.project}`));
  }

  if (options.passThroughArgs && options.passThroughArgs.length > 0) {
    console.log(colors.gray(`Arguments: ${options.passThroughArgs.join(" ")}`));
  }

  const workingDirectory = options.workingDirectory || process.cwd();
  console.log(colors.gray(`Directory: ${workingDirectory}`));

  if (claudeConfigDir) {
    console.log(colors.gray(`CLAUDE_CONFIG_DIR: ${claudeConfigDir}`));
  }
  console.log();

  try {
    // Configure environment
    const env = { ...process.env };
    if (claudeConfigDir) {
      env.CLAUDE_CONFIG_DIR = claudeConfigDir;

      // Ensure Claude config directory exists
      await fs.mkdir(claudeConfigDir, { recursive: true });
    }

    console.log(colors.blue("🚀 Launching Claude Code..."));

    // Build arguments
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

    // Execute Claude directly
    const child = spawn("claude", claudeArgs, {
      stdio: "inherit",
      env,
      cwd: workingDirectory,
    });

    return new Promise((resolve) => {
      child.on("exit", (code, signal) => {
        if (code === 0) {
          console.log();
          console.log(colors.green("✅ Claude Code completed successfully"));
        } else if (signal) {
          console.log();
          console.log(
            colors.yellow(`⚠️  Claude Code terminated by signal: ${signal}`),
          );
        } else {
          console.log();
          console.log(colors.red(`❌ Claude Code exited with code: ${code}`));
        }
        resolve();
      });

      child.on("error", (error) => {
        console.log();
        console.error(
          colors.red(`❌ Failed to start Claude: ${error.message}`),
        );
        resolve();
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
        console.log(
          colors.yellow("⚠️  MCF CLI is already installed at ~/.local/bin/mcf"),
        );
        console.log(
          "Use --force to overwrite or run 'mcf --version' to check current version",
        );
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
    console.log(
      "• Profiles control CLAUDE_CONFIG_DIR for different Claude configurations",
    );
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
    profiles.forEach((profile) => {
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
  const hasDoubleDash = args.includes("--");
  const beforeDoubleDash = hasDoubleDash
    ? args.slice(0, args.indexOf("--"))
    : args;

  if (
    args.length === 0 ||
    beforeDoubleDash.includes("--help") ||
    beforeDoubleDash.includes("-h")
  ) {
    showHelp();
    return;
  }

  if (
    beforeDoubleDash.includes("--version") ||
    beforeDoubleDash.includes("-V")
  ) {
    showVersion();
    return;
  }

  const command = args[0];

  try {
    switch (command) {
      case "init":
        // Parse init command options
        const initOptions = {
          force: args.includes("--force") || args.includes("-f"),
          branch: null,
          "no-shallow": args.includes("--no-shallow"),
        };

        // Find branch option
        const branchIndex = args.findIndex(
          (arg) => arg === "--branch" || arg === "-b",
        );
        if (branchIndex >= 0 && branchIndex + 1 < args.length) {
          initOptions.branch = args[branchIndex + 1];
        }

        // Get directory argument (first non-option argument)
        const initArgs = args
          .slice(1)
          .filter((arg) => !arg.startsWith("-") && arg !== initOptions.branch);

        await initCommand(null, initArgs, initOptions);
        break;

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
          passThroughArgs: [],
        };

        // Find -- separator
        const separatorIndex = args.indexOf("--");
        const runArgs =
          separatorIndex >= 0 ? args.slice(1, separatorIndex) : args.slice(1);
        const passThrough =
          separatorIndex >= 0 ? args.slice(separatorIndex + 1) : [];

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
        const installOptions = {
          force: args.includes("--force") || args.includes("-f"),
        };
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
