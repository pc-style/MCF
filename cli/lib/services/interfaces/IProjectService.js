// IProjectService interface for JavaScript runtime
export class IProjectService {
  async createProject(options, path) {
    throw new Error("Method 'createProject' must be implemented");
  }

  async getProject(identifier) {
    throw new Error("Method 'getProject' must be implemented");
  }

  async listProjects(options) {
    throw new Error("Method 'listProjects' must be implemented");
  }

  async updateProject(projectId, updates) {
    throw new Error("Method 'updateProject' must be implemented");
  }

  async deleteProject(projectId, deleteFiles) {
    throw new Error("Method 'deleteProject' must be implemented");
  }

  async projectExists(identifier) {
    throw new Error("Method 'projectExists' must be implemented");
  }

  async getCurrentProject() {
    throw new Error("Method 'getCurrentProject' must be implemented");
  }

  async setCurrentProject(projectId) {
    throw new Error("Method 'setCurrentProject' must be implemented");
  }

  async discoverProjects(path, options) {
    throw new Error("Method 'discoverProjects' must be implemented");
  }

  async initializeFromTemplate(template, options) {
    throw new Error("Method 'initializeFromTemplate' must be implemented");
  }

  async createWorkspace(name, path) {
    throw new Error("Method 'createWorkspace' must be implemented");
  }

  async getWorkspace(identifier) {
    throw new Error("Method 'getWorkspace' must be implemented");
  }

  async listWorkspaces() {
    throw new Error("Method 'listWorkspaces' must be implemented");
  }

  async addProjectToWorkspace(workspaceId, projectId) {
    throw new Error("Method 'addProjectToWorkspace' must be implemented");
  }

  async removeProjectFromWorkspace(workspaceId, projectId) {
    throw new Error("Method 'removeProjectFromWorkspace' must be implemented");
  }

  async validateProject(project) {
    throw new Error("Method 'validateProject' must be implemented");
  }

  async getProjectStats() {
    throw new Error("Method 'getProjectStats' must be implemented");
  }
}



