import { CommandMetadata } from "../contracts/CommandMetadata.js";

export interface BaseCommand<
  InputType = unknown,
  OutputType = unknown,
  ErrorType = Error,
> {
  // Metadata for the command
  metadata: CommandMetadata<InputType, OutputType, ErrorType>;

  // Execute the command with given arguments
  execute(input: InputType): Promise<OutputType>;

  // Validate input before execution
  validate(input: InputType): Promise<boolean>;

  // Get help text for the command
  help(): string;

  // Get description of the command
  description(): string;

  // Get current command metadata (useful for dynamic metadata)
  getMetadata(): CommandMetadata<InputType, OutputType, ErrorType>;
}
