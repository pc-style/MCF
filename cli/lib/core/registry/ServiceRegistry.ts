import { LoggerFactory } from "../logger/LoggerFactory.js";

/**
 * CLIConfig interface for MCF CLI specific configuration
 */
export interface CLIConfig {
  verbose?: boolean;
  debug?: boolean;
}

/**
 * Service factory type definition
 */
export type ServiceFactory<T> = () => T;

/**
 * ServiceRegistry for MCF CLI
 * Manages service registration, dependency injection, and lifecycle
 */
export class ServiceRegistry {
  private static instance: ServiceRegistry;
  private services: Map<string, any> = new Map();
  private factories: Map<string, ServiceFactory<any>> = new Map();
  private logger = LoggerFactory.getLogger("ServiceRegistry");
  private config: CLIConfig = {};

  /**
   * Private constructor to enforce singleton pattern
   */
  private constructor() {}

  /**
   * Get singleton instance of ServiceRegistry
   */
  public static getInstance(): ServiceRegistry {
    if (!ServiceRegistry.instance) {
      ServiceRegistry.instance = new ServiceRegistry();
    }
    return ServiceRegistry.instance;
  }

  /**
   * Initialize ServiceRegistry with optional configuration
   * @param config Optional CLI configuration
   */
  public initialize(config?: CLIConfig): void {
    this.config = config || {};
    this.logger.info("ServiceRegistry initialized", { config: this.config });
  }

  /**
   * Register a service factory
   * @param key Unique service identifier
   * @param factory Service factory function
   */
  public registerService<T>(key: string, factory: ServiceFactory<T>): void {
    if (this.factories.has(key)) {
      this.logger.warn(`Service ${key} already registered. Overwriting.`);
    }
    this.factories.set(key, factory);
    this.logger.debug(`Registered service: ${key}`);
  }

  /**
   * Get a service instance, creating it if not already instantiated
   * @param key Service identifier
   * @returns Service instance
   */
  public getService<T>(key: string): T {
    // If service already instantiated, return it
    if (this.services.has(key)) {
      return this.services.get(key);
    }

    // If factory exists, create and cache service
    const factory = this.factories.get(key);
    if (!factory) {
      throw new Error(`Service not found: ${key}`);
    }

    const service = factory();
    this.services.set(key, service);
    return service;
  }

  /**
   * Check if a service is registered
   * @param key Service identifier
   * @returns Boolean indicating service registration status
   */
  public hasService(key: string): boolean {
    return this.factories.has(key);
  }

  /**
   * Clear all registered services
   */
  public reset(): void {
    this.services.clear();
    this.factories.clear();
    this.logger.info("ServiceRegistry reset");
  }

  /**
   * Get current configuration
   * @returns Current CLI configuration
   */
  public getConfig(): CLIConfig {
    return { ...this.config };
  }
}

// Export singleton instance for convenient access
export const serviceRegistry = ServiceRegistry.getInstance();
