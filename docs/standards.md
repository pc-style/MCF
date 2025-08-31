# MCF CLI Development Standards

**Based on UCM Micro-Block Development Guidelines**

## Overview

This document establishes development standards for MCF CLI based on UCM micro-block architecture patterns. These standards ensure consistency, maintainability, and proper implementation of architectural patterns across the codebase.

## Critical Core Principles

### 🚨 ES MODULES ONLY 🚨

MCF CLI uses ES Modules exclusively. **NEVER use CommonJS**:

- ❌ **FORBIDDEN**: `require()`, `module.exports`, `exports`
- ✅ **REQUIRED**: `import`, `export`, `export default`

**Why ES Modules are mandatory:**

- Modern JavaScript standard (not legacy CommonJS)
- Better tree-shaking and dead code elimination
- Static analysis for bundlers and tools
- Top-level await support
- Improved TypeScript integration
- Future-proof architecture

### Contract-First Development

All CLI commands and services must define their contracts before implementation:

1. **Define Input/Output Interfaces** - Start with TypeScript interfaces
2. **Define Error Contracts** - Create specific error classes
3. **Define Metadata** - Specify dependencies and capabilities
4. **Implement Logic** - Write the actual implementation
5. **Test Contracts** - Verify interface compliance

## CLI Command Development Guidelines

### CLI Command Structure

Every CLI command must follow the micro-block pattern:

```typescript
// 1. Input Contract
interface RunCommandInput {
  dangerousSkip?: boolean;
  continue?: boolean;
  profile?: string;
  passThroughArgs?: string[];
}

// 2. Output Contract
interface RunCommandOutput {
  claudeExitCode: number;
  executionTime: number;
  configUsed: string;
}

// 3. Error Contract
class RunCommandError extends BaseError {
  constructor(message: string, code?: string, details?: Record<string, any>) {
    super(message, code, details);
    this.name = "RunCommandError";
  }
}

// 4. Command Implementation
export class RunCommand
  implements OutputCommand<RunCommandInput, RunCommandOutput, RunCommandError>
{
  static readonly metadata: CommandMetadata = {
    name: "RunCommand",
    description: "Executes Claude Code with configuration and flags",
    category: "run",
    inputType: "RunCommandInput",
    outputType: "RunCommandOutput",
    errorType: "RunCommandError",
    version: "1.0.0",
    contractVersion: "1.0",
    dependencies: {
      services: ["IClaudeService", "IConfigurationService"],
      commands: [],
      external: ["child_process"],
    },
  };

  constructor(
    public input?: RunCommandInput,
    private logger?: ICommandLogger,
    private services?: Record<string, any>,
  ) {
    // Initialize services and validate dependencies
  }

  validate(): void {
    // Implement validation with detailed error feedback
  }

  async execute(): Promise<RunCommandOutput> {
    this.validate();
    // Implement command logic
  }

  getMetadata(): CommandMetadata {
    return RunCommand.metadata;
  }
}
```

### CLI Command Design Rules

1. **Self-Contained**: All related types in one file
2. **Static Metadata**: Accessible before instantiation
3. **Interface Dependencies**: Use service interfaces, not concrete classes
4. **Error-Throwing Validation**: Throw detailed errors, don't return booleans
5. **Dependency Declaration**: Explicitly declare all dependencies in metadata

## Service Development Guidelines

### Service Portability

**CRITICAL**: Services must be portable and reusable across projects:

```typescript
// ✅ CORRECT: Portable service with own configuration interface
export interface ConfigurationServiceConfig {
  configDirectory: string;
  defaultProfile: string;
  validateProfiles: boolean;
}

export class ConfigurationService
  extends BaseService
  implements IConfigurationService
{
  constructor(
    private config: ConfigurationServiceConfig, // Service-specific config
    private logger: ILogger,
    private fileSystemService?: IFileSystemService,
  ) {
    super();
  }
}

// ❌ WRONG: Service coupled to CLI-specific globals
export class BadConfigService extends BaseService {
  constructor(
    private cliGlobals: CLIGlobals, // Project-specific dependency
    private logger: ILogger,
  ) {
    super();
  }
}
```

### Service Design Rules

1. **Own Configuration Interface**: Each service defines its config interface
2. **Constructor Injection**: All dependencies via constructor
3. **Rich Metadata**: Comprehensive metadata for AI-driven selection
4. **Health Monitoring**: Implement isHealthy() method
5. **Graceful Cleanup**: Implement destroy() method

## ServiceRegistry Guidelines

### Access Patterns

**✅ CORRECT: Always use ServiceRegistry as primary access point**

