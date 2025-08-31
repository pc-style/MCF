import { ServiceMetadata } from "../contracts/ServiceMetadata.js";

export abstract class BaseService {
  protected readonly metadata: ServiceMetadata;

  constructor(metadata: ServiceMetadata) {
    this.metadata = {
      id: metadata.id,
      name: metadata.name,
      version: metadata.version,
      description: metadata.description || "",
    };
  }

  public getMetadata(): ServiceMetadata {
    return { ...this.metadata };
  }

  public abstract initialize(): Promise<void>;
  public abstract dispose(): Promise<void>;

  protected validateInitialization(): void {
    if (!this.metadata.id) {
      throw new Error("Service initialization failed: Missing service ID");
    }
  }

  protected logServiceEvent(event: string, details?: unknown): void {
    console.debug(`[${this.metadata.name}] ${event}`, details || "");
  }
}
