// BaseCommand should be a class with overridable methods
export class BaseCommand {
  constructor(metadata) {
    this.metadata = metadata;
  }

  // Default implementation for execute method
  async execute(input) {
    throw new Error("Method 'execute' must be implemented");
  }

  // Default implementation for validate method
  async validate(input) {
    return true;
  }

  // Default implementation for help method
  help() {
    return "No help text available";
  }

  // Default implementation for description method
  description() {
    return this.metadata.description || "No description available";
  }

  // Default implementation for getMetadata method
  getMetadata() {
    return { ...this.metadata };
  }
}