```typescript
// Pattern 1: Constructor Injection (Frequent Users)
export class ProjectService extends BaseService {
  constructor(
    config: ProjectServiceConfig,
    logger: ILogger,
    configurationService: IConfigurationService, // Injected by ServiceRegistry
  ) {
    super();
  }
}

// Pattern 2: Runtime Access (Occasional Users)
export class InstallCommand extends OutputCommand {
  private async setupConfiguration(): Promise<any> {
    const serviceRegistry = ServiceRegistry.getInstance();
    const configService = serviceRegistry.get<IConfigurationService>(
      "IConfigurationService",
    );

    return await configService.createDefaultProfile();
  }
}
```

**❌ WRONG: Direct service instantiation**

```typescript
// NEVER DO THIS
const configService = new ConfigurationService(config, logger);
const projectService = new ProjectService();
```

## CLI-Specific Guidelines

### Command Line Argument Parsing

All CLI commands must use the same argument parsing approach:

```typescript
// ✅ CORRECT: Standardized CLI argument handling
export class RunCommand extends OutputCommand {
  private parseArguments(args: string[]): RunCommandInput {
    const result: RunCommandInput = {
      dangerousSkip: false,
      continue: false,
      profile: undefined,
      passThroughArgs: [],
    };

    const separatorIndex = args.indexOf("--");
    const mcfArgs = separatorIndex >= 0 ? args.slice(0, separatorIndex) : args;
    const passThrough =
      separatorIndex >= 0 ? args.slice(separatorIndex + 1) : [];

    // Parse MCF-specific flags
    for (let i = 0; i < mcfArgs.length; i++) {
      switch (mcfArgs[i]) {
        case "-d":
          result.dangerousSkip = true;
          break;
        case "-c":
          result.continue = true;
          break;
        case "--profile":
          result.profile = mcfArgs[++i];
          break;
      }
    }

    result.passThroughArgs = passThrough;
    return result;
  }
}
```

### Process Management

CLI commands that spawn child processes must follow these patterns:

```typescript
// ✅ CORRECT: Proper child process handling
export class ClaudeService extends BaseService implements IClaudeService {
  async runClaude(options: ClaudeRunOptions): Promise<ClaudeRunResult> {
    const args = this.buildClaudeArguments(options);

    return new Promise((resolve, reject) => {
      const child = spawn("claude", args, {
        stdio: "inherit",
        env: { ...process.env, ...options.environment },
      });

      child.on("exit", (code) => {
        resolve({ exitCode: code || 0 });
      });

      child.on("error", (error) => {
        reject(
          new ClaudeServiceError(`Failed to start Claude: ${error.message}`),
        );
      });
    });
  }
}
```

### Configuration File Management

Configuration handling must be consistent across the CLI:

```typescript
// ✅ CORRECT: Configuration service pattern
export class ConfigurationService extends BaseService {
  async saveProfile(name: string, profile: MCFProfile): Promise<void> {
    const profilePath = path.join(
      this.config.configDirectory,
      "profiles",
      `${name}.json`,
    );

    await this.fileSystemService.ensureDirectory(path.dirname(profilePath));
    await this.fileSystemService.writeJSON(profilePath, profile);

    this.logger.info(`Profile '${name}' saved to ${profilePath}`);
  }

  async loadProfile(name: string): Promise<MCFProfile> {
    const profilePath = path.join(
      this.config.configDirectory,
      "profiles",
      `${name}.json`,
    );

    if (!(await this.fileSystemService.exists(profilePath))) {
      throw new ConfigurationServiceError(
        `Profile '${name}' not found`,
        "PROFILE_NOT_FOUND",
      );
    }

    return await this.fileSystemService.readJSON<MCFProfile>(profilePath);
  }
}
```

## Testing Guidelines

### CLI Command Testing

```typescript
describe("RunCommand", () => {
  let mockClaudeService: jest.Mocked<IClaudeService>;
  let mockConfigService: jest.Mocked<IConfigurationService>;

  beforeEach(() => {
    mockClaudeService = createMockClaudeService();
    mockConfigService = createMockConfigurationService();
  });

  it("should parse dangerous skip flag correctly", async () => {
    const command = new RunCommand(
      { dangerousSkip: true, continue: false },
      mockLogger,
      {
        IClaudeService: mockClaudeService,
        IConfigurationService: mockConfigService,
      },
    );

    await command.execute();
    expect(mockClaudeService.runClaude).toHaveBeenCalledWith(
      expect.objectContaining({ dangerousSkip: true }),
    );
  });

  it("should handle pass-through arguments", async () => {
    const command = new RunCommand(
      { passThroughArgs: ["--debug", "--verbose"] },
      mockLogger,
      {
        IClaudeService: mockClaudeService,
        IConfigurationService: mockConfigService,
      },
    );

    await command.execute();
    expect(mockClaudeService.runClaude).toHaveBeenCalledWith(
      expect.objectContaining({
        additionalArgs: ["--debug", "--verbose"],
      }),
    );
  });
});
```

