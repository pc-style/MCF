#!/usr/bin/env node

const { program } = require('commander');
const chalk = require('chalk');
const pkg = require('../package.json');

// Import command modules
const install = require('../lib/install');
const setup = require('../lib/setup');
const run = require('../lib/run');
const templates = require('../lib/templates');
const status = require('../lib/status');

program
  .name('mcf')
  .description('MCF (Multi Component Framework) CLI - Installation, configuration and setup tool')
  .version(pkg.version);

// Install command
program
  .command('install')
  .description('Install MCF framework')
  .option('-y, --yes', 'Skip interactive prompts and proceed automatically')
  .action(install);

// Setup command  
program
  .command('setup')
  .description('Configure MCF after installation')
  .action(setup);

// Run command
program
  .command('run')
  .description('Start MCF session')
  .action(run);

// Templates command
program
  .command('templates')
  .alias('t')
  .description('Manage MCF templates')
  .argument('[action]', 'Action to perform (list, init, info)')
  .argument('[name]', 'Template name')
  .action(templates);

// Status command
program
  .command('status')
  .description('Check MCF installation status')
  .action(status);

// Handle unknown commands
program
  .configureOutput({
    writeErr: (str) => process.stderr.write(chalk.red(str))
  });

program.parse();

// Show help if no command provided
if (!process.argv.slice(2).length) {
  program.outputHelp();
}