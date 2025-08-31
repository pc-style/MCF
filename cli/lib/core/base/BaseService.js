import { ServiceMetadata } from "../contracts/ServiceMetadata.js";
import crypto from "crypto";

export class BaseService {
  constructor(metadata = {}) {
    const serviceClass = this.constructor.name;
    this.metadata = {
      id: metadata.id || crypto.randomUUID(),
      name: metadata.name || serviceClass,
      version: metadata.version || "1.0.0",
      description: metadata.description || `${serviceClass} service`,
    };
  }

  getMetadata() {
    return { ...this.metadata };
  }

  initialize() {
    throw new Error("Method 'initialize' must be implemented by subclass");
  }

  dispose() {
    throw new Error("Method 'dispose' must be implemented by subclass");
  }

  validateInitialization() {
    if (!this.metadata.id) {
      throw new Error("Service initialization failed: Missing service ID");
    }
  }

  logServiceEvent(event, details) {
    console.debug(`[${this.metadata.name}] ${event}`, details || "");
  }
}
