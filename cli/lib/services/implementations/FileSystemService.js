import { BaseService } from "../../core/base/BaseService.js";
import { ILogger } from "../../core/logging/ILogger.js";
import { LoggerFactory } from "../../core/logging/LoggerFactory.js";
import { CLIError } from "../../types/CLITypes.js";
import fs from "fs/promises";
import path from "path";
import { constants } from "fs";

/**
 * File system service implementation for MCF CLI
 * Provides cross-platform file operations with permission management
 */
export class FileSystemService extends BaseService {
  constructor(config, logger) {
    super();
    this.config = config || {};
    this.logger = logger || LoggerFactory.getLogger("FileSystemService");
    this.baseDirectory = config?.baseDirectory || process.cwd();
    this.defaultEncoding = config?.defaultEncoding || "utf-8";
    this.createParentDirectories = config?.createParentDirectories !== false;
    this.maxFileSize = config?.maxFileSize || 10 * 1024 * 1024; // 10MB default
  }

  /**
   * Read a file as text
   */
  async readFile(filePath, encoding) {
    try {
      const fullPath = this.resolvePath(filePath);
      const enc = encoding || this.defaultEncoding;

      this.logger.debug(`Reading file: ${fullPath}`);
      const content = await fs.readFile(fullPath, enc);

      // Check file size limit
      if (content.length > this.maxFileSize) {
        throw new CLIError(
          `File size exceeds maximum limit of ${this.maxFileSize} bytes`,
          "FILE_TOO_LARGE"
        );
      }

      return content;
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to read file '${filePath}': ${message}`);
      throw new CLIError(
        `Failed to read file: ${message}`,
        "FILE_READ_FAILED",
        { filePath }
      );
    }
  }

  /**
   * Read a file as JSON and parse it
   */
  async readJSON(filePath) {
    try {
      const content = await this.readFile(filePath, "utf-8");
      return JSON.parse(content);
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to read JSON file '${filePath}': ${message}`);
      throw new CLIError(
        `Failed to read JSON file: ${message}`,
        "JSON_READ_FAILED",
        { filePath }
      );
    }
  }

  /**
   * Write text to a file
   */
  async writeFile(filePath, content, encoding) {
    try {
      const fullPath = this.resolvePath(filePath);
      const enc = encoding || this.defaultEncoding;

      // Ensure parent directory exists if configured
      if (this.createParentDirectories) {
        const dirPath = this.getDirectoryName(fullPath);
        await this.ensureDirectory(dirPath);
      }

      this.logger.debug(`Writing file: ${fullPath}`);
      await fs.writeFile(fullPath, content, enc);

      // Set default permissions if configured
      if (this.config?.permissions?.defaultFileMode) {
        await this.setPermissions(fullPath, this.config.permissions.defaultFileMode);
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to write file '${filePath}': ${message}`);
      throw new CLIError(
        `Failed to write file: ${message}`,
        "FILE_WRITE_FAILED",
        { filePath }
      );
    }
  }

  /**
   * Write an object to a file as JSON
   */
  async writeJSON(filePath, data) {
    try {
      const jsonContent = JSON.stringify(data, null, 2);
      await this.writeFile(filePath, jsonContent, "utf-8");
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to write JSON file '${filePath}': ${message}`);
      throw new CLIError(
        `Failed to write JSON file: ${message}`,
        "JSON_WRITE_FAILED",
        { filePath }
      );
    }
  }

  /**
   * Check if a file or directory exists
   */
  async exists(filePath) {
    try {
      const fullPath = this.resolvePath(filePath);
      await fs.access(fullPath, constants.F_OK);
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Get file/directory statistics
   */
  async stat(filePath) {
    try {
      const fullPath = this.resolvePath(filePath);
      return await fs.stat(fullPath);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get stats for '${filePath}': ${message}`);
      throw new CLIError(
        `Failed to get file stats: ${message}`,
        "STAT_FAILED",
        { path: filePath }
      );
    }
  }

  /**
   * Create a directory (recursively)
   */
  async createDirectory(dirPath) {
    try {
      const fullPath = this.resolvePath(dirPath);
      this.logger.debug(`Creating directory: ${fullPath}`);
      await fs.mkdir(fullPath, { recursive: true });

      // Set default permissions if configured
      if (this.config?.permissions?.defaultDirectoryMode) {
        await this.setPermissions(fullPath, this.config.permissions.defaultDirectoryMode);
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to create directory '${dirPath}': ${message}`);
      throw new CLIError(
        `Failed to create directory: ${message}`,
        "DIRECTORY_CREATE_FAILED",
        { dirPath }
      );
    }
  }

  /**
   * Ensure a directory exists (alias for createDirectory)
   */
  async ensureDirectory(dirPath) {
    return this.createDirectory(dirPath);
  }

  /**
   * Remove a file or directory
   */
  async remove(filePath) {
    try {
      const fullPath = this.resolvePath(filePath);
      const stats = await this.stat(fullPath);

      if (stats.isDirectory()) {
        this.logger.debug(`Removing directory: ${fullPath}`);
        await fs.rm(fullPath, { recursive: true, force: true });
      } else {
        this.logger.debug(`Removing file: ${fullPath}`);
        await fs.unlink(fullPath);
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to remove '${filePath}': ${message}`);
      throw new CLIError(
        `Failed to remove: ${message}`,
        "REMOVE_FAILED",
        { path: filePath }
      );
    }
  }

  /**
   * Copy a file or directory
   */
  async copy(source, destination) {
    try {
      const fullSource = this.resolvePath(source);
      const fullDest = this.resolvePath(destination);

      // Ensure destination directory exists
      if (this.createParentDirectories) {
        const destDir = this.getDirectoryName(fullDest);
        await this.ensureDirectory(destDir);
      }

      this.logger.debug(`Copying from ${fullSource} to ${fullDest}`);
      await fs.cp(fullSource, fullDest, { recursive: true });
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to copy '${source}' to '${destination}': ${message}`);
      throw new CLIError(
        `Failed to copy: ${message}`,
        "COPY_FAILED",
        { source, destination }
      );
    }
  }

  /**
   * Move/rename a file or directory
   */
  async move(source, destination) {
    try {
      const fullSource = this.resolvePath(source);
      const fullDest = this.resolvePath(destination);

      // Ensure destination directory exists
      if (this.createParentDirectories) {
        const destDir = this.getDirectoryName(fullDest);
        await this.ensureDirectory(destDir);
      }

      this.logger.debug(`Moving from ${fullSource} to ${fullDest}`);
      await fs.rename(fullSource, fullDest);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to move '${source}' to '${destination}': ${message}`);
      throw new CLIError(
        `Failed to move: ${message}`,
        "MOVE_FAILED",
        { source, destination }
      );
    }
  }

  /**
   * List directory contents
   */
  async listDirectory(dirPath) {
    try {
      const fullPath = this.resolvePath(dirPath);
      this.logger.debug(`Listing directory: ${fullPath}`);
      return await fs.readdir(fullPath);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to list directory '${dirPath}': ${message}`);
      throw new CLIError(
        `Failed to list directory: ${message}`,
        "LIST_DIRECTORY_FAILED",
        { dirPath }
      );
    }
  }

  /**
   * Get the current working directory
   */
  getCurrentDirectory() {
    return process.cwd();
  }

  /**
   * Change the current working directory
   */
  changeDirectory(dirPath) {
    try {
      const fullPath = this.resolvePath(dirPath);
      process.chdir(fullPath);
      this.logger.debug(`Changed directory to: ${fullPath}`);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to change directory to '${dirPath}': ${message}`);
      throw new CLIError(
        `Failed to change directory: ${message}`,
        "CHANGE_DIRECTORY_FAILED",
        { dirPath }
      );
    }
  }

  /**
   * Resolve a path to absolute
   */
  resolvePath(filePath) {
    if (this.isAbsolutePath(filePath)) {
      return filePath;
    }
    return path.resolve(this.baseDirectory, filePath);
  }

  /**
   * Get the directory name of a path
   */
  getDirectoryName(filePath) {
    return path.dirname(filePath);
  }

  /**
   * Get the base name of a path
   */
  getBaseName(filePath, ext) {
    return path.basename(filePath, ext);
  }

  /**
   * Join path segments
   */
  joinPath(...paths) {
    return path.join(...paths);
  }

  /**
   * Check if path is absolute
   */
  isAbsolutePath(filePath) {
    return path.isAbsolute(filePath);
  }

  /**
   * Get file extension
   */
  getExtension(filePath) {
    return path.extname(filePath);
  }

  /**
   * Check if path is a directory
   */
  async isDirectory(filePath) {
    try {
      const stats = await this.stat(filePath);
      return stats.isDirectory();
    } catch {
      return false;
    }
  }

  /**
   * Check if path is a file
   */
  async isFile(filePath) {
    try {
      const stats = await this.stat(filePath);
      return stats.isFile();
    } catch {
      return false;
    }
  }

  /**
   * Get file size in bytes
   */
  async getFileSize(filePath) {
    try {
      const stats = await this.stat(filePath);
      return stats.size;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get file size for '${filePath}': ${message}`);
      throw new CLIError(
        `Failed to get file size: ${message}`,
        "GET_FILE_SIZE_FAILED",
        { path: filePath }
      );
    }
  }

  /**
   * Get file modification time
   */
  async getModificationTime(filePath) {
    try {
      const stats = await this.stat(filePath);
      return stats.mtime;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get modification time for '${filePath}': ${message}`);
      throw new CLIError(
        `Failed to get modification time: ${message}`,
        "GET_MODIFICATION_TIME_FAILED",
        { path: filePath }
      );
    }
  }

  /**
   * Set file permissions (Unix-like)
   */
  async setPermissions(filePath, mode) {
    try {
      const fullPath = this.resolvePath(filePath);
      await fs.chmod(fullPath, mode);
      this.logger.debug(`Set permissions ${mode.toString(8)} on ${fullPath}`);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to set permissions on '${filePath}': ${message}`);
      // Don't throw error for permission setting failures on some platforms
      this.logger.warn(`Could not set permissions: ${message}`);
    }
  }

  /**
   * Get file permissions (Unix-like)
   */
  async getPermissions(filePath) {
    try {
      const stats = await this.stat(filePath);
      return stats.mode;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get permissions for '${filePath}': ${message}`);
      throw new CLIError(
        `Failed to get permissions: ${message}`,
        "GET_PERMISSIONS_FAILED",
        { path: filePath }
      );
    }
  }

  /**
   * Initialize the service
   */
  async onInit() {
    this.logger.info("FileSystemService initialized", {
      baseDirectory: this.baseDirectory,
      createParentDirectories: this.createParentDirectories,
      maxFileSize: this.maxFileSize
    });
  }
}
