import { Type } from "../../types/CLITypes.js";

export interface CommandOptionMetadata {
  name: string;
  description: string;
  type: "string" | "boolean" | "number";
  required?: boolean;
  default?: unknown;
  aliases?: string[];
}

export interface CommandDependency {
  services?: string[];
  commands?: string[];
  external?: string[];
}

export interface CommandMetadata<
  InputType = unknown,
  OutputType = unknown,
  ErrorType = Error,
> {
  name: string;
  description: string;
  category?: string;

  // Versioning
  version: string;
  contractVersion: string;

  // Type information
  inputType?: Type<InputType>;
  outputType?: Type<OutputType>;
  errorType?: Type<ErrorType>;

  // Dependencies
  dependencies?: CommandDependency;

  // CLI-specific fields
  options?: CommandOptionMetadata[];
  aliases?: string[];
  flags?: string[];

  // Help and documentation
  helpText?: string;
  examples?: string[];
}
