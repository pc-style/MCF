#!/usr/bin/env node

/**
 * Test script for MCF CLI standalone version
 */

import { spawn } from 'child_process';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

console.log('🧪 Testing MCF CLI Standalone...');

const mcfScript = path.join(__dirname, 'mcf-standalone.js');

async function runTest(name, args, expectedInOutput) {
  return new Promise((resolve, reject) => {
    console.log(`\n🔍 Testing: ${name}`);
    
    const child = spawn('node', [mcfScript, ...args], { stdio: 'pipe' });
    let output = '';
    let errorOutput = '';
    
    child.stdout.on('data', (data) => {
      output += data.toString();
    });
    
    child.stderr.on('data', (data) => {
      errorOutput += data.toString();
    });
    
    child.on('close', (code) => {
      const allOutput = output + errorOutput;
      
      if (allOutput.includes(expectedInOutput)) {
        console.log('✅ Test passed');
        resolve(true);
      } else {
        console.log('❌ Test failed');
        console.log(`Expected: "${expectedInOutput}"`);
        console.log(`Got output: "${allOutput.substring(0, 200)}..."`);
        resolve(false);
      }
    });
    
    child.on('error', (error) => {
      console.log('❌ Test error:', error.message);
      resolve(false);
    });
  });
}

async function runTests() {
  const tests = [
    { name: 'Help command', args: ['--help'], expected: 'MCF (My Claude Flow) CLI' },
    { name: 'Version command', args: ['--version'], expected: '1.0.0' },
    { name: 'Config list', args: ['config', 'list'], expected: 'MCF Configuration Profiles' },
    { name: 'Config show mcf', args: ['config', 'show', 'mcf'], expected: 'Config Directory' },
    { name: 'Project list', args: ['project', 'list'], expected: 'MCF Projects' },
    { name: 'Status command', args: ['status'], expected: 'MCF Status Check' }
  ];
  
  let passed = 0;
  let failed = 0;
  
  for (const test of tests) {
    const result = await runTest(test.name, test.args, test.expected);
    if (result) {
      passed++;
    } else {
      failed++;
    }
  }
  
  console.log(`\n📊 Test Results: ${passed} passed, ${failed} failed`);
  
  if (failed === 0) {
    console.log('\n🎉 All tests passed!');
    console.log('\nStandalone MCF CLI is ready!');
    console.log('File: cli/mcf-standalone.js');
    console.log('Size: ~24KB, 854 lines');
    console.log('\nUsage:');
    console.log('• node cli/mcf-standalone.js run --config mcf');
    console.log('• node cli/mcf-standalone.js config list');
    process.exit(0);
  } else {
    console.log('\n❌ Some tests failed');
    process.exit(1);
  }
}

runTests();
