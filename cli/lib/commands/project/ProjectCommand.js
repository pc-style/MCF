import { BaseCommand } from "../../core/interfaces/BaseCommand.js";
import { ServiceRegistry } from "../../core/registry/ServiceRegistry.js";
import { ILogger } from "../../core/logging/ILogger.js";
import { LoggerFactory } from "../../core/logging/LoggerFactory.js";
import { CLIError } from "../../types/CLITypes.js";
import chalk from "chalk";

/**
 * ProjectCommand - MCF CLI Project Management
 * Provides subcommands for managing MCF projects and workspaces
 */
export class ProjectCommand extends BaseCommand {
  constructor(serviceRegistry) {
    super();
    this.serviceRegistry = serviceRegistry;
    this.logger = LoggerFactory.getLogger("ProjectCommand");
    this.projectService = null;
    this.fileSystemService = null;
  }

  static get metadata() {
    return {
      name: "ProjectCommand",
      description: "Manage MCF projects and workspaces",
      category: "project",
      version: "1.0.0",
      dependencies: {
        services: ["IProjectService", "IFileSystemService"],
        commands: [],
        external: []
      }
    };
  }

  async initialize() {
    try {
      this.projectService = this.serviceRegistry.getService("IProjectService");
      this.fileSystemService = this.serviceRegistry.getService("IFileSystemService");
      this.logger.debug("ProjectCommand initialized with services");
    } catch (error) {
      this.logger.error("Failed to initialize ProjectCommand", error);
      throw new CLIError(
        "Failed to initialize project services",
        "PROJECT_COMMAND_INIT_FAILED"
      );
    }
  }

  async execute(args = []) {
    await this.initialize();

    if (args.length === 0) {
      return this.showHelp();
    }

    const [subcommand, ...subArgs] = args;

    switch (subcommand.toLowerCase()) {
      case "list":
      case "ls":
        return await this.listProjects();
      case "show":
      case "info":
        return await this.showProject(subArgs[0]);
      case "create":
      case "new":
        return await this.createProject(subArgs);
      case "delete":
      case "del":
      case "remove":
      case "rm":
        return await this.deleteProject(subArgs);
      case "switch":
      case "cd":
        return await this.switchProject(subArgs[0]);
      case "current":
        return await this.showCurrentProject();
      case "discover":
        return await this.discoverProjects(subArgs);
      case "stats":
        return await this.showStats();
      case "workspace":
      case "ws":
        return await this.handleWorkspaceCommand(subArgs);
      default:
        console.log(chalk.red(`Unknown project subcommand: ${subcommand}`));
        return this.showHelp();
    }
  }

  async listProjects() {
    try {
      const projects = await this.projectService.listProjects();

      if (projects.length === 0) {
        console.log(chalk.yellow("No projects found."));
        console.log("Create one with: mcf project create <name> [description]");
        console.log("Or discover existing projects: mcf project discover");
        return;
      }

      console.log(chalk.blue("MCF Projects:"));
      console.log();

      const currentProject = await this.projectService.getCurrentProject();
      const currentProjectId = currentProject?.id;

      for (const project of projects) {
        const marker = project.id === currentProjectId ? chalk.green(" (current)") : "";
        const envColor = this.getEnvironmentColor(project.environment);
        console.log(`  ${chalk.cyan(project.id)}${marker}`);
        console.log(`    Name: ${project.name}`);
        console.log(`    Environment: ${envColor(project.environment)}`);
        console.log(`    Path: ${chalk.gray(project.path)}`);
        if (project.description) {
          console.log(`    Description: ${project.description}`);
        }
        console.log();
      }

      console.log(`Total: ${projects.length} project(s)`);
    } catch (error) {
      console.error(chalk.red(`Failed to list projects: ${error.message}`));
      throw error;
    }
  }

