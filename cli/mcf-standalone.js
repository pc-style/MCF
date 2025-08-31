#!/usr/bin/env node

import { program } from "commander";
import chalk from "chalk";
import path from "path";
import { fileURLToPath } from "url";
import { spawn } from "child_process";
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
// CORE CLASSES AND INTERFACES
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
    console.log(chalk.yellow(`[WARN] ${message}`), data ? JSON.stringify(data) : "");
  }

  error(message, data) {
    console.log(chalk.red(`[ERROR] ${message}`), data ? JSON.stringify(data) : "");
  }

  debug(message, data) {
    if (process.env.DEBUG) {
      console.log(chalk.gray(`[DEBUG] ${message}`), data ? JSON.stringify(data) : "");
    }
  }
}

/**
 * Configuration service for profile management
 */
class ConfigurationService {
  constructor() {
    this.logger = new Logger("ConfigurationService");
    this.profilesDir = path.join(__dirname, ".mcf", "profiles");
    this.defaultProfilePath = path.join(__dirname, ".mcf", "default-profile.json");
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

/**
 * Project service for project management
 */
class ProjectService {
  constructor() {
    this.logger = new Logger("ProjectService");
    this.projectFileName = ".mcf-project.json";
  }

  generateProjectId(name) {
    return name.toLowerCase().replace(/[^a-z0-9]/g, "-").replace(/-+/g, "-").replace(/^-|-$/g, "");
  }

  async createProject(options, customPath) {
    const projectId = this.generateProjectId(options.name);
    const projectPath = customPath || path.join(process.cwd(), options.name);

    // Create project directory
    await fs.mkdir(projectPath, { recursive: true });

    const project = {
      id: projectId,
      name: options.name,
      description: options.description || `Project ${options.name}`,
      path: projectPath,
      environment: options.environment || "development",
      createdAt: new Date(),
      lastModified: new Date(),
      config: {
        profile: options.profile,
        workspace: options.workspace,
        settings: options.settings || {}
      },
      metadata: {
        version: "1.0.0",
        author: process.env.USER || "unknown",
        tags: [],
        custom: {}
      }
    };

    // Save project file
    const projectFilePath = path.join(projectPath, this.projectFileName);
    await fs.writeFile(projectFilePath, JSON.stringify(project, null, 2), "utf-8");

    this.logger.info(`Project '${projectId}' created successfully at ${projectPath}`);
    return project;
  }

  async listProjects() {
    // Simple implementation - just check current directory for projects
    try {
      const entries = await fs.readdir(process.cwd());
      const projects = [];

      for (const entry of entries) {
        const entryPath = path.join(process.cwd(), entry);
        const projectFilePath = path.join(entryPath, this.projectFileName);

        try {
          await fs.access(projectFilePath);
          const projectData = await fs.readFile(projectFilePath, "utf-8");
          const project = JSON.parse(projectData);
          project.path = entryPath;
          projects.push(project);
        } catch {
          // Not a project directory
        }
      }

      return projects;
    } catch {
      return [];
    }
  }

  async getCurrentProject() {
    try {
      const projectFilePath = path.join(process.cwd(), this.projectFileName);
      const projectData = await fs.readFile(projectFilePath, "utf-8");
      const project = JSON.parse(projectData);
      project.path = process.cwd();
      return project;
    } catch {
      return null;
    }
  }
}

// ============================================================================
// COMMAND IMPLEMENTATIONS
// ============================================================================

/**
 * Configuration command implementation
 */
class ConfigCommand {
  constructor() {
    this.configService = new ConfigurationService();
  }

  async execute(args = []) {
    if (args.length === 0) {
      return this.showHelp();
    }

    const [subcommand, ...subArgs] = args;

    switch (subcommand.toLowerCase()) {
      case "list":
      case "ls":
        return await this.listProfiles();
      case "show":
      case "get":
        return await this.showProfile(subArgs[0]);
      case "create":
      case "new":
        return await this.createProfile(subArgs);
      case "delete":
      case "del":
      case "remove":
      case "rm":
        return await this.deleteProfile(subArgs[0]);
      default:
        console.log(chalk.red(`Unknown config subcommand: ${subcommand}`));
        return this.showHelp();
    }
  }

  async listProfiles() {
    const profileIds = await this.configService.listProfiles();

    if (profileIds.length === 0) {
      console.log(chalk.yellow("No configuration profiles found."));
      console.log("Create one with: mcf config create <name> [environment]");
      return;
    }

    console.log(chalk.blue("MCF Configuration Profiles:"));
    console.log();

    for (const profileId of profileIds) {
      console.log(`  ${chalk.cyan(profileId)}`);
    }

    console.log();
    console.log(`Total: ${profileIds.length} profile(s)`);
  }

  async showProfile(profileId) {
    if (!profileId) {
      console.log(chalk.red("Profile ID is required"));
      console.log("Usage: mcf config show <profile-id>");
      return;
    }

    const profile = await this.configService.loadProfile(profileId);

    if (!profile) {
      console.log(chalk.red(`Profile '${profileId}' not found`));
      return;
    }

    console.log(chalk.blue(`Profile: ${profileId}`));
    console.log(chalk.gray("─".repeat(50)));
    console.log(`Name: ${chalk.cyan(profile.name)}`);
    console.log(`Environment: ${chalk.cyan(profile.environment || "undefined")}`);

    if (profile.description) {
      console.log(`Description: ${profile.description}`);
    }

    if (profile.config?.claude?.configDirectory) {
      console.log();
      console.log(chalk.blue("Claude Configuration:"));
      console.log(`  Config Directory: ${chalk.green(profile.config.claude.configDirectory)}`);
    }
  }

  async createProfile(args) {
    const [name, environment = "development"] = args;

    if (!name) {
      console.log(chalk.red("Profile name is required"));
      console.log("Usage: mcf config create <name> [environment]");
      return;
    }

    const validEnvs = ["development", "production", "staging", "test"];
    if (!validEnvs.includes(environment)) {
      console.log(chalk.red(`Invalid environment: ${environment}`));
      console.log(`Valid environments: ${validEnvs.join(", ")}`);
      return;
    }

    const profileId = this.configService.generateProfileId(name);
    if (await this.configService.profileExists(profileId)) {
      console.log(chalk.red(`Profile '${profileId}' already exists`));
      return;
    }

    const profile = await this.configService.createProfile(name, environment);
    console.log(chalk.green(`Profile '${profile.id}' created successfully`));
  }

  async deleteProfile(profileId) {
    if (!profileId) {
      console.log(chalk.red("Profile ID is required"));
      console.log("Usage: mcf config delete <profile-id>");
      return;
    }

    if (!(await this.configService.profileExists(profileId))) {
      console.log(chalk.red(`Profile '${profileId}' not found`));
      return;
    }

    const success = await this.configService.deleteProfile(profileId);
    if (success) {
      console.log(chalk.green(`Profile '${profileId}' deleted successfully`));
    }
  }

  showHelp() {
    console.log(chalk.blue("MCF Config Command Help"));
    console.log(chalk.gray("─".repeat(50)));
    console.log();
    console.log("Manage MCF configuration profiles");
    console.log();
    console.log(chalk.blue("Available subcommands:"));
    console.log();
    console.log("  list, ls                    List all profiles");
    console.log("  show, get <id>             Show profile details");
    console.log("  create, new <name> [env]   Create new profile");
    console.log("  delete, del, rm <id>       Delete profile");
    console.log();
    console.log(chalk.blue("Examples:"));
    console.log();
    console.log("  mcf config list");
    console.log("  mcf config create myapp production");
    console.log("  mcf config show myapp");
    console.log();
  }
}

/**
 * Project command implementation
 */
class ProjectCommand {
  constructor() {
    this.projectService = new ProjectService();
  }

  async execute(args = []) {
    if (args.length === 0) {
      return this.showHelp();
    }

    const [subcommand, ...subArgs] = args;

    switch (subcommand.toLowerCase()) {
      case "list":
      case "ls":
        return await this.listProjects();
      case "show":
      case "info":
        return await this.showProject(subArgs[0]);
      case "create":
      case "new":
        return await this.createProject(subArgs);
      case "current":
        return await this.showCurrentProject();
      default:
        console.log(chalk.red(`Unknown project subcommand: ${subcommand}`));
        return this.showHelp();
    }
  }

  async listProjects() {
    const projects = await this.projectService.listProjects();

    if (projects.length === 0) {
      console.log(chalk.yellow("No projects found."));
      console.log("Create one with: mcf project create <name> [description]");
      return;
    }

    console.log(chalk.blue("MCF Projects:"));
    console.log();

    for (const project of projects) {
      console.log(`  ${chalk.cyan(project.id)}`);
      console.log(`    Name: ${project.name}`);
      console.log(`    Environment: ${project.environment}`);
      console.log(`    Path: ${chalk.gray(project.path)}`);
      console.log();
    }

    console.log(`Total: ${projects.length} project(s)`);
  }

  async showProject(projectId) {
    if (!projectId) {
      console.log(chalk.red("Project ID is required"));
      return;
    }

    const projects = await this.projectService.listProjects();
    const project = projects.find(p => p.id === projectId);

    if (!project) {
      console.log(chalk.red(`Project '${projectId}' not found`));
      return;
    }

    console.log(chalk.blue(`Project: ${project.id}`));
    console.log(`Name: ${project.name}`);
    console.log(`Path: ${project.path}`);
    console.log(`Environment: ${project.environment}`);
  }

  async createProject(args) {
    const [name, ...descriptionParts] = args;
    const description = descriptionParts.join(" ");

    if (!name) {
      console.log(chalk.red("Project name is required"));
      return;
    }

    const options = {
      name,
      description: description || undefined,
      environment: "development"
    };

    const project = await this.projectService.createProject(options);
    console.log(chalk.green(`Project '${project.id}' created successfully`));
    console.log(`Path: ${project.path}`);
  }

  async showCurrentProject() {
    const currentProject = await this.projectService.getCurrentProject();

    if (!currentProject) {
      console.log(chalk.yellow("No current project"));
      return;
    }

    console.log(chalk.blue("Current Project:"));
    console.log(`ID: ${chalk.cyan(currentProject.id)}`);
    console.log(`Name: ${currentProject.name}`);
    console.log(`Path: ${chalk.gray(currentProject.path)}`);
  }

  showHelp() {
    console.log(chalk.blue("MCF Project Command Help"));
    console.log(chalk.gray("─".repeat(50)));
    console.log();
    console.log("  list, ls                    List all projects");
    console.log("  show, info <id>             Show project details");
    console.log("  create, new <name> [desc]   Create new project");
    console.log("  current                     Show current project");
    console.log();
  }
}

/**
 * Run command implementation
 */
class RunCommand {
  constructor() {
    this.configService = new ConfigurationService();
    this.projectService = new ProjectService();
    this.logger = new Logger("RunCommand");
  }

  async execute(args = []) {
    const parsedArgs = this.parseArguments(args);

    try {
      // Load profile configuration if specified
      let profileConfig = null;
      if (parsedArgs.profile) {
        profileConfig = await this.configService.loadProfile(parsedArgs.profile);
        if (!profileConfig) {
          console.log(chalk.yellow(`⚠️  Profile '${parsedArgs.profile}' not found, using defaults`));
        } else {
          console.log(chalk.blue(`🔧 Using profile: ${parsedArgs.profile}`));
        }
      }

      // Determine working directory
      let workingDirectory = parsedArgs.workingDirectory || process.cwd();

      // Check if we're in a project directory
      const currentProject = await this.projectService.getCurrentProject();
      if (currentProject && !parsedArgs.workingDirectory) {
        workingDirectory = currentProject.path;
        console.log(chalk.blue(`📁 Using project directory: ${workingDirectory}`));
      }

      // Configure environment
      const env = { ...process.env };

      // Apply profile configuration for CLAUDE_CONFIG_DIR
      if (profileConfig?.config?.claude?.configDirectory) {
        env.CLAUDE_CONFIG_DIR = profileConfig.config.claude.configDirectory;
      }

      // Show execution details
      console.log(chalk.blue.bold("🚀 Starting Claude Code"));
      console.log();

      if (parsedArgs.profile) {
        console.log(chalk.gray(`Profile: ${parsedArgs.profile}`));
      }

      if (parsedArgs.debug) {
        console.log(chalk.gray("Debug mode: enabled"));
      }

      if (parsedArgs.projectName) {
        console.log(chalk.gray(`Project: ${parsedArgs.projectName}`));
      }

      if (parsedArgs.passThroughArgs && parsedArgs.passThroughArgs.length > 0) {
        console.log(chalk.gray(`Arguments: ${parsedArgs.passThroughArgs.join(" ")}`));
      }

      console.log(chalk.gray(`Directory: ${workingDirectory}`));

      if (env.CLAUDE_CONFIG_DIR) {
        console.log(chalk.gray(`🗂️  CLAUDE_CONFIG_DIR: ${env.CLAUDE_CONFIG_DIR}`));
      }

      console.log();

      // Build Claude arguments
      const args = [];
      if (parsedArgs.dangerousSkip) {
        args.push("--dangerously-skip-permissions");
      }
      if (parsedArgs.debug) {
        args.push("--debug");
      }
      if (parsedArgs.projectName) {
        args.push("--project", parsedArgs.projectName);
      }
      if (parsedArgs.passThroughArgs) {
        args.push(...parsedArgs.passThroughArgs);
      }

      console.log(chalk.blue("🚀 Launching Claude Code..."));

      // Execute Claude directly
      const child = spawn("claude", args, {
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
            console.log(chalk.green(`✅ Claude Code completed successfully`));
          } else if (signal) {
            console.log();
            console.log(chalk.yellow(`⚠️  Claude Code terminated by signal: ${signal}`));
          } else {
            console.log();
            console.log(chalk.red(`❌ Claude Code exited with code: ${code}`));
          }

          resolve(result);
        });

        child.on("error", (error) => {
          console.log();
          console.error(chalk.red(`❌ Failed to start Claude: ${error.message}`));
          resolve({
            exitCode: 1,
            success: false,
            error: error.message
          });
        });
      });

    } catch (error) {
      console.error(chalk.red(`❌ Failed to run Claude Code: ${error.message}`));
      throw error;
    }
  }

  parseArguments(args) {
    const result = {
      debug: false,
      profile: null,
      projectName: null,
      workingDirectory: null,
      dangerousSkip: false,
      passThroughArgs: []
    };

    // Find separator between MCF args and pass-through args
    const separatorIndex = args.indexOf("--");
    const mcfArgs = separatorIndex >= 0 ? args.slice(0, separatorIndex) : args;
    const passThrough = separatorIndex >= 0 ? args.slice(separatorIndex + 1) : [];

    // Parse MCF-specific flags
    for (let i = 0; i < mcfArgs.length; i++) {
      const arg = mcfArgs[i];

      switch (arg) {
        case "-d":
        case "--debug":
          result.debug = true;
          break;
        case "-c":
        case "--config":
          if (i + 1 < mcfArgs.length) {
            result.profile = mcfArgs[++i];
          }
          break;
        case "-p":
        case "--project":
          if (i + 1 < mcfArgs.length) {
            result.projectName = mcfArgs[++i];
          }
          break;
        case "-w":
        case "--working-directory":
          if (i + 1 < mcfArgs.length) {
            result.workingDirectory = mcfArgs[++i];
          }
          break;
        case "--dangerous-skip":
          result.dangerousSkip = true;
          break;
      }
    }

    result.passThroughArgs = passThrough;
    return result;
  }
}

