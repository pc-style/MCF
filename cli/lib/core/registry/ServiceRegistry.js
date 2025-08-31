import { LoggerFactory } from "../logging/LoggerFactory.js";

/**
 * Service factory type definition
 * @template T
 * @typedef {function(): T} ServiceFactory
 */

/**
 * @typedef {Object} CLIConfig
 * @property {boolean} [verbose]
 * @property {boolean} [debug]
 */

/**
 * ServiceRegistry for MCF CLI
 * Manages service registration, dependency injection, and lifecycle
 */
export class ServiceRegistry {
  static instance;

  /** @type {Map<string, any>} */
  services = new Map();

  /** @type {Map<string, ServiceFactory>} */
  factories = new Map();

  /** @type {Object} */
  logger = LoggerFactory.getLogger("ServiceRegistry");

  /** @type {CLIConfig} */
  config = {};

  /**
   * Private constructor to enforce singleton pattern
   */
  constructor() {}

  /**
   * Get singleton instance of ServiceRegistry
   * @returns {ServiceRegistry}
   */
  static getInstance() {
    if (!ServiceRegistry.instance) {
      ServiceRegistry.instance = new ServiceRegistry();
    }
    return ServiceRegistry.instance;
  }

  /**
   * Initialize ServiceRegistry with optional configuration
   * @param {CLIConfig} [config] Optional CLI configuration
   */
  initialize(config) {
    this.config = config || {};
    this.logger.info("ServiceRegistry initialized", { config: this.config });
  }

  /**
   * Register a service factory
   * @param {string} key Unique service identifier
   * @param {ServiceFactory} factory Service factory function
   */
  registerService(key, factory) {
    if (this.factories.has(key)) {
      this.logger.warn(`Service ${key} already registered. Overwriting.`);
    }
    this.factories.set(key, factory);
    this.logger.debug(`Registered service: ${key}`);
  }

  /**
   * Get a service instance, creating it if not already instantiated
   * @param {string} key Service identifier
   * @returns {any} Service instance
   */
  getService(key) {
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
   * @param {string} key Service identifier
   * @returns {boolean} Boolean indicating service registration status
   */
  hasService(key) {
    return this.factories.has(key);
  }

  /**
   * Clear all registered services
   */
  reset() {
    this.services.clear();
    this.factories.clear();
    this.logger.info("ServiceRegistry reset");
  }

  /**
   * Get current configuration
   * @returns {CLIConfig} Current CLI configuration
   */
  getConfig() {
    return { ...this.config };
  }
}

// Export singleton instance for convenient access
export const serviceRegistry = ServiceRegistry.getInstance();
