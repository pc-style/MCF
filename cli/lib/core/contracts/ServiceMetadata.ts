import { MCFProfile } from "../../types/MCFProfile.js";

export interface ServiceMetadata<ConfigType = Record<string, unknown>> {
  // Universal Core Metadata (UCM) Fields
  id: string;
  name: string;
  version: string;
  description?: string;

  // CLI-specific Configuration
  config?: {
    profile?: MCFProfile;
    options?: ConfigType;
  };

  // Service Dependencies
  dependencies?: {
    services?: string[];
    external?: string[];
  };

  // CLI Runtime Information
  runtime?: {
    environment?: string;
    requiredPermissions?: string[];
  };

  // Documentation and Help
  helpText?: string;
  documentation?: string;
}