  async showProject(projectId) {
    if (!projectId) {
      console.log(chalk.red("Project ID is required"));
      console.log("Usage: mcf project show <project-id>");
      return;
    }

    try {
      const project = await this.projectService.getProject(projectId);

      if (!project) {
        console.log(chalk.red(`Project '${projectId}' not found`));
        return;
      }

      console.log(chalk.blue(`Project: ${project.id}`));
      console.log(chalk.gray("─".repeat(50)));
      console.log(`Name: ${chalk.cyan(project.name)}`);
      console.log(`Environment: ${this.getEnvironmentColor(project.environment)(project.environment)}`);
      console.log(`Path: ${chalk.gray(project.path)}`);

      if (project.description) {
        console.log(`Description: ${project.description}`);
      }

      console.log(`Created: ${new Date(project.createdAt).toLocaleString()}`);
      console.log(`Last Modified: ${new Date(project.lastModified).toLocaleString()}`);

      if (project.metadata) {
        console.log();
        console.log(chalk.blue("Metadata:"));
        if (project.metadata.version) {
          console.log(`  Version: ${project.metadata.version}`);
        }
        if (project.metadata.author) {
          console.log(`  Author: ${project.metadata.author}`);
        }
        if (project.metadata.tags && project.metadata.tags.length > 0) {
          console.log(`  Tags: ${project.metadata.tags.join(", ")}`);
        }
      }

      if (project.config) {
        console.log();
        console.log(chalk.blue("Configuration:"));
        if (project.config.profile) {
          console.log(`  Profile: ${project.config.profile}`);
        }
        if (project.config.workspace) {
          console.log(`  Workspace: ${project.config.workspace}`);
        }
        if (project.config.settings && Object.keys(project.config.settings).length > 0) {
          console.log(`  Settings: ${JSON.stringify(project.config.settings, null, 2)}`);
        }
      }
    } catch (error) {
      console.error(chalk.red(`Failed to show project: ${error.message}`));
      throw error;
    }
  }

  async createProject(args) {
    const [name, ...descriptionParts] = args;
    const description = descriptionParts.join(" ");

    if (!name) {
      console.log(chalk.red("Project name is required"));
      console.log("Usage: mcf project create <name> [description]");
      return;
    }

    try {
      const options = {
        name,
        description: description || undefined,
        environment: "development"
      };

      const project = await this.projectService.createProject(options);
      console.log(chalk.green(`Project '${project.id}' created successfully`));
      console.log(`Path: ${project.path}`);
      console.log(`Environment: ${project.environment}`);
      console.log();
      console.log("You can now:");
      console.log(`  • Switch to project: mcf project switch ${project.id}`);
      console.log(`  • Show details: mcf project show ${project.id}`);
      console.log(`  • Edit settings: mcf project edit ${project.id}`);
    } catch (error) {
      console.error(chalk.red(`Failed to create project: ${error.message}`));
      throw error;
    }
  }

  async deleteProject(args) {
    const [projectId, ...flags] = args;
    const deleteFiles = flags.includes("--files") || flags.includes("-f");

    if (!projectId) {
      console.log(chalk.red("Project ID is required"));
      console.log("Usage: mcf project delete <project-id> [--files]");
      console.log("  --files: Also delete project files (default: metadata only)");
      return;
    }

    try {
      // Confirm deletion if deleting files
      if (deleteFiles) {
        console.log(chalk.yellow(`⚠️  This will permanently delete the project directory and all files:`));
        console.log(chalk.red(projectId));
        console.log();
        console.log("This action cannot be undone!");
        console.log();
        console.log("Use Ctrl+C to cancel, or press Enter to continue...");
        process.stdin.setRawMode(true);
        process.stdin.resume();
        await new Promise(resolve => {
          process.stdin.once('data', () => resolve());
        });
        process.stdin.setRawMode(false);
        process.stdin.pause();
      }

      const success = await this.projectService.deleteProject(projectId, deleteFiles);

      if (success) {
        if (deleteFiles) {
          console.log(chalk.green(`Project '${projectId}' and all files deleted successfully`));
        } else {
          console.log(chalk.green(`Project '${projectId}' metadata deleted successfully`));
          console.log(chalk.gray("Project files preserved at original location"));
        }
      }
    } catch (error) {
      console.error(chalk.red(`Failed to delete project: ${error.message}`));
      throw error;
    }
  }

