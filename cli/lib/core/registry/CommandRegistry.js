import { BaseService } from "../base/BaseService.js";
import { ServiceRegistry } from "./ServiceRegistry.js";
import { BaseCommand } from "../interfaces/BaseCommand.js";
import { LoggerFactory } from "../logging/LoggerFactory.js";
import { fileURLToPath } from "url";
import path from "path";
import fs from "fs/promises";

export class CommandRegistry extends BaseService {
  /** @type {CommandRegistry} */
  static instance;

  /** @type {Map<string, () => Promise<BaseCommand>>} */
  commands = new Map();

  /** @type {Object} */
  logger;

  /**
   * @param {ServiceRegistry} serviceRegistry
   */
  constructor(serviceRegistry) {
    super();
    this.serviceRegistry = serviceRegistry;
    this.logger = LoggerFactory.getLogger(CommandRegistry.name);
  }

  /**
   * Get singleton instance
   * @param {ServiceRegistry} serviceRegistry
   * @returns {CommandRegistry}
   */
  static getInstance(serviceRegistry) {
    if (!CommandRegistry.instance) {
      CommandRegistry.instance = new CommandRegistry(serviceRegistry);
    }
    return CommandRegistry.instance;
  }

  /**
   * Load commands from commands directory
   * @returns {Promise<void>}
   */
  async loadCommands() {
    const __filename = fileURLToPath(import.meta.url);
    const __dirname = path.dirname(__filename);
    const commandsBasePath = path.resolve(__dirname, "../../commands");

    try {
      const commandCategories = await fs.readdir(commandsBasePath);

      for (const category of commandCategories) {
        const categoryPath = path.join(commandsBasePath, category);
        const commandFiles = await fs.readdir(categoryPath);

        for (const commandFile of commandFiles) {
          if (commandFile.endsWith(".ts") || commandFile.endsWith(".js")) {
            const commandName = path.basename(
              commandFile,
              path.extname(commandFile),
            );
            await this.registerCommand(category, commandName);
          }
        }
      }
    } catch (error) {
      this.logger.error(
        `Error loading commands: ${error instanceof Error ? error.message : error}`,
      );
    }
  }

  /**
   * Register a command
   * @param {string} category
   * @param {string} commandName
   * @returns {Promise<void>}
   */
  async registerCommand(category, commandName) {
    const fullCommandPath = `../../commands/${category}/${commandName}.js`;

    const lazyLoadCommand = async () => {
      try {
        const commandModule = await import(fullCommandPath);
        const CommandClass = commandModule.default;

        // Dependency injection for command
        const commandInstance = new CommandClass(this.serviceRegistry);

        return commandInstance;
      } catch (error) {
        this.logger.error(
          `Failed to load command ${category}/${commandName}: ${error instanceof Error ? error.message : error}`,
        );
        throw error;
      }
    };

    // Use category/commandName as the key for command registration
    const commandKey = `${category}/${commandName}`;
    this.commands.set(commandKey, lazyLoadCommand);
    this.logger.info(`Registered command: ${commandKey}`);
  }

  /**
   * Get a command instance
   * @param {string} category
   * @param {string} commandName
   * @returns {Promise<BaseCommand | null>}
   */
  async getCommand(category, commandName) {
    const commandKey = `${category}/${commandName}`;
    const commandLoader = this.commands.get(commandKey);

    if (!commandLoader) {
      this.logger.warn(`Command not found: ${commandKey}`);
      return null;
    }

    try {
      return await commandLoader();
    } catch (error) {
      this.logger.error(
        `Error instantiating command ${commandKey}: ${error instanceof Error ? error.message : error}`,
      );
      return null;
    }
  }

  /**
   * Get list of registered commands
   * @returns {string[]}
   */
  getRegisteredCommands() {
    return Array.from(this.commands.keys());
  }

  /**
   * Initialize commands
   * @returns {Promise<void>}
   */
  async onInit() {
    await this.loadCommands();
  }
}
