import { BaseService } from "../../core/base/BaseService.js";
import { ILogger } from "../../core/logging/ILogger.js";
import { LoggerFactory } from "../../core/logging/LoggerFactory.js";
import { CLIError } from "../../types/CLITypes.js";
import fs from "fs/promises";
import path from "path";

/**
 * Configuration service implementation for MCF CLI
 * Handles profile management with JSON-based storage and validation
 */
export class ConfigurationService extends BaseService {
  constructor(config, logger) {
    super();
    this.config = config;
    this.logger = logger || LoggerFactory.getLogger("ConfigurationService");
    this.profilesDir = config.profilesDirectory || path.join(config.configDirectory, "profiles");
    this.defaultProfilePath = path.join(config.configDirectory, "default-profile.json");
  }

  /**
   * Save a profile to storage
   */
  async saveProfile(profile) {
    try {
      // Validate profile before saving
      if (this.config.validateProfiles !== false) {
        const validation = await this.validateProfile(profile);
        if (!validation.isValid) {
          throw new CLIError(
            `Profile validation failed: ${validation.errors.join(", ")}`,
            "PROFILE_VALIDATION_FAILED",
            { errors: validation.errors, warnings: validation.warnings }
          );
        }

        if (validation.warnings.length > 0) {
          this.logger.warn(`Profile warnings: ${validation.warnings.join(", ")}`);
        }
      }

      // Ensure profiles directory exists
      await this.ensureProfilesDirectory();

      // Generate filename from profile ID
      const profilePath = this.getProfilePath(profile.id);

      // Write profile to file
      const profileData = JSON.stringify(profile, null, 2);
      await fs.writeFile(profilePath, profileData, "utf-8");

      this.logger.info(`Profile '${profile.id}' saved to ${profilePath}`);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to save profile '${profile.id}': ${message}`);
      throw new CLIError(
        `Failed to save profile: ${message}`,
        "PROFILE_SAVE_FAILED",
        { profileId: profile.id }
      );
    }
  }

  /**
   * Load a profile from storage
   */
  async loadProfile(profileId) {
    try {
      const profilePath = this.getProfilePath(profileId);

      // Check if profile exists
      try {
        await fs.access(profilePath);
      } catch {
        this.logger.debug(`Profile '${profileId}' not found at ${profilePath}`);
        return null;
      }

      // Read and parse profile
      const profileData = await fs.readFile(profilePath, "utf-8");
      const profile = JSON.parse(profileData);

      // Validate loaded profile
      if (profile.id !== profileId) {
        throw new CLIError(
          `Profile ID mismatch: expected '${profileId}', got '${profile.id}'`,
          "PROFILE_ID_MISMATCH"
        );
      }

      this.logger.debug(`Profile '${profileId}' loaded successfully`);
      return profile;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to load profile '${profileId}': ${message}`);
      throw new CLIError(
        `Failed to load profile: ${message}`,
        "PROFILE_LOAD_FAILED",
        { profileId }
      );
    }
  }

  /**
   * List all available profiles
   */
  async listProfiles() {
    try {
      await this.ensureProfilesDirectory();

      const entries = await fs.readdir(this.profilesDir);
      const profileIds = [];

      for (const entry of entries) {
        if (entry.endsWith(".json")) {
          const profileId = entry.replace(".json", "");
          profileIds.push(profileId);
        }
      }

      this.logger.debug(`Found ${profileIds.length} profiles: ${profileIds.join(", ")}`);
      return profileIds;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to list profiles: ${message}`);
      throw new CLIError(
        `Failed to list profiles: ${message}`,
        "PROFILE_LIST_FAILED"
      );
    }
  }

  /**
   * Delete a profile from storage
   */
  async deleteProfile(profileId) {
    try {
      const profilePath = this.getProfilePath(profileId);

      // Check if profile exists
      try {
        await fs.access(profilePath);
      } catch {
        this.logger.warn(`Profile '${profileId}' not found, nothing to delete`);
        return false;
      }

      // Check if this is the default profile
      const defaultProfileId = await this.getDefaultProfileId();
      if (defaultProfileId === profileId) {
        throw new CLIError(
          `Cannot delete default profile '${profileId}'`,
          "CANNOT_DELETE_DEFAULT_PROFILE"
        );
      }

      // Delete the profile file
      await fs.unlink(profilePath);

      this.logger.info(`Profile '${profileId}' deleted successfully`);
      return true;
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to delete profile '${profileId}': ${message}`);
      throw new CLIError(
        `Failed to delete profile: ${message}`,
        "PROFILE_DELETE_FAILED",
        { profileId }
      );
    }
  }

  /**
   * Check if a profile exists
   */
  async profileExists(profileId) {
    try {
      const profilePath = this.getProfilePath(profileId);
      await fs.access(profilePath);
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Get the default profile
   */
  async getDefaultProfile() {
    try {
      const defaultProfileId = await this.getDefaultProfileId();

      if (!defaultProfileId) {
        // Create and return a default profile
        return await this.createDefaultProfile();
      }

      const profile = await this.loadProfile(defaultProfileId);
      if (!profile) {
        throw new CLIError(
          `Default profile '${defaultProfileId}' not found`,
          "DEFAULT_PROFILE_NOT_FOUND"
        );
      }

      return profile;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get default profile: ${message}`);
      throw new CLIError(
        `Failed to get default profile: ${message}`,
        "DEFAULT_PROFILE_LOAD_FAILED"
      );
    }
  }

  /**
   * Set the default profile
   */
  async setDefaultProfile(profileId) {
    try {
      // Verify profile exists
      if (!(await this.profileExists(profileId))) {
        throw new CLIError(
          `Profile '${profileId}' does not exist`,
          "PROFILE_NOT_FOUND"
        );
      }

      // Write default profile setting
      const defaultConfig = { defaultProfileId: profileId };
      await fs.writeFile(this.defaultProfilePath, JSON.stringify(defaultConfig, null, 2), "utf-8");

      this.logger.info(`Default profile set to '${profileId}'`);
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to set default profile '${profileId}': ${message}`);
      throw new CLIError(
        `Failed to set default profile: ${message}`,
        "DEFAULT_PROFILE_SET_FAILED",
        { profileId }
      );
    }
  }

  /**
   * Get the current default profile ID
   */
  async getDefaultProfileId() {
    try {
      const data = await fs.readFile(this.defaultProfilePath, "utf-8");
      const config = JSON.parse(data);
      return config.defaultProfileId || null;
    } catch {
      return null;
    }
  }

  /**
   * Validate a profile structure
   */
  async validateProfile(profile) {
    const errors = [];
    const warnings = [];

    // Validate required fields
    if (!profile.id || typeof profile.id !== "string" || profile.id.trim() === "") {
      errors.push("Profile ID is required and must be a non-empty string");
    }

    if (!profile.name || typeof profile.name !== "string" || profile.name.trim() === "") {
      errors.push("Profile name is required and must be a non-empty string");
    }

    if (!profile.environment || !["development", "production", "staging", "test"].includes(profile.environment)) {
      errors.push("Profile environment must be one of: development, production, staging, test");
    }

    // Validate config structure
    if (profile.config) {
      if (profile.config.timeout !== undefined && (typeof profile.config.timeout !== "number" || profile.config.timeout < 0)) {
        errors.push("Config timeout must be a positive number");
      }

      if (profile.config.maxRetries !== undefined && (typeof profile.config.maxRetries !== "number" || profile.config.maxRetries < 0)) {
        errors.push("Config maxRetries must be a non-negative number");
      }
    }

    // Validate permissions
    if (profile.permissions) {
      const validPermissions = ["read", "write", "execute", "admin", "network", "system"];
      if (profile.permissions.allowedServices) {
        for (const service of profile.permissions.allowedServices) {
          if (!validPermissions.includes(service)) {
            errors.push(`Invalid permission: ${service}`);
          }
        }
      }
    }

    return {
      isValid: errors.length === 0,
      errors,
      warnings
    };
  }

  /**
   * Create a new profile with default values
   */
  async createProfile(name, environment) {
    const profileId = this.generateProfileId(name);

    const profile = {
      id: profileId,
      name,
      description: `Profile for ${environment} environment`,
      environment,
      config: {
        timeout: 30000,
        maxRetries: 3,
        logLevel: "info"
      },
      version: "1.0.0",
      lastUpdated: new Date()
    };

    await this.saveProfile(profile);
    this.logger.info(`Created new profile '${profileId}' for environment '${environment}'`);

    return profile;
  }

  /**
   * Clone an existing profile
   */
  async cloneProfile(sourceProfileId, newProfileId) {
    const sourceProfile = await this.loadProfile(sourceProfileId);
    if (!sourceProfile) {
      throw new CLIError(
        `Source profile '${sourceProfileId}' not found`,
        "SOURCE_PROFILE_NOT_FOUND"
      );
    }

    const clonedProfile = {
      ...sourceProfile,
      id: newProfileId,
      name: `${sourceProfile.name} (Copy)`,
      description: `Cloned from ${sourceProfile.name}`,
      lastUpdated: new Date()
    };

    await this.saveProfile(clonedProfile);
    this.logger.info(`Cloned profile '${sourceProfileId}' to '${newProfileId}'`);

    return clonedProfile;
  }

  /**
   * Ensure profiles directory exists
   */
  async ensureProfilesDirectory() {
    try {
      await fs.mkdir(this.profilesDir, { recursive: true });
    } catch (error) {
      if (error.code !== "EEXIST") {
        throw error;
      }
    }
  }

  /**
   * Get profile file path
   */
  getProfilePath(profileId) {
    return path.join(this.profilesDir, `${profileId}.json`);
  }

  /**
   * Generate a profile ID from name
   */
  generateProfileId(name) {
    return name.toLowerCase().replace(/[^a-z0-9]/g, "-").replace(/-+/g, "-").replace(/^-|-$/g, "");
  }

  /**
   * Create a default profile
   */
  async createDefaultProfile() {
    const defaultProfile = await this.createProfile("default", "development");

    // Set as default
    await this.setDefaultProfile(defaultProfile.id);

    return defaultProfile;
  }

  /**
   * Initialize the service
   */
  async onInit() {
    await this.ensureProfilesDirectory();
    this.logger.info("ConfigurationService initialized", {
      configDirectory: this.config.configDirectory,
      profilesDirectory: this.profilesDir
    });
  }
}



