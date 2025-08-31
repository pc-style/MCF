/**
 * Claude Code execution options
 */
export interface ClaudeRunOptions {
  /** Working directory for Claude execution */
  workingDirectory?: string;

  /** Enable dangerous skip permissions flag */
  dangerousSkip?: boolean;

  /** Enable debug mode */
  debug?: boolean;

  /** Project name for Claude */
  project?: string;

  /** Custom configuration directory */
  configDirectory?: string;

  /** Anthropic base URL override */
  anthropicBaseUrl?: string;

  /** Anthropic auth token */
  authToken?: string;

  /** Anthropic model override */
  model?: string;

  /** Small fast model override */
  smallFastModel?: string;

  /** Additional arguments to pass through */
  additionalArgs?: string[];

  /** Interactive mode */
  interactive?: boolean;

  /** Environment variables to set */
  environment?: Record<string, string>;

  /** Timeout for execution (in milliseconds) */
  timeout?: number;
}

/**
 * Claude Code execution result
 */
export interface ClaudeRunResult {
  /** Exit code from Claude process */
  exitCode: number;

  /** Execution time in milliseconds */
  executionTime: number;

  /** Signal that terminated the process (if any) */
  signal?: string;

  /** Whether execution was successful */
  success: boolean;

  /** Output from Claude (if captured) */
  output?: string;

  /** Error output from Claude (if captured) */
  errorOutput?: string;
}

/**
 * Claude Code version information
 */
export interface ClaudeVersionInfo {
  /** Claude CLI version */
  version: string;

  /** Whether Claude is installed and accessible */
  installed: boolean;

  /** Path to Claude executable */
  executablePath?: string;

  /** Supported models */
  supportedModels?: string[];
}

/**
 * Claude environment configuration
 */
export interface ClaudeEnvironment {
  /** Anthropic base URL */
  ANTHROPIC_BASE_URL?: string;

  /** Anthropic auth token */
  ANTHROPIC_AUTH_TOKEN?: string;

  /** Claude config directory */
  CLAUDE_CONFIG_DIR?: string;

  /** Anthropic model */
  ANTHROPIC_MODEL?: string;

  /** Small fast model */
  ANTHROPIC_SMALL_FAST_MODEL?: string;
}

/**
 * Claude service interface for MCF CLI
 * Handles direct integration with Claude Code CLI
 */
export interface IClaudeService {
  /**
   * Execute Claude Code with specified options
   * @param options Execution options
   */
  runClaude(options?: ClaudeRunOptions): Promise<ClaudeRunResult>;

  /**
   * Get Claude Code version information
   */
  getVersion(): Promise<ClaudeVersionInfo>;

  /**
   * Check if Claude Code is installed and accessible
   */
  isInstalled(): Promise<boolean>;

  /**
   * Get the path to Claude executable
   */
  getExecutablePath(): Promise<string | null>;

  /**
   * Validate Claude configuration
   * @param config Configuration to validate
   */
  validateConfiguration(config: ClaudeRunOptions): Promise<ClaudeValidationResult>;

  /**
   * Build Claude command arguments from options
   * @param options Execution options
   */
  buildClaudeArguments(options: ClaudeRunOptions): string[];

  /**
   * Configure Claude environment variables
   * @param options Execution options
   */
  configureEnvironment(options: ClaudeRunOptions): Promise<ClaudeEnvironment>;

  /**
   * Get default Claude execution options
   */
  getDefaultOptions(): ClaudeRunOptions;

  /**
   * Kill a running Claude process
   * @param signal Signal to send (default: SIGTERM)
   */
  killClaude(signal?: string): Promise<boolean>;

  /**
   * Check if Claude process is currently running
   */
  isRunning(): Promise<boolean>;

  /**
   * Get Claude process information
   */
  getProcessInfo(): Promise<ClaudeProcessInfo | null>;
}

/**
 * Claude configuration validation result
 */
export interface ClaudeValidationResult {
  isValid: boolean;
  errors: string[];
  warnings: string[];
  suggestions: string[];
}

/**
 * Claude process information
 */
export interface ClaudeProcessInfo {
  /** Process ID */
  pid: number;

  /** Whether process is running */
  running: boolean;

  /** Start time */
  startTime: Date;

  /** Command line arguments */
  args: string[];

  /** Working directory */
  cwd: string;
}

/**
 * Claude service configuration
 */
export interface ClaudeServiceConfig {
  /** Default Claude executable path */
  defaultExecutablePath?: string;

  /** Default timeout for Claude execution */
  defaultTimeout?: number;

  /** Whether to validate Claude installation on startup */
  validateOnStartup?: boolean;

  /** Supported Claude models */
  supportedModels?: string[];

  /** Default environment configuration */
  defaultEnvironment?: Partial<ClaudeEnvironment>;
}


