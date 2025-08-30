const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');
const chalk = require('chalk');
const ora = require('ora');

module.exports = async function run() {
  console.log(chalk.blue.bold('🚀 Starting MCF Session'));
  console.log();
  
  const mcfDir = path.join(os.homedir(), 'mcf');
  const claudeMcfScript = path.join(mcfDir, 'claude-mcf.sh');
  const settingsFile = path.join(mcfDir, '.claude', 'settings.json');
  
  // Check if MCF is installed
  if (!fs.existsSync(mcfDir)) {
    console.log(chalk.red('❌ MCF is not installed.'));
    console.log(chalk.blue('Run'), chalk.yellow('mcf install'), chalk.blue('first.'));
    return;
  }
  
  if (!fs.existsSync(claudeMcfScript)) {
    console.log(chalk.red('❌ MCF runner script not found.'));
    console.log(chalk.blue('Try running'), chalk.yellow('mcf install'), chalk.blue('to reinstall.'));
    return;
  }
  
  if (!fs.existsSync(settingsFile)) {
    console.log(chalk.yellow('⚠️  MCF is not configured.'));
    console.log(chalk.blue('Run'), chalk.yellow('mcf setup'), chalk.blue('to configure MCF first.'));
    return;
  }
  
  const spinner = ora('Preparing MCF session...').start();
  
  try {
    // Check if this is first run
    const bookmarksDir = path.join(mcfDir, '.claude', 'bookmarks');
    const firstRunFile = path.join(bookmarksDir, '.first-run.txt');
    const isFirstRun = !fs.existsSync(firstRunFile);
    
    if (isFirstRun) {
      spinner.info('First run detected - authentication may be required');
      console.log(chalk.blue('🔐 You may need to authenticate with Claude'));
    }
    
    spinner.text = 'Starting Claude with MCF configuration...';
    
    // Execute the claude-mcf.sh script
    const runProcess = spawn('bash', [claudeMcfScript], {
      stdio: 'inherit', // Pass through stdin/stdout/stderr to allow interactive use
      cwd: process.cwd(),
      env: {
        ...process.env,
        MCF_CLI_MODE: 'true' // Let the script know it's being run via CLI
      }
    });
    
    runProcess.on('close', (code) => {
      if (code === 0) {
        console.log();
        console.log(chalk.green('✅ MCF session completed successfully'));
      } else {
        console.log();
        console.log(chalk.red(`❌ MCF session ended with exit code: ${code}`));
        if (code === 1) {
          console.log(chalk.blue('This might be normal if you exited Claude manually.'));
        }
      }
    });
    
    runProcess.on('error', (error) => {
      spinner.fail('Failed to start MCF session');
      console.error(chalk.red('❌ Error starting MCF:'), error.message);
      
      // Provide helpful error messages
      if (error.code === 'ENOENT') {
        console.log();
        console.log(chalk.blue('Possible solutions:'));
        console.log('  • Make sure bash is installed and in your PATH');
        console.log('  • Try running', chalk.yellow('mcf install'), 'to reinstall MCF');
      }
      
      process.exit(1);
    });
    
    // Clear spinner when process starts successfully
    setTimeout(() => {
      spinner.stop();
    }, 1000);
    
  } catch (error) {
    spinner.fail('Failed to start MCF session');
    console.error(chalk.red('❌ Error:'), error.message);
    process.exit(1);
  }
};