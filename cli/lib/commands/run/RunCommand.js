import { BaseCommand } from "../../core/interfaces/BaseCommand.js";
import { ServiceRegistry } from "../../core/registry/ServiceRegistry.js";
import { ILogger } from "../../core/logging/ILogger.js";
import { LoggerFactory } from "../../core/logging/LoggerFactory.js";
import { CLIError } from "../../types/CLITypes.js";
import chalk from "chalk";

/**
 * RunCommand - MCF CLI Claude Code Runner
 * Provides enhanced Claude Code execution with flags and pass-through arguments
 */
export class RunCommand extends BaseCommand {
  constructor(serviceRegistry) {
    super();
    this.serviceRegistry = serviceRegistry;
    this.logger = LoggerFactory.getLogger("RunCommand");
    this.claudeService = null;
    this.configService = null;
    this.projectService = null;
  }

  static get metadata() {
    return {
      name: "RunCommand",
      description: "Execute Claude Code with configuration and flags",
      category: "run",
      version: "1.0.0",
      dependencies: {
        services: ["IClaudeService", "IConfigurationService", "IProjectService"],
        commands: [],
        external: []
      }
    };
  }

  async initialize() {
    try {
      this.claudeService = this.serviceRegistry.getService("IClaudeService");
      this.configService = this.serviceRegistry.getService("IConfigurationService");
      this.projectService = this.serviceRegistry.getService("IProjectService");
      this.logger.debug("RunCommand initialized with services");
    } catch (error) {
      this.logger.error("Failed to initialize RunCommand", error);
      throw new CLIError(
        "Failed to initialize Claude services",
        "RUN_COMMAND_INIT_FAILED"
      );
    }
  }

