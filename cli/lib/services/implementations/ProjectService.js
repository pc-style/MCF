import { BaseService } from "../../core/base/BaseService.js";
import { ILogger } from "../../core/logging/ILogger.js";
import { LoggerFactory } from "../../core/logging/LoggerFactory.js";
import { CLIError } from "../../types/CLITypes.js";
import fs from "fs/promises";
import path from "path";

/**
 * Project service implementation for MCF CLI
 * Handles project lifecycle management and workspace operations
 */
export class ProjectService extends BaseService {
  constructor(config, logger, fileSystemService) {
    super();
    this.config = config || {};
    this.logger = logger || LoggerFactory.getLogger("ProjectService");
    this.fileSystemService = fileSystemService;

    this.defaultProjectPath = config?.defaultProjectPath || process.cwd();
    this.workspacePath = config?.workspacePath || path.join(process.cwd(), ".mcf", "workspaces");
    this.maxDiscoveryDepth = config?.maxDiscoveryDepth || 3;
    this.projectFileName = config?.projectFileName || ".mcf-project.json";
    this.workspaceFileName = config?.workspaceFileName || ".mcf-workspace.json";
    this.autoDiscover = config?.autoDiscover !== false;
  }

  /**
   * Create a new project
   */
  async createProject(options, customPath) {
    try {
      // Generate project ID and path
      const projectId = this.generateProjectId(options.name);
      const projectPath = customPath || path.join(this.defaultProjectPath, options.name);

      // Check if project already exists
      if (await this.projectExists(projectId)) {
        throw new CLIError(
          `Project '${projectId}' already exists`,
          "PROJECT_ALREADY_EXISTS"
        );
      }

      // Ensure project directory exists
      await this.fileSystemService.ensureDirectory(projectPath);

      // Create project information
      const project = {
        id: projectId,
        name: options.name,
        description: options.description || `Project ${options.name}`,
        path: projectPath,
        environment: options.environment || "development",
        createdAt: new Date(),
        lastModified: new Date(),
        config: {
          profile: options.profile,
          workspace: options.workspace,
          settings: options.settings || {}
        },
        metadata: {
          version: "1.0.0",
          author: process.env.USER || "unknown",
          tags: [],
          custom: {}
        }
      };

      // Validate project structure
      const validation = await this.validateProject(project);
      if (!validation.isValid) {
        throw new CLIError(
          `Project validation failed: ${validation.errors.join(", ")}`,
          "PROJECT_VALIDATION_FAILED",
          { errors: validation.errors }
        );
      }

      // Save project file
      const projectFilePath = path.join(projectPath, this.projectFileName);
      await this.fileSystemService.writeJSON(projectFilePath, project);

      this.logger.info(`Project '${projectId}' created successfully at ${projectPath}`);
      return project;
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to create project '${options.name}': ${message}`);
      throw new CLIError(
        `Failed to create project: ${message}`,
        "PROJECT_CREATE_FAILED",
        { projectName: options.name }
      );
    }
  }

  /**
   * Get project information by ID or path
   */
  async getProject(identifier) {
    try {
      // Try to find by ID first
      const projectPath = await this.findProjectPath(identifier);
      if (!projectPath) {
        return null;
      }

      const projectFilePath = path.join(projectPath, this.projectFileName);

      // Check if project file exists
      if (!(await this.fileSystemService.exists(projectFilePath))) {
        return null;
      }

      // Read and parse project file
      const project = await this.fileSystemService.readJSON(projectFilePath);

      // Validate project structure
      project.path = projectPath; // Ensure path is set correctly
      return project;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get project '${identifier}': ${message}`);
      throw new CLIError(
        `Failed to get project: ${message}`,
        "PROJECT_GET_FAILED",
        { identifier }
      );
    }
  }

  /**
   * List all projects
   */
  async listProjects(options = {}) {
    try {
      const projects = [];
      const searchPath = options.path || this.defaultProjectPath;
      const maxDepth = options.maxDepth || this.maxDiscoveryDepth;

      // Discover projects recursively
      const discoveredProjects = await this.discoverProjects(searchPath, {
        ...options,
        maxDepth
      });

      return discoveredProjects;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to list projects: ${message}`);
      throw new CLIError(
        `Failed to list projects: ${message}`,
        "PROJECT_LIST_FAILED"
      );
    }
  }

  /**
   * Update project information
   */
  async updateProject(projectId, updates) {
    try {
      const project = await this.getProject(projectId);
      if (!project) {
        throw new CLIError(
          `Project '${projectId}' not found`,
          "PROJECT_NOT_FOUND"
        );
      }

      // Apply updates
      const updatedProject = {
        ...project,
        ...updates,
        lastModified: new Date()
      };

      // Validate updated project
      const validation = await this.validateProject(updatedProject);
      if (!validation.isValid) {
        throw new CLIError(
          `Project validation failed: ${validation.errors.join(", ")}`,
          "PROJECT_VALIDATION_FAILED",
          { errors: validation.errors }
        );
      }

      // Save updated project
      const projectFilePath = path.join(project.path, this.projectFileName);
      await this.fileSystemService.writeJSON(projectFilePath, updatedProject);

      this.logger.info(`Project '${projectId}' updated successfully`);
      return updatedProject;
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to update project '${projectId}': ${message}`);
      throw new CLIError(
        `Failed to update project: ${message}`,
        "PROJECT_UPDATE_FAILED",
        { projectId }
      );
    }
  }

  /**
   * Delete a project
   */
  async deleteProject(projectId, deleteFiles = false) {
    try {
      const project = await this.getProject(projectId);
      if (!project) {
        this.logger.warn(`Project '${projectId}' not found, nothing to delete`);
        return false;
      }

      if (deleteFiles) {
        // Delete entire project directory
        await this.fileSystemService.remove(project.path);
        this.logger.info(`Project '${projectId}' and its files deleted`);
      } else {
        // Just remove project file
        const projectFilePath = path.join(project.path, this.projectFileName);
        await this.fileSystemService.remove(projectFilePath);
        this.logger.info(`Project '${projectId}' metadata deleted (files preserved)`);
      }

      return true;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to delete project '${projectId}': ${message}`);
      throw new CLIError(
        `Failed to delete project: ${message}`,
        "PROJECT_DELETE_FAILED",
        { projectId }
      );
    }
  }

  /**
   * Check if a project exists
   */
  async projectExists(identifier) {
    try {
      const project = await this.getProject(identifier);
      return project !== null;
    } catch {
      return false;
    }
  }

  /**
   * Get the current working project
   */
  async getCurrentProject() {
    try {
      // Check if current directory has a project
      const currentDir = process.cwd();
      const projectFilePath = path.join(currentDir, this.projectFileName);

      if (await this.fileSystemService.exists(projectFilePath)) {
        const project = await this.fileSystemService.readJSON(projectFilePath);
        project.path = currentDir;
        return project;
      }

      return null;
    } catch (error) {
      this.logger.debug(`No current project found: ${error.message}`);
      return null;
    }
  }

  /**
   * Set the current working project
   */
  async setCurrentProject(projectId) {
    try {
      const project = await this.getProject(projectId);
      if (!project) {
        throw new CLIError(
          `Project '${projectId}' not found`,
          "PROJECT_NOT_FOUND"
        );
      }

      // Change to project directory
      await this.fileSystemService.changeDirectory(project.path);
      this.logger.info(`Switched to project '${projectId}' at ${project.path}`);
    } catch (error) {
      if (error instanceof CLIError) {
        throw error;
      }

      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to set current project '${projectId}': ${message}`);
      throw new CLIError(
        `Failed to set current project: ${message}`,
        "SET_CURRENT_PROJECT_FAILED",
        { projectId }
      );
    }
  }

  /**
   * Discover projects in a directory
   */
  async discoverProjects(searchPath, options = {}) {
    try {
      const projects = [];
      const maxDepth = options.maxDepth || this.maxDiscoveryDepth;
      const visited = new Set();

      const discoverRecursive = async (currentPath, depth = 0) => {
        if (depth > maxDepth || visited.has(currentPath)) {
          return;
        }

        visited.add(currentPath);

        try {
          // Check if current directory is a project
          const projectFilePath = path.join(currentPath, this.projectFileName);
          if (await this.fileSystemService.exists(projectFilePath)) {
            const project = await this.fileSystemService.readJSON(projectFilePath);
            project.path = currentPath;

            // Apply filters if specified
            if (this.matchesFilter(project, options.filter)) {
              projects.push(project);
            }
          }

          // Continue discovery if recursive and not at max depth
          if (options.recursive !== false && depth < maxDepth) {
            const entries = await this.fileSystemService.listDirectory(currentPath);

            for (const entry of entries) {
              const entryPath = path.join(currentPath, entry);

              // Skip hidden directories unless explicitly included
              if (!options.includeHidden && entry.startsWith('.') && entry !== '.') {
                continue;
              }

              // Only recurse into directories
              if (await this.fileSystemService.isDirectory(entryPath)) {
                await discoverRecursive(entryPath, depth + 1);
              }
            }
          }
        } catch (error) {
          // Skip directories we can't read
          this.logger.debug(`Skipping ${currentPath}: ${error.message}`);
        }
      };

      await discoverRecursive(searchPath);
      return projects;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to discover projects in '${searchPath}': ${message}`);
      throw new CLIError(
        `Failed to discover projects: ${message}`,
        "PROJECT_DISCOVERY_FAILED",
        { searchPath }
      );
    }
  }

  /**
   * Initialize a project from a template
   */
  async initializeFromTemplate(template, options) {
    try {
      // For now, create a basic project
      // In a full implementation, this would copy template files
      this.logger.info(`Initializing project from template '${template}'`);
      return await this.createProject(options);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to initialize from template '${template}': ${message}`);
      throw new CLIError(
        `Failed to initialize from template: ${message}`,
        "TEMPLATE_INIT_FAILED",
        { template }
      );
    }
  }

  /**
   * Create a new workspace
   */
  async createWorkspace(name, workspacePath) {
    try {
      const workspaceId = this.generateWorkspaceId(name);

      // Ensure workspace directory exists
      await this.fileSystemService.ensureDirectory(workspacePath);

      const workspace = {
        id: workspaceId,
        name,
        path: workspacePath,
        projects: [],
        settings: {}
      };

      // Save workspace file
      const workspaceFilePath = path.join(workspacePath, this.workspaceFileName);
      await this.fileSystemService.writeJSON(workspaceFilePath, workspace);

      this.logger.info(`Workspace '${workspaceId}' created at ${workspacePath}`);
      return workspace;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to create workspace '${name}': ${message}`);
      throw new CLIError(
        `Failed to create workspace: ${message}`,
        "WORKSPACE_CREATE_FAILED",
        { workspaceName: name }
      );
    }
  }

  /**
   * Get workspace information
   */
  async getWorkspace(identifier) {
    try {
      const workspacePath = await this.findWorkspacePath(identifier);
      if (!workspacePath) {
        return null;
      }

      const workspaceFilePath = path.join(workspacePath, this.workspaceFileName);

      if (!(await this.fileSystemService.exists(workspaceFilePath))) {
        return null;
      }

      const workspace = await this.fileSystemService.readJSON(workspaceFilePath);
      workspace.path = workspacePath;

      return workspace;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get workspace '${identifier}': ${message}`);
      throw new CLIError(
        `Failed to get workspace: ${message}`,
        "WORKSPACE_GET_FAILED",
        { identifier }
      );
    }
  }

  /**
   * List all workspaces
   */
  async listWorkspaces() {
    try {
      // For now, return empty array as workspace management is basic
      // In a full implementation, this would discover workspaces
      return [];
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to list workspaces: ${message}`);
      throw new CLIError(
        `Failed to list workspaces: ${message}`,
        "WORKSPACE_LIST_FAILED"
      );
    }
  }

  /**
   * Add a project to a workspace
   */
  async addProjectToWorkspace(workspaceId, projectId) {
    // Basic implementation - in a full version this would update workspace metadata
    this.logger.info(`Added project '${projectId}' to workspace '${workspaceId}'`);
  }

  /**
   * Remove a project from a workspace
   */
  async removeProjectFromWorkspace(workspaceId, projectId) {
    // Basic implementation - in a full version this would update workspace metadata
    this.logger.info(`Removed project '${projectId}' from workspace '${workspaceId}'`);
  }

  /**
   * Validate project structure
   */
  async validateProject(project) {
    const errors = [];
    const warnings = [];
    const suggestions = [];

    // Validate required fields
    if (!project.id || typeof project.id !== "string" || project.id.trim() === "") {
      errors.push("Project ID is required and must be a non-empty string");
    }

    if (!project.name || typeof project.name !== "string" || project.name.trim() === "") {
      errors.push("Project name is required and must be a non-empty string");
    }

    if (!project.path || typeof project.path !== "string" || project.path.trim() === "") {
      errors.push("Project path is required and must be a non-empty string");
    }

    if (!project.environment || !["development", "production", "staging", "test"].includes(project.environment)) {
      errors.push("Project environment must be one of: development, production, staging, test");
    }

    if (!project.createdAt) {
      errors.push("Project creation date is required");
    }

    // Validate directory exists
    if (project.path && !(await this.fileSystemService.exists(project.path))) {
      warnings.push(`Project directory '${project.path}' does not exist`);
    }

    return {
      isValid: errors.length === 0,
      errors,
      warnings,
      suggestions
    };
  }

  /**
   * Get project statistics
   */
  async getProjectStats() {
    try {
      const projects = await this.listProjects();
      const workspaces = await this.listWorkspaces();

      const projectsByEnvironment = {
        development: 0,
        production: 0,
        staging: 0,
        test: 0
      };

      projects.forEach(project => {
        if (projectsByEnvironment[project.environment] !== undefined) {
          projectsByEnvironment[project.environment]++;
        }
      });

      // Sort projects by last modified date
      const recentProjects = projects
        .sort((a, b) => new Date(b.lastModified) - new Date(a.lastModified))
        .slice(0, 5);

      // Calculate workspace stats
      const averageProjectsPerWorkspace = workspaces.length > 0
        ? projects.length / workspaces.length
        : 0;

      const largestWorkspace = workspaces.length > 0
        ? workspaces.reduce((largest, current) =>
            current.projects.length > largest.projects.length ? current : largest
          )
        : { name: "None", projectCount: 0 };

      return {
        totalProjects: projects.length,
        totalWorkspaces: workspaces.length,
        projectsByEnvironment,
        recentProjects,
        workspaceStats: {
          averageProjectsPerWorkspace,
          largestWorkspace: {
            name: largestWorkspace.name,
            projectCount: largestWorkspace.projects.length
          }
        }
      };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Unknown error";
      this.logger.error(`Failed to get project stats: ${message}`);
      throw new CLIError(
        `Failed to get project stats: ${message}`,
        "PROJECT_STATS_FAILED"
      );
    }
  }

  /**
   * Find project path by ID
   */
  async findProjectPath(identifier) {
    // If identifier is an absolute path, use it directly
    if (path.isAbsolute(identifier)) {
      return identifier;
    }

    // If identifier contains path separators, treat as relative path
    if (identifier.includes('/') || identifier.includes('\\')) {
      return path.resolve(identifier);
    }

    // Otherwise, search for project by ID
    const searchPaths = [
      this.defaultProjectPath,
      process.cwd()
    ];

    for (const searchPath of searchPaths) {
      try {
        const projects = await this.discoverProjects(searchPath, {
          recursive: true,
          maxDepth: this.maxDiscoveryDepth
        });

        const project = projects.find(p => p.id === identifier);
        if (project) {
          return project.path;
        }
      } catch {
        // Continue searching
      }
    }

    return null;
  }

  /**
   * Find workspace path by ID
   */
  async findWorkspacePath(identifier) {
    // Basic implementation - in a full version this would search for workspaces
    if (path.isAbsolute(identifier)) {
      return identifier;
    }

    return path.resolve(this.workspacePath, identifier);
  }

  /**
   * Check if project matches filter criteria
   */
  matchesFilter(project, filter) {
    if (!filter) {
      return true;
    }

    if (filter.environment && project.environment !== filter.environment) {
      return false;
    }

    if (filter.hasConfig !== undefined) {
      const hasConfig = project.config && Object.keys(project.config).length > 0;
      if (hasConfig !== filter.hasConfig) {
        return false;
      }
    }

    if (filter.tags && filter.tags.length > 0) {
      const projectTags = project.metadata?.tags || [];
      if (!filter.tags.some(tag => projectTags.includes(tag))) {
        return false;
      }
    }

    return true;
  }

  /**
   * Generate a project ID from name
   */
  generateProjectId(name) {
    return name.toLowerCase().replace(/[^a-z0-9]/g, "-").replace(/-+/g, "-").replace(/^-|-$/g, "");
  }

  /**
   * Generate a workspace ID from name
   */
  generateWorkspaceId(name) {
    return name.toLowerCase().replace(/[^a-z0-9]/g, "-").replace(/-+/g, "-").replace(/^-|-$/g, "");
  }

  /**
   * Initialize the service
   */
  async onInit() {
    this.logger.info("ProjectService initialized", {
      defaultProjectPath: this.defaultProjectPath,
      workspacePath: this.workspacePath,
      autoDiscover: this.autoDiscover
    });
  }
}



