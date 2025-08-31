import { BaseService } from "../base/BaseService.js";
import { ServiceRegistry } from "./ServiceRegistry.js";
import { BaseCommand } from "../interfaces/BaseCommand.js";
import { CommandMetadata } from "../contracts/CommandMetadata.js";
import { LoggerFactory } from "../logging/LoggerFactory.js";
import { ILogger } from "../logging/ILogger.js";
import { fileURLToPath } from "url";
import path from "path";
import fs from "fs/promises";

export class CommandRegistry extends BaseService {
  private static instance: CommandRegistry;
  private commands: Map<string, () => Promise<BaseCommand>> = new Map();
  private logger: ILogger;

  private constructor(private serviceRegistry: ServiceRegistry) {
    super();
    this.logger = LoggerFactory.getLogger(CommandRegistry.name);
  }

  public static getInstance(serviceRegistry: ServiceRegistry): CommandRegistry {
    if (!CommandRegistry.instance) {
      CommandRegistry.instance = new CommandRegistry(serviceRegistry);
    }
    return CommandRegistry.instance;
  }

  public async loadCommands(): Promise<void> {
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

  private async registerCommand(
    category: string,
    commandName: string,
  ): Promise<void> {
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

  public async getCommand(
    category: string,
    commandName: string,
  ): Promise<BaseCommand | null> {
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

  public getRegisteredCommands(): string[] {
    return Array.from(this.commands.keys());
  }

  protected async onInit(): Promise<void> {
    await this.loadCommands();
  }
}