  async switchProject(projectId) {
    if (!projectId) {
      console.log(chalk.red("Project ID is required"));
      console.log("Usage: mcf project switch <project-id>");
      return;
    }

    try {
      await this.projectService.setCurrentProject(projectId);
      console.log(chalk.green(`Switched to project '${projectId}'`));
      console.log(`Working directory: ${process.cwd()}`);
    } catch (error) {
      console.error(chalk.red(`Failed to switch project: ${error.message}`));
      throw error;
    }
  }

  async showCurrentProject() {
    try {
      const currentProject = await this.projectService.getCurrentProject();

      if (!currentProject) {
        console.log(chalk.yellow("No current project"));
        console.log("You can:");
        console.log("  • Switch to a project: mcf project switch <project-id>");
        console.log("  • Create a new project: mcf project create <name>");
        console.log("  • List available projects: mcf project list");
        return;
      }

      console.log(chalk.blue("Current Project:"));
      console.log(`ID: ${chalk.cyan(currentProject.id)}`);
      console.log(`Name: ${currentProject.name}`);
      console.log(`Environment: ${this.getEnvironmentColor(currentProject.environment)(currentProject.environment)}`);
      console.log(`Path: ${chalk.gray(currentProject.path)}`);
    } catch (error) {
      console.error(chalk.red(`Failed to get current project: ${error.message}`));
      throw error;
    }
  }

  async discoverProjects(args) {
    const [searchPath = "."] = args;

    try {
      console.log(chalk.blue(`Discovering projects in: ${searchPath}`));
      const projects = await this.projectService.discoverProjects(searchPath);

      if (projects.length === 0) {
        console.log(chalk.yellow("No projects found in the specified path"));
        return;
      }

      console.log(chalk.green(`Found ${projects.length} project(s):`));
      console.log();

      for (const project of projects) {
        console.log(`  ${chalk.cyan(project.id)} - ${project.name}`);
        console.log(`    Path: ${chalk.gray(project.path)}`);
        console.log(`    Environment: ${this.getEnvironmentColor(project.environment)(project.environment)}`);
        console.log();
      }
    } catch (error) {
      console.error(chalk.red(`Failed to discover projects: ${error.message}`));
      throw error;
    }
  }

  async showStats() {
    try {
      const stats = await this.projectService.getProjectStats();

      console.log(chalk.blue("Project Statistics"));
      console.log(chalk.gray("─".repeat(50)));
      console.log(`Total Projects: ${chalk.cyan(stats.totalProjects)}`);
      console.log(`Total Workspaces: ${chalk.cyan(stats.totalWorkspaces)}`);
      console.log();

      console.log(chalk.blue("Projects by Environment:"));
      Object.entries(stats.projectsByEnvironment).forEach(([env, count]) => {
        const envColor = this.getEnvironmentColor(env);
        console.log(`  ${envColor(env)}: ${count}`);
      });

      if (stats.workspaceStats.largestWorkspace.projectCount > 0) {
        console.log();
        console.log(chalk.blue("Workspace Stats:"));
        console.log(`  Average projects per workspace: ${stats.workspaceStats.averageProjectsPerWorkspace.toFixed(1)}`);
        console.log(`  Largest workspace: ${stats.workspaceStats.largestWorkspace.name} (${stats.workspaceStats.largestWorkspace.projectCount} projects)`);
      }

      if (stats.recentProjects.length > 0) {
        console.log();
        console.log(chalk.blue("Recent Projects:"));
        stats.recentProjects.slice(0, 3).forEach(project => {
          const timeAgo = this.getTimeAgo(new Date(project.lastModified));
          console.log(`  ${chalk.cyan(project.id)} - ${timeAgo}`);
        });
      }
    } catch (error) {
      console.error(chalk.red(`Failed to get project stats: ${error.message}`));
      throw error;
    }
  }

  async handleWorkspaceCommand(args) {
    const [wsCommand, ...wsArgs] = args;

    switch (wsCommand?.toLowerCase()) {
      case "list":
      case "ls":
        return await this.listWorkspaces();
      case "create":
        return await this.createWorkspace(wsArgs);
      case "show":
        return await this.showWorkspace(wsArgs[0]);
      default:
        console.log(chalk.red("Workspace command required"));
        console.log("Available workspace commands:");
        console.log("  list, ls          List all workspaces");
        console.log("  create <name>     Create new workspace");
        console.log("  show <id>         Show workspace details");
        return;
    }
  }

