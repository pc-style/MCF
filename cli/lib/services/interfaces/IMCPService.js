// IMCPService interface for JavaScript runtime
export class IMCPService {
  async registerServer(config) {
    throw new Error("Method 'registerServer' must be implemented");
  }

  async unregisterServer(serverId) {
    throw new Error("Method 'unregisterServer' must be implemented");
  }

  async getServer(serverId) {
    throw new Error("Method 'getServer' must be implemented");
  }

  async listServers() {
    throw new Error("Method 'listServers' must be implemented");
  }

  async startServer(serverId, options) {
    throw new Error("Method 'startServer' must be implemented");
  }

  async stopServer(serverId, force) {
    throw new Error("Method 'stopServer' must be implemented");
  }

  async restartServer(serverId) {
    throw new Error("Method 'restartServer' must be implemented");
  }

  async isServerRunning(serverId) {
    throw new Error("Method 'isServerRunning' must be implemented");
  }

  async getServerHealth(serverId) {
    throw new Error("Method 'getServerHealth' must be implemented");
  }

  async getAllServerHealth() {
    throw new Error("Method 'getAllServerHealth' must be implemented");
  }

  async startAutoStartServers() {
    throw new Error("Method 'startAutoStartServers' must be implemented");
  }

  async stopAllServers(force) {
    throw new Error("Method 'stopAllServers' must be implemented");
  }

  async installServer(packageName, options) {
    throw new Error("Method 'installServer' must be implemented");
  }

  async updateServer(serverId) {
    throw new Error("Method 'updateServer' must be implemented");
  }

  async removeServer(serverId, keepConfig) {
    throw new Error("Method 'removeServer' must be implemented");
  }

  async getServerLogs(serverId, lines) {
    throw new Error("Method 'getServerLogs' must be implemented");
  }

  async validateServerConfig(config) {
    throw new Error("Method 'validateServerConfig' must be implemented");
  }

  async getServiceStats() {
    throw new Error("Method 'getServiceStats' must be implemented");
  }

  async exportConfigurations() {
    throw new Error("Method 'exportConfigurations' must be implemented");
  }

  async importConfigurations(configData) {
    throw new Error("Method 'importConfigurations' must be implemented");
  }
}


