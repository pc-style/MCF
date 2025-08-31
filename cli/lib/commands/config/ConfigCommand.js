import { BaseCommand } from "../../core/interfaces/BaseCommand.js";
import { ServiceRegistry } from "../../core/registry/ServiceRegistry.js";
import { ILogger } from "../../core/logging/ILogger.js";
import { LoggerFactory } from "../../core/logging/LoggerFactory.js";
import { CLIError } from "../../types/CLITypes.js";
import chalk from "chalk";

/**
 * ConfigCommand - MCF CLI Configuration Profile Management
 * Provides subcommands for managing MCF configuration profiles
 */
export class ConfigCommand extends BaseCommand {
  constructor(serviceRegistry) {
    super();
    this.serviceRegistry = serviceRegistry;
    this.logger = LoggerFactory.getLogger("ConfigCommand");
    this.configService = null;
    this.fileSystemService = null;
  }

  static get metadata() {
    return {
      name: "ConfigCommand",
      description: "Manage MCF configuration profiles",
      category: "config",
      version: "1.0.0",
      dependencies: {
        services: ["IConfigurationService", "IFileSystemService"],
        commands: [],
        external: []
      }
    };
  }

  async initialize() {
    try {
      this.configService = this.serviceRegistry.getService("IConfigurationService");
      this.fileSystemService = this.serviceRegistry.getService("IFileSystemService");
      this.logger.debug("ConfigCommand initialized with services");
    } catch (error) {
      this.logger.error("Failed to initialize ConfigCommand", error);
      throw new CLIError(
        "Failed to initialize configuration services",
        "CONFIG_COMMAND_INIT_FAILED"
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
      case "set-default":
      case "default":
        return await this.setDefaultProfile(subArgs[0]);
      case "edit":
        return await this.editProfile(subArgs[0]);
      case "clone":
      case "copy":
        return await this.cloneProfile(subArgs);
      case "validate":
        return await this.validateProfile(subArgs[0]);
      default:
        console.log(chalk.red(`Unknown config subcommand: ${subcommand}`));
        return this.showHelp();
    }
  }

  async listProfiles() {
    try {
      const profileIds = await this.configService.listProfiles();
      const defaultProfileId = await this.configService.getDefaultProfileId();

      if (profileIds.length === 0) {
        console.log(chalk.yellow("No configuration profiles found."));
        console.log("Create one with: mcf config create <name> [environment]");
        return;
      }

      console.log(chalk.blue("MCF Configuration Profiles:"));
      console.log();

      for (const profileId of profileIds) {
        const marker = profileId === defaultProfileId ? chalk.green(" (default)") : "";
        console.log(`  ${chalk.cyan(profileId)}${marker}`);
      }

      console.log();
      console.log(`Total: ${profileIds.length} profile(s)`);
    } catch (error) {
      console.error(chalk.red(`Failed to list profiles: ${error.message}`));
      throw error;
    }
  }

  async showProfile(profileId) {
    if (!profileId) {
      console.log(chalk.red("Profile ID is required"));
      console.log("Usage: mcf config show <profile-id>");
      return;
    }

    try {
      const profile = await this.configService.loadProfile(profileId);

      if (!profile) {
        console.log(chalk.red(`Profile '${profileId}' not found`));
        return;
      }

      console.log(chalk.blue(`Profile: ${profileId}`));
      console.log(chalk.gray("─".repeat(50)));
      console.log(`Name: ${chalk.cyan(profile.name)}`);
      console.log(`Environment: ${chalk.cyan(profile.environment)}`);

      if (profile.description) {
        console.log(`Description: ${profile.description}`);
      }

      if (profile.version) {
        console.log(`Version: ${profile.version}`);
      }

      if (profile.lastUpdated) {
        console.log(`Last Updated: ${new Date(profile.lastUpdated).toLocaleString()}`);
      }

      console.log();
      console.log(chalk.blue("Configuration:"));

      if (profile.config) {
        console.log(`  Timeout: ${profile.config.timeout || "default"}ms`);
        console.log(`  Max Retries: ${profile.config.maxRetries || "default"}`);
        console.log(`  Log Level: ${profile.config.logLevel || "default"}`);

        // Display Claude configuration
        if (profile.config.claude) {
          console.log();
          console.log(chalk.blue("Claude Configuration:"));
          if (profile.config.claude.configDirectory) {
            console.log(`  Config Directory: ${chalk.green(profile.config.claude.configDirectory)}`);
          }
          if (profile.config.claude.model) {
            console.log(`  Model: ${chalk.green(profile.config.claude.model)}`);
          }
          if (profile.config.claude.baseUrl) {
            console.log(`  Base URL: ${chalk.green(profile.config.claude.baseUrl)}`);
          }
          if (profile.config.claude.authToken) {
            console.log(`  Auth Token: ${chalk.green("***configured***")}`);
          }
        }
      }

      if (profile.permissions) {
        console.log();
        console.log(chalk.blue("Permissions:"));

        if (profile.permissions.allowedServices) {
          console.log(`  Allowed Services: ${profile.permissions.allowedServices.join(", ")}`);
        }

        if (profile.permissions.blockedServices) {
          console.log(`  Blocked Services: ${profile.permissions.blockedServices.join(", ")}`);
        }
      }
    } catch (error) {
      console.error(chalk.red(`Failed to show profile: ${error.message}`));
      throw error;
    }
  }

  async createProfile(args) {
    const [name, environment = "development"] = args;

    if (!name) {
      console.log(chalk.red("Profile name is required"));
      console.log("Usage: mcf config create <name> [environment]");
      console.log("Environments: development, production, staging, test");
      return;
    }

    try {
      // Validate environment
      const validEnvs = ["development", "production", "staging", "test"];
      if (!validEnvs.includes(environment)) {
        console.log(chalk.red(`Invalid environment: ${environment}`));
        console.log(`Valid environments: ${validEnvs.join(", ")}`);
        return;
      }

      // Check if profile already exists
      const profileId = this.generateProfileId(name);
      if (await this.configService.profileExists(profileId)) {
        console.log(chalk.red(`Profile '${profileId}' already exists`));
        return;
      }

      const profile = await this.configService.createProfile(name, environment);
      console.log(chalk.green(`Profile '${profile.id}' created successfully`));
      console.log(`Environment: ${environment}`);
      console.log();
      console.log("You can now:");
      console.log(`  • Set as default: mcf config set-default ${profile.id}`);
      console.log(`  • Edit settings: mcf config edit ${profile.id}`);
      console.log(`  • Show details: mcf config show ${profile.id}`);
    } catch (error) {
      console.error(chalk.red(`Failed to create profile: ${error.message}`));
      throw error;
    }
  }

  async deleteProfile(profileId) {
    if (!profileId) {
      console.log(chalk.red("Profile ID is required"));
      console.log("Usage: mcf config delete <profile-id>");
      return;
    }

    try {
      // Check if profile exists
      if (!(await this.configService.profileExists(profileId))) {
        console.log(chalk.red(`Profile '${profileId}' not found`));
        return;
      }

      // Check if it's the default profile
      const defaultProfileId = await this.configService.getDefaultProfileId();
      if (profileId === defaultProfileId) {
        console.log(chalk.red(`Cannot delete default profile '${profileId}'`));
        console.log("Set a different default first: mcf config set-default <other-profile>");
        return;
      }

      const success = await this.configService.deleteProfile(profileId);
      if (success) {
        console.log(chalk.green(`Profile '${profileId}' deleted successfully`));
      }
    } catch (error) {
      console.error(chalk.red(`Failed to delete profile: ${error.message}`));
      throw error;
    }
  }

  async setDefaultProfile(profileId) {
    if (!profileId) {
      console.log(chalk.red("Profile ID is required"));
      console.log("Usage: mcf config set-default <profile-id>");
      return;
    }

    try {
      // Check if profile exists
      if (!(await this.configService.profileExists(profileId))) {
        console.log(chalk.red(`Profile '${profileId}' not found`));
        return;
      }

      await this.configService.setDefaultProfile(profileId);
      console.log(chalk.green(`Default profile set to '${profileId}'`));
    } catch (error) {
      console.error(chalk.red(`Failed to set default profile: ${error.message}`));
      throw error;
    }
  }

  async editProfile(profileId) {
    if (!profileId) {
      console.log(chalk.red("Profile ID is required"));
      console.log("Usage: mcf config edit <profile-id>");
      return;
    }

    console.log(chalk.yellow("Profile editing not yet implemented"));
    console.log("For now, you can manually edit the JSON files in ~/.mcf/profiles/");
    console.log(`Profile location: ~/.mcf/profiles/${profileId}.json`);
  }

  async cloneProfile(args) {
    const [sourceProfileId, newProfileId] = args;

    if (!sourceProfileId || !newProfileId) {
      console.log(chalk.red("Both source and new profile IDs are required"));
      console.log("Usage: mcf config clone <source-profile> <new-profile>");
      return;
    }

    try {
      // Check if source profile exists
      if (!(await this.configService.profileExists(sourceProfileId))) {
        console.log(chalk.red(`Source profile '${sourceProfileId}' not found`));
        return;
      }

      // Check if new profile already exists
      if (await this.configService.profileExists(newProfileId)) {
        console.log(chalk.red(`Profile '${newProfileId}' already exists`));
        return;
      }

      const clonedProfile = await this.configService.cloneProfile(sourceProfileId, newProfileId);
      console.log(chalk.green(`Profile '${sourceProfileId}' cloned to '${newProfileId}'`));
      console.log(`New profile ID: ${clonedProfile.id}`);
    } catch (error) {
      console.error(chalk.red(`Failed to clone profile: ${error.message}`));
      throw error;
    }
  }

  async validateProfile(profileId) {
    if (!profileId) {
      console.log(chalk.red("Profile ID is required"));
      console.log("Usage: mcf config validate <profile-id>");
      return;
    }

    try {
      const profile = await this.configService.loadProfile(profileId);

      if (!profile) {
        console.log(chalk.red(`Profile '${profileId}' not found`));
        return;
      }

      const validation = await this.configService.validateProfile(profile);

      if (validation.isValid) {
        console.log(chalk.green(`✓ Profile '${profileId}' is valid`));

        if (validation.warnings.length > 0) {
          console.log(chalk.yellow("Warnings:"));
          validation.warnings.forEach(warning => {
            console.log(`  • ${warning}`);
          });
        }
      } else {
        console.log(chalk.red(`✗ Profile '${profileId}' has validation errors:`));
        validation.errors.forEach(error => {
          console.log(`  • ${error}`);
        });
      }
    } catch (error) {
      console.error(chalk.red(`Failed to validate profile: ${error.message}`));
      throw error;
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
    console.log("  set-default, default <id>  Set default profile");
    console.log("  edit <id>                  Edit profile (not implemented)");
    console.log("  clone, copy <src> <new>    Clone profile");
    console.log("  validate <id>              Validate profile");
    console.log();
    console.log(chalk.blue("Examples:"));
    console.log();
    console.log("  mcf config list");
    console.log("  mcf config create myapp production");
    console.log("  mcf config show myapp-production");
    console.log("  mcf config set-default myapp-production");
    console.log("  mcf config clone myapp-prod myapp-staging");
    console.log();
    return { success: true };
  }

  generateProfileId(name) {
    return name.toLowerCase().replace(/[^a-z0-9]/g, "-").replace(/-+/g, "-").replace(/^-|-$/g, "");
  }

  getMetadata() {
    return ConfigCommand.metadata;
  }
}