// ============================================================================
// MAIN CLI SETUP
// ============================================================================

// Initialize services
const configService = new ConfigurationService();
const projectService = new ProjectService();

// Configure program
program
  .name("mcf")
  .description(pkg.description)
  .version(pkg.version);

// Config command
program
  .command("config")
  .description("Manage MCF configuration profiles")
  .argument("[subcommand]", "Config subcommand")
  .argument("[args...]", "Subcommand arguments")
  .action(async (subcommand, args, options) => {
    try {
      const configCommand = new ConfigCommand();
      const allArgs = subcommand ? [subcommand, ...args] : [];
      await configCommand.execute(allArgs);
    } catch (error) {
      console.error(chalk.red(`Config command failed: ${error.message}`));
      process.exit(1);
    }
  });

// Project command
program
  .command("project")
  .description("Manage MCF projects and workspaces")
  .argument("[subcommand]", "Project subcommand")
  .argument("[args...]", "Subcommand arguments")
  .action(async (subcommand, args, options) => {
    try {
      const projectCommand = new ProjectCommand();
      const allArgs = subcommand ? [subcommand, ...args] : [];
      await projectCommand.execute(allArgs);
    } catch (error) {
      console.error(chalk.red(`Project command failed: ${error.message}`));
      process.exit(1);
    }
  });

