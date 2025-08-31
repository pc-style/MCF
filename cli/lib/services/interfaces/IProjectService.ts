import { CLIEnvironment } from "../../types/CLITypes.js";

/**
 * Project information structure
 */
export interface ProjectInfo {
  id: string;
  name: string;
  description?: string;
  path: string;
  environment: CLIEnvironment;
  createdAt: Date;
  lastModified: Date;
  config?: {
    profile?: string;
    workspace?: string;
    settings?: Record<string, any>;
  };
  metadata?: {
    version?: string;
    author?: string;
    tags?: string[];
    custom?: Record<string, any>;
  };
}

/**
 * Project creation options
 */
export interface ProjectCreateOptions {
  name: string;
  description?: string;
  environment?: CLIEnvironment;
  template?: string;
  profile?: string;
  workspace?: string;
  settings?: Record<string, any>;
}

/**
 * Project discovery options
 */
export interface ProjectDiscoveryOptions {
  path?: string;
  recursive?: boolean;
  maxDepth?: number;
  includeHidden?: boolean;
  filter?: {
    environment?: CLIEnvironment;
    hasConfig?: boolean;
    tags?: string[];
  };
}

/**
 * Project workspace information
 */
export interface ProjectWorkspace {
  id: string;
  name: string;
  path: string;
  projects: string[];
  defaultProject?: string;
  settings?: Record<string, any>;
}

/**
 * Project service interface for MCF CLI
 * Handles project lifecycle management and workspace operations
 */
export interface IProjectService {
  /**
   * Create a new project
   * @param options Project creation options
   * @param path Optional custom path for project creation
   */
  createProject(options: ProjectCreateOptions, path?: string): Promise<ProjectInfo>;

  /**
   * Get project information by ID or path
   * @param identifier Project ID or path
   */
  getProject(identifier: string): Promise<ProjectInfo | null>;

  /**
   * List all projects in a workspace or directory
   * @param options Discovery options
   */
  listProjects(options?: ProjectDiscoveryOptions): Promise<ProjectInfo[]>;

  /**
   * Update project information
   * @param projectId Project identifier
   * @param updates Partial project information to update
   */
  updateProject(projectId: string, updates: Partial<ProjectInfo>): Promise<ProjectInfo>;

  /**
   * Delete a project
   * @param projectId Project identifier
   * @param deleteFiles Whether to delete project files (default: false)
   */
  deleteProject(projectId: string, deleteFiles?: boolean): Promise<boolean>;

  /**
   * Check if a project exists
   * @param identifier Project ID or path
   */
  projectExists(identifier: string): Promise<boolean>;

  /**
   * Get the current working project
   */
  getCurrentProject(): Promise<ProjectInfo | null>;

  /**
   * Set the current working project
   * @param projectId Project identifier
   */
  setCurrentProject(projectId: string): Promise<void>;

  /**
   * Discover projects in a directory
   * @param path Directory to scan
   * @param options Discovery options
   */
  discoverProjects(path: string, options?: ProjectDiscoveryOptions): Promise<ProjectInfo[]>;

  /**
   * Initialize a project from a template
   * @param template Template name or path
   * @param options Project creation options
   */
  initializeFromTemplate(template: string, options: ProjectCreateOptions): Promise<ProjectInfo>;

  /**
   * Create a new workspace
   * @param name Workspace name
   * @param path Workspace path
   */
  createWorkspace(name: string, path: string): Promise<ProjectWorkspace>;

  /**
   * Get workspace information
   * @param identifier Workspace ID or path
   */
  getWorkspace(identifier: string): Promise<ProjectWorkspace | null>;

  /**
   * List all workspaces
   */
  listWorkspaces(): Promise<ProjectWorkspace[]>;

  /**
   * Add a project to a workspace
   * @param workspaceId Workspace identifier
   * @param projectId Project identifier
   */
  addProjectToWorkspace(workspaceId: string, projectId: string): Promise<void>;

  /**
   * Remove a project from a workspace
   * @param workspaceId Workspace identifier
   * @param projectId Project identifier
   */
  removeProjectFromWorkspace(workspaceId: string, projectId: string): Promise<void>;

  /**
   * Validate project structure
   * @param project Project information to validate
   */
  validateProject(project: ProjectInfo): Promise<ProjectValidationResult>;

  /**
   * Get project statistics
   */
  getProjectStats(): Promise<ProjectStats>;
}

/**
 * Project validation result
 */
export interface ProjectValidationResult {
  isValid: boolean;
  errors: string[];
  warnings: string[];
  suggestions: string[];
}

/**
 * Project statistics
 */
export interface ProjectStats {
  totalProjects: number;
  totalWorkspaces: number;
  projectsByEnvironment: Record<CLIEnvironment, number>;
  recentProjects: ProjectInfo[];
  workspaceStats: {
    averageProjectsPerWorkspace: number;
    largestWorkspace: {
      name: string;
      projectCount: number;
    };
  };
}

/**
 * Project service configuration
 */
export interface ProjectServiceConfig {
  defaultProjectPath?: string;
  workspacePath?: string;
  maxDiscoveryDepth?: number;
  projectFileName?: string;
  workspaceFileName?: string;
  autoDiscover?: boolean;
  validation?: {
    strict?: boolean;
    checkDependencies?: boolean;
  };
}