  async listWorkspaces() {
    try {
      const workspaces = await this.projectService.listWorkspaces();
      console.log(chalk.yellow("Workspace management is under development"));
      console.log(`Found ${workspaces.length} workspace(s)`);
    } catch (error) {
      console.error(chalk.red(`Failed to list workspaces: ${error.message}`));
      throw error;
    }
  }

  async createWorkspace(args) {
    const [name, workspacePath] = args;

    if (!name) {
      console.log(chalk.red("Workspace name is required"));
      console.log("Usage: mcf project workspace create <name> [path]");
      return;
    }

    try {
      const path = workspacePath || `./${name}`;
      const workspace = await this.projectService.createWorkspace(name, path);
      console.log(chalk.green(`Workspace '${workspace.id}' created successfully`));
      console.log(`Path: ${workspace.path}`);
    } catch (error) {
      console.error(chalk.red(`Failed to create workspace: ${error.message}`));
      throw error;
    }
  }

  async showWorkspace(workspaceId) {
    if (!workspaceId) {
      console.log(chalk.red("Workspace ID is required"));
      console.log("Usage: mcf project workspace show <workspace-id>");
      return;
    }

    try {
      const workspace = await this.projectService.getWorkspace(workspaceId);

      if (!workspace) {
        console.log(chalk.red(`Workspace '${workspaceId}' not found`));
        return;
      }

      console.log(chalk.blue(`Workspace: ${workspace.id}`));
      console.log(`Name: ${workspace.name}`);
      console.log(`Path: ${chalk.gray(workspace.path)}`);
      console.log(`Projects: ${workspace.projects.length}`);
    } catch (error) {
      console.error(chalk.red(`Failed to show workspace: ${error.message}`));
      throw error;
    }
  }

  showHelp() {
    console.log(chalk.blue("MCF Project Command Help"));
    console.log(chalk.gray("─".repeat(50)));
    console.log();
    console.log("Manage MCF projects and workspaces");
    console.log();
    console.log(chalk.blue("Project Commands:"));
    console.log();
    console.log("  list, ls                    List all projects");
    console.log("  show, info <id>             Show project details");
    console.log("  create, new <name> [desc]   Create new project");
    console.log("  delete, del <id> [--files]  Delete project");
    console.log("  switch, cd <id>             Switch to project");
    console.log("  current                     Show current project");
    console.log("  discover [path]             Discover projects in path");
    console.log("  stats                       Show project statistics");
    console.log();
    console.log(chalk.blue("Workspace Commands:"));
    console.log();
    console.log("  workspace list, ws ls       List workspaces");
    console.log("  workspace create <name>     Create workspace");
    console.log("  workspace show <id>         Show workspace details");
    console.log();
    console.log(chalk.blue("Examples:"));
    console.log();
    console.log("  mcf project list");
    console.log("  mcf project create myapp 'My awesome app'");
    console.log("  mcf project switch myapp");
    console.log("  mcf project show myapp");
    console.log("  mcf project discover ./projects");
    console.log("  mcf project workspace create myworkspace");
    console.log();
    return { success: true };
  }

  getEnvironmentColor(environment) {
    switch (environment) {
      case "development":
        return chalk.blue;
      case "production":
        return chalk.red;
      case "staging":
        return chalk.yellow;
      case "test":
        return chalk.green;
      default:
        return chalk.gray;
    }
  }

  getTimeAgo(date) {
    const now = new Date();
    const diffMs = now - date;
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

    if (diffDays === 0) {
      return "today";
    } else if (diffDays === 1) {
      return "yesterday";
    } else if (diffDays < 7) {
      return `${diffDays} days ago`;
    } else if (diffDays < 30) {
      const weeks = Math.floor(diffDays / 7);
      return `${weeks} week${weeks > 1 ? 's' : ''} ago`;
    } else {
      const months = Math.floor(diffDays / 30);
      return `${months} month${months > 1 ? 's' : ''} ago`;
    }
  }

  getMetadata() {
    return ProjectCommand.metadata;
  }
}