### Service Testing

```typescript
describe("ConfigurationService", () => {
  it("should implement IConfigurationService contract", () => {
    const config: ConfigurationServiceConfig = {
      configDirectory: "/tmp/mcf-test",
      defaultProfile: "default",
      validateProfiles: true,
    };

    const service = new ConfigurationService(
      config,
      mockLogger,
      mockFileSystemService,
    );
    expect(service).toImplementInterface("IConfigurationService");
  });

  it("should save and load profiles correctly", async () => {
    const service = new ConfigurationService(
      testConfig,
      mockLogger,
      mockFileSystemService,
    );
    const profile: MCFProfile = {
      name: "test",
      claude: { flags: ["-d"], environment: {} },
      mcp: { servers: ["serena"] },
    };

    await service.saveProfile("test", profile);
    const loaded = await service.loadProfile("test");

    expect(loaded).toEqual(profile);
  });
});
```

## File Structure

```
cli/
├── bin/                           # CLI entry point
│   └── mcf.js                    # Main CLI script
├── lib/                          # Core implementation
│   ├── commands/                 # Command micro-blocks
│   │   ├── install/              # Installation commands
│   │   │   └── InstallCommand.ts
│   │   ├── run/                  # Runtime commands
│   │   │   └── RunCommand.ts
│   │   ├── config/               # Configuration commands
│   │   │   └── ConfigCommand.ts
│   │   └── project/              # Project management commands
│   │       └── ProjectCommand.ts
│   ├── services/                 # Infrastructure services
│   │   ├── interfaces/           # Service contracts
│   │   │   ├── IClaudeService.ts
│   │   │   ├── IConfigurationService.ts
│   │   │   └── IFileSystemService.ts
│   │   └── implementations/      # Service implementations
│   │       ├── ClaudeService.ts
│   │       ├── ConfigurationService.ts
│   │       └── FileSystemService.ts
│   └── core/                     # Registry and base classes
│       ├── registry/             # Service and command registries
│       │   ├── ServiceRegistry.ts
│       │   └── CommandRegistry.ts
│       ├── base/                 # Base command and service classes
│       │   ├── BaseCommand.ts
│       │   └── BaseService.ts
│       └── contracts/            # Interfaces and type definitions
│           ├── CommandMetadata.ts
│           └── ServiceMetadata.ts
└── types/                        # Shared type definitions
    ├── MCFProfile.ts
    └── ClaudeOptions.ts
```

## Naming Conventions

1. **Commands**: `{Verb}{Noun}Command` (e.g., `RunCommand`, `ConfigCommand`)
2. **Services**: `{Purpose}Service` (e.g., `ClaudeService`, `ConfigurationService`)
3. **Interfaces**: `I{ServiceName}` (e.g., `IClaudeService`, `IConfigurationService`)
4. **Errors**: `{ComponentName}Error` (e.g., `RunCommandError`, `ClaudeServiceError`)
5. **Types**: `{Purpose}{Input|Output|Config}` (e.g., `RunCommandInput`, `MCFProfile`)
6. **Files**: PascalCase for classes, camelCase for types

## Best Practices Summary

### Do's

1. ✅ Use ES Modules exclusively
2. ✅ Define contracts before implementation
3. ✅ Access services through ServiceRegistry
4. ✅ Keep commands self-contained
5. ✅ Use error-throwing validation
6. ✅ Implement comprehensive metadata
7. ✅ Test contracts and behaviors
8. ✅ Follow naming conventions
9. ✅ Handle child processes properly
10. ✅ Implement proper configuration management

### Don'ts

1. ❌ Never use CommonJS syntax
2. ❌ Never access services directly
3. ❌ Never skip error handling
4. ❌ Never couple services to CLI globals
5. ❌ Never instantiate commands directly
6. ❌ Never return boolean from validation
7. ❌ Never ignore TypeScript errors
8. ❌ Never create circular dependencies
9. ❌ Never bypass the registry system
10. ❌ Never hardcode configuration paths

These standards ensure consistent, maintainable, and scalable development within the MCF CLI micro-block architecture.

---

_Based on UCM Development Guidelines_  
_Source: [development-guidelines.md](https://ucm.utaba.ai/browse/utaba/main/guidance/development/development-guidelines.md)_
