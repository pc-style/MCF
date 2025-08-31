// IClaudeService interface for JavaScript runtime
export class IClaudeService {
  async runClaude(options) {
    throw new Error("Method 'runClaude' must be implemented");
  }

  async getVersion() {
    throw new Error("Method 'getVersion' must be implemented");
  }

  async isInstalled() {
    throw new Error("Method 'isInstalled' must be implemented");
  }

  async getExecutablePath() {
    throw new Error("Method 'getExecutablePath' must be implemented");
  }

  async validateConfiguration(config) {
    throw new Error("Method 'validateConfiguration' must be implemented");
  }

  buildClaudeArguments(options) {
    throw new Error("Method 'buildClaudeArguments' must be implemented");
  }

  async configureEnvironment(options) {
    throw new Error("Method 'configureEnvironment' must be implemented");
  }

  getDefaultOptions() {
    throw new Error("Method 'getDefaultOptions' must be implemented");
  }

  async killClaude(signal) {
    throw new Error("Method 'killClaude' must be implemented");
  }

  async isRunning() {
    throw new Error("Method 'isRunning' must be implemented");
  }

  async getProcessInfo() {
    throw new Error("Method 'getProcessInfo' must be implemented");
  }
}


