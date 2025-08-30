#!/usr/bin/env node

/**
 * MCF CLI - Entry point
 * 
 * This module provides programmatic access to MCF CLI functionality.
 * For command-line usage, use bin/mcf.js
 */

const install = require('./lib/install');
const setup = require('./lib/setup');
const run = require('./lib/run');
const templates = require('./lib/templates');
const status = require('./lib/status');

module.exports = {
  install,
  setup,
  run,
  templates,
  status
};

// If called directly, delegate to CLI
if (require.main === module) {
  require('./bin/mcf.js');
}