  async execute(args = []) {
    await this.initialize();

    // Parse arguments - support both direct args and -- separator
    const parsedArgs = this.parseArguments(args);

    try {
      // Get default profile if no profile specified
      let profile = parsedArgs.profile;
      if (!profile) {
        const defaultProfileId = await this.configService.getDefaultProfileId();
        if (defaultProfileId) {
          profile = defaultProfileId;
        }
      }

      // Load profile configuration if available
      let profileConfig = null;
      if (profile) {
        profileConfig = await this.configService.loadProfile(profile);
        if (!profileConfig) {
          console.log(chalk.yellow(`⚠️  Profile '${profile}' not found, using defaults`));
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

      // Build Claude execution options
      const claudeOptions = {
        workingDirectory,
        dangerousSkip: parsedArgs.dangerousSkip,
        debug: parsedArgs.debug,
        project: parsedArgs.projectName,
        interactive: parsedArgs.interactive,
        timeout: parsedArgs.timeout,
        additionalArgs: parsedArgs.passThroughArgs
      };

      // Apply profile configuration
      if (profileConfig) {
        console.log(chalk.blue(`🔧 Using profile: ${profile}`));

        // Apply environment variables from profile
        if (profileConfig.config?.environment) {
          claudeOptions.environment = {
            ...claudeOptions.environment,
            ...profileConfig.config.environment
          };
        }

        // Apply Claude-specific settings from profile
        if (profileConfig.config?.claude) {
          const claudeProfile = profileConfig.config.claude;
          if (claudeProfile.model) {
            claudeOptions.model = claudeProfile.model;
          }
          if (claudeProfile.smallFastModel) {
            claudeOptions.smallFastModel = claudeProfile.smallFastModel;
          }
          if (claudeProfile.baseUrl) {
            claudeOptions.anthropicBaseUrl = claudeProfile.baseUrl;
          }
          if (claudeProfile.authToken) {
            claudeOptions.authToken = claudeProfile.authToken;
          }
          if (claudeProfile.configDirectory) {
            claudeOptions.configDirectory = claudeProfile.configDirectory;
          }
        }
      }

      // Show execution details
      console.log(chalk.blue.bold("🚀 Starting Claude Code"));
      console.log();

      if (profile) {
        console.log(chalk.gray(`Profile: ${profile}`));
      }

      if (parsedArgs.dangerousSkip) {
        console.log(chalk.yellow("⚠️  Dangerous skip enabled"));
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
      console.log();

      // Debug: Show what we're about to execute
      console.log(chalk.gray(`🔧 Executing: claude ${claudeOptions.additionalArgs?.join(' ') || ''}`));
      console.log(chalk.gray(`📁 Working directory: ${workingDirectory}`));
      console.log(chalk.gray(`⚙️  Interactive mode: ${claudeOptions.interactive !== false ? 'enabled' : 'disabled'}`));
      if (claudeOptions.configDirectory) {
        console.log(chalk.gray(`🗂️  CLAUDE_CONFIG_DIR: ${claudeOptions.configDirectory}`));
      }
      console.log();

      // For now, use direct spawning instead of the service to avoid hanging
      const { spawn } = await import("child_process");

      // Configure environment
      const env = {
        ...process.env
      };

      if (claudeOptions.configDirectory) {
        env.CLAUDE_CONFIG_DIR = claudeOptions.configDirectory;
      }

      console.log(chalk.blue("🚀 Launching Claude Code..."));

      // Build arguments
      const args = [];
      if (claudeOptions.dangerousSkip) {
        args.push("--dangerously-skip-permissions");
      }
      if (claudeOptions.debug) {
        args.push("--debug");
      }
      if (claudeOptions.project) {
        args.push("--project", claudeOptions.project);
      }
      if (claudeOptions.additionalArgs) {
        args.push(...claudeOptions.additionalArgs);
      }

      // Execute Claude directly
      const child = spawn("claude", args, {
        stdio: "inherit",
        env,
        cwd: claudeOptions.workingDirectory
      });

      return new Promise((resolve) => {
        child.on("exit", (code, signal) => {
          const result = {
            exitCode: code || 0,
            success: (code || 0) === 0,
            signal,
            executionTime: Date.now() - Date.now()
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

  /**
   * Parse command line arguments
   */
  parseArguments(args) {
    const result = {
      dangerousSkip: false,
      debug: false,
      profile: null,
      projectName: null,
      workingDirectory: null,
      timeout: null,
      interactive: true,
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
          } else {
            throw new CLIError("--config requires a profile name", "MISSING_CONFIG_VALUE");
          }
          break;

        case "-p":
        case "--project":
          if (i + 1 < mcfArgs.length) {
            result.projectName = mcfArgs[++i];
          } else {
            throw new CLIError("--project requires a project name", "MISSING_PROJECT_VALUE");
          }
          break;

        case "-w":
        case "--working-directory":
          if (i + 1 < mcfArgs.length) {
            result.workingDirectory = mcfArgs[++i];
          } else {
            throw new CLIError("--working-directory requires a path", "MISSING_WORKDIR_VALUE");
          }
          break;

        case "-t":
        case "--timeout":
          if (i + 1 < mcfArgs.length) {
            const timeoutValue = parseInt(mcfArgs[++i]);
            if (isNaN(timeoutValue) || timeoutValue <= 0) {
              throw new CLIError("--timeout requires a positive number (milliseconds)", "INVALID_TIMEOUT");
            }
            result.timeout = timeoutValue;
          } else {
            throw new CLIError("--timeout requires a value", "MISSING_TIMEOUT_VALUE");
          }
          break;

        case "--no-interactive":
          result.interactive = false;
          break;

        case "--dangerous-skip":
          result.dangerousSkip = true;
          break;

        default:
          // Unknown argument
          if (arg.startsWith("-")) {
            console.log(chalk.yellow(`⚠️  Unknown option: ${arg}`));
          } else {
            // Treat as positional argument (could be profile name)
            if (!result.profile && mcfArgs.length === 1) {
              result.profile = arg;
            } else {
              console.log(chalk.yellow(`⚠️  Ignoring unexpected argument: ${arg}`));
            }
          }
      }
    }

    result.passThroughArgs = passThrough;
    return result;
  }

  showHelp() {
    console.log(chalk.blue("MCF Run Command Help"));
    console.log(chalk.gray("─".repeat(50)));
    console.log();
    console.log("Execute Claude Code with enhanced configuration and flags");
    console.log();
    console.log(chalk.blue("Options:"));
    console.log();
    console.log("  -d, --debug                    Enable debug mode");
    console.log("  -c, --config <profile>         Use specific configuration profile");
    console.log("  -p, --project <name>           Set Claude project name");
    console.log("  -w, --working-directory <path> Set working directory");
    console.log("  -t, --timeout <ms>            Set execution timeout (milliseconds)");
    console.log("  --no-interactive               Run in non-interactive mode");
    console.log("  --dangerous-skip               Skip permission checks (dangerous)");
    console.log("  -- <args...>                   Pass remaining arguments to Claude");
    console.log();
    console.log(chalk.blue("Examples:"));
    console.log();
    console.log("  mcf run                           # Start Claude with default settings");
    console.log("  mcf run --debug                   # Start Claude in debug mode");
    console.log("  mcf run --config myprofile        # Use specific profile");
    console.log("  mcf run --project myproject       # Set Claude project");
    console.log("  mcf run -- --help                 # Pass --help to Claude directly");
    console.log("  mcf run --debug -- --verbose      # Debug mode + pass --verbose to Claude");
    console.log();
    console.log(chalk.blue("Profile Configuration:"));
    console.log();
    console.log("Profiles can be configured with:");
    console.log("  • Environment variables (ANTHROPIC_BASE_URL, etc.)");
    console.log("  • Model settings (claude-3-5-sonnet, etc.)");
    console.log("  • Authentication tokens");
    console.log("  • Custom working directories");
    console.log();
    console.log("Use 'mcf config create <name>' to create profiles.");
    console.log();
    return { success: true };
  }

  getMetadata() {
    return RunCommand.metadata;
  }
}
