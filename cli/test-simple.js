#!/usr/bin/env node

/**
 * Simple test script for MCF CLI
 */

const { spawn } = require('child_process');
const path = require('path');

console.log('🧪 Testing MCF CLI...');

const mcfScript = path.join(__dirname, 'bin', 'mcf.js');

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
    { name: 'Help command', args: ['--help'], expected: 'MCF (Multi Component Framework) CLI' },
    { name: 'Version command', args: ['--version'], expected: '1.0.0' },
    { name: 'Status command', args: ['status'], expected: '📊 MCF Status Check' },
    { name: 'Templates command', args: ['templates'], expected: '📚 MCF Templates' }
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
    console.log('\nCLI is ready for use:');
    console.log('• npm install -g @pc-style/mcf-cli');
    console.log('• npx @pc-style/mcf-cli <command>');
    process.exit(0);
  } else {
    console.log('\n❌ Some tests failed');
    process.exit(1);
  }
}

runTests();