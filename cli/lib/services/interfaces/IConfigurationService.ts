import { MCFProfile } from "../../types/MCFProfile.js";
import { CLIEnvironment } from "../../types/CLITypes.js";

/**
 * Configuration service interface for MCF CLI
 * Handles profile management, storage, and validation
 */
export interface IConfigurationService {
  /**
   * Save a profile to storage
   * @param profile Profile to save
   */
  saveProfile(profile: MCFProfile): Promise<void>;

  /**
   * Load a profile from storage
   * @param profileId Profile identifier
   * @returns Profile instance or null if not found
   */
  loadProfile(profileId: string): Promise<MCFProfile | null>;

  /**
   * List all available profiles
   * @returns Array of profile identifiers
   */
  listProfiles(): Promise<string[]>;

  /**
   * Delete a profile from storage
   * @param profileId Profile identifier
   */
  deleteProfile(profileId: string): Promise<boolean>;

  /**
   * Check if a profile exists
   * @param profileId Profile identifier
   */
  profileExists(profileId: string): Promise<boolean>;

  /**
   * Get the default profile
   */
  getDefaultProfile(): Promise<MCFProfile>;

  /**
   * Set the default profile
   * @param profileId Profile identifier
   */
  setDefaultProfile(profileId: string): Promise<void>;

  /**
   * Get the current default profile ID
   */
  getDefaultProfileId(): Promise<string | null>;

  /**
   * Validate a profile structure
   * @param profile Profile to validate
   * @returns Validation result
   */
  validateProfile(profile: MCFProfile): Promise<ProfileValidationResult>;

  /**
   * Create a new profile with default values
   * @param name Profile name
   * @param environment Target environment
   */
  createProfile(name: string, environment: CLIEnvironment): Promise<MCFProfile>;

  /**
   * Clone an existing profile
   * @param sourceProfileId Source profile identifier
   * @param newProfileId New profile identifier
   */
  cloneProfile(sourceProfileId: string, newProfileId: string): Promise<MCFProfile>;
}

/**
 * Profile validation result
 */
export interface ProfileValidationResult {
  isValid: boolean;
  errors: string[];
  warnings: string[];
}

/**
 * Configuration service configuration
 */
export interface ConfigurationServiceConfig {
  configDirectory: string;
  profilesDirectory?: string;
  defaultProfileName?: string;
  validateProfiles?: boolean;
  maxProfileSize?: number; // in bytes
}