// Run command
program
  .command("run")
  .description("Execute Claude Code with configuration and flags")
  .option("-d, --debug", "Enable debug mode")
  .option("-c, --config <profile>", "Use specific configuration profile")
  .option("-p, --project <name>", "Set Claude project name")
  .option("-w, --working-directory <path>", "Set working directory")
  .option("--dangerous-skip", "Skip permission checks (dangerous)")
  .allowUnknownOption()
  .action(async (options, command) => {
    try {
      const runCommand = new RunCommand();
      
      // Get all arguments after "run"
      const runArgsIndex = process.argv.indexOf("run") + 1;
      const allArgs = process.argv.slice(runArgsIndex);

      await runCommand.execute(allArgs);
    } catch (error) {
      console.error(chalk.red(`Run command failed: ${error.message}`));
      process.exit(1);
    }
  });

// Install command (self-installing)
program
  .command("install")
  .description("Install MCF CLI to ~/.local/bin/mcf")
  .option("-f, --force", "Force overwrite existing installation")
  .action(async (options) => {
    console.log(chalk.blue.bold("🚀 MCF CLI Self-Installer"));
    console.log();
    
    try {
      const homeDir = os.homedir();
      const localBinDir = path.join(homeDir, ".local", "bin");
      const targetPath = path.join(localBinDir, "mcf");
      
      // Ensure ~/.local/bin exists
      console.log(chalk.blue("📁 Ensuring ~/.local/bin directory exists..."));
      await fs.mkdir(localBinDir, { recursive: true });
      
      // Check if already installed
      if (!options.force) {
        try {
          await fs.access(targetPath);
          console.log(chalk.yellow("⚠️  MCF CLI is already installed at ~/.local/bin/mcf"));
          console.log("Use --force to overwrite or run 'mcf --version' to check current version");
          return;
        } catch {
          // File doesn't exist, proceed with installation
        }
      }
      
      // Copy this file to ~/.local/bin/mcf
      console.log(chalk.blue("📋 Installing MCF CLI to ~/.local/bin/mcf..."));
      const currentScript = __filename;
      await fs.copyFile(currentScript, targetPath);
      
      // Make executable
      await fs.chmod(targetPath, 0o755);
      
      // Verify installation
      try {
        const child = spawn(targetPath, ["--version"], { stdio: "pipe" });
        await new Promise((resolve, reject) => {
          let output = "";
          child.stdout.on("data", (data) => {
            output += data.toString();
          });
          
          child.on("exit", (code) => {
            if (code === 0 && output.includes("1.0.0")) {
              resolve();
            } else {
              reject(new Error("Installation verification failed"));
            }
          });
          
          child.on("error", reject);
        });
        
        console.log(chalk.green("✅ MCF CLI installed successfully!"));
        console.log();
        console.log(chalk.blue("Installation complete:"));
        console.log(`📁 Location: ${chalk.cyan(targetPath)}`);
        console.log(`🔧 Version: ${chalk.cyan("1.0.0")}`);
        console.log();
        console.log(chalk.blue("Next steps:"));
        console.log("1. Add ~/.local/bin to your PATH if not already done:");
        console.log(chalk.cyan('   export PATH="$HOME/.local/bin:$PATH"'));
        console.log("2. Reload your shell or run:");
        console.log(chalk.cyan("   source ~/.zshrc"));
        console.log("3. Test the installation:");
        console.log(chalk.cyan("   mcf --version"));
        console.log(chalk.cyan("   mcf config list"));
        console.log(chalk.cyan("   mcf run --config mcf"));
        console.log();
        console.log(chalk.blue("Profile management:"));
        console.log("• Your existing profiles are preserved in cli/.mcf/profiles/");
        console.log("• Create new profiles with: mcf config create <name>");
        console.log("• Use profiles with: mcf run --config <profile-name>");
        
      } catch (error) {
        console.log(chalk.red("❌ Installation verification failed"));
        console.log(chalk.red(`Error: ${error.message}`));
        console.log();
        console.log(chalk.yellow("The file was copied but may not be working correctly."));
        console.log(`Please check: ${targetPath}`);
      }
      
    } catch (error) {
      console.log(chalk.red("❌ Installation failed"));
      console.log(chalk.red(`Error: ${error.message}`));
      console.log();
      console.log(chalk.blue("Manual installation:"));
      console.log(`1. Copy ${__filename} to ~/.local/bin/mcf`);
      console.log("2. Make it executable: chmod +x ~/.local/bin/mcf");
      console.log("3. Add ~/.local/bin to your PATH");
      process.exit(1);
    }
  });

// Status command (simplified)
program
  .command("status")
  .description("Check MCF installation status")
  .action(async (options) => {
    console.log(chalk.blue.bold("📊 MCF Status Check"));
    console.log();
    
    // Check Claude installation
    try {
      const child = spawn("claude", ["--version"], { stdio: "pipe" });
      await new Promise((resolve) => {
        child.on("exit", (code) => {
          if (code === 0) {
            console.log(chalk.green("✅ Claude Code is installed"));
          } else {
            console.log(chalk.red("❌ Claude Code not found"));
          }
          resolve();
        });
        child.on("error", () => {
          console.log(chalk.red("❌ Claude Code not found"));
          resolve();
        });
      });
    } catch {
      console.log(chalk.red("❌ Claude Code not found"));
    }

    // Check profiles
    const profiles = await configService.listProfiles();
    console.log(`📁 Configuration profiles: ${profiles.length}`);
    
    // Check projects
    const projects = await projectService.listProjects();
    console.log(`📂 Projects found: ${projects.length}`);
    
    console.log();
    console.log(chalk.blue("MCF CLI is ready to use!"));
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
