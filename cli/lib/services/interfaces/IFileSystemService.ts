import { Stats } from "fs";

/**
 * File system service interface for MCF CLI
 * Provides cross-platform file operations with permission management
 */
export interface IFileSystemService {
  /**
   * Read a file as text
   * @param filePath Path to the file
   * @param encoding Text encoding (default: 'utf-8')
   */
  readFile(filePath: string, encoding?: BufferEncoding): Promise<string>;

  /**
   * Read a file as JSON and parse it
   * @param filePath Path to the JSON file
   */
  readJSON<T = any>(filePath: string): Promise<T>;

  /**
   * Write text to a file
   * @param filePath Path to the file
   * @param content Content to write
   * @param encoding Text encoding (default: 'utf-8')
   */
  writeFile(filePath: string, content: string, encoding?: BufferEncoding): Promise<void>;

  /**
   * Write an object to a file as JSON
   * @param filePath Path to the file
   * @param data Object to serialize
   */
  writeJSON(filePath: string, data: any): Promise<void>;

  /**
   * Check if a file or directory exists
   * @param path Path to check
   */
  exists(path: string): Promise<boolean>;

  /**
   * Get file/directory statistics
   * @param path Path to check
   */
  stat(path: string): Promise<Stats>;

  /**
   * Create a directory (recursively)
   * @param dirPath Directory path to create
   */
  createDirectory(dirPath: string): Promise<void>;

  /**
   * Ensure a directory exists (alias for createDirectory)
   * @param dirPath Directory path to ensure
   */
  ensureDirectory(dirPath: string): Promise<void>;

  /**
   * Remove a file or directory
   * @param path Path to remove
   */
  remove(path: string): Promise<void>;

  /**
   * Copy a file or directory
   * @param source Source path
   * @param destination Destination path
   */
  copy(source: string, destination: string): Promise<void>;

  /**
   * Move/rename a file or directory
   * @param source Source path
   * @param destination Destination path
   */
  move(source: string, destination: string): Promise<void>;

  /**
   * List directory contents
   * @param dirPath Directory path
   */
  listDirectory(dirPath: string): Promise<string[]>;

  /**
   * Get the current working directory
   */
  getCurrentDirectory(): string;

  /**
   * Change the current working directory
   * @param dirPath New working directory
   */
  changeDirectory(dirPath: string): void;

  /**
   * Resolve a path to absolute
   * @param path Path to resolve
   */
  resolvePath(path: string): string;

  /**
   * Get the directory name of a path
   * @param path Path to get directory from
   */
  getDirectoryName(path: string): string;

  /**
   * Get the base name of a path
   * @param path Path to get base name from
   * @param ext Extension to remove (optional)
   */
  getBaseName(path: string, ext?: string): string;

  /**
   * Join path segments
   * @param paths Path segments to join
   */
  joinPath(...paths: string[]): string;

  /**
   * Check if path is absolute
   * @param path Path to check
   */
  isAbsolutePath(path: string): boolean;

  /**
   * Get file extension
   * @param path Path to get extension from
   */
  getExtension(path: string): string;

  /**
   * Check if path is a directory
   * @param path Path to check
   */
  isDirectory(path: string): Promise<boolean>;

  /**
   * Check if path is a file
   * @param path Path to check
   */
  isFile(path: string): Promise<boolean>;

  /**
   * Get file size in bytes
   * @param path Path to check
   */
  getFileSize(path: string): Promise<number>;

  /**
   * Get file modification time
   * @param path Path to check
   */
  getModificationTime(path: string): Promise<Date>;

  /**
   * Set file permissions (Unix-like)
   * @param path Path to modify
   * @param mode Permission mode (e.g., 0o755)
   */
  setPermissions(path: string, mode: number): Promise<void>;

  /**
   * Get file permissions (Unix-like)
   * @param path Path to check
   */
  getPermissions(path: string): Promise<number>;
}

/**
 * File system service configuration
 */
export interface FileSystemServiceConfig {
  baseDirectory?: string;
  createParentDirectories?: boolean;
  defaultEncoding?: BufferEncoding;
  maxFileSize?: number; // in bytes
  permissions?: {
    defaultFileMode?: number;
    defaultDirectoryMode?: number;
  };
}



