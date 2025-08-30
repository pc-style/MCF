const { spawn } = require('child_process');
const path = require('path');
const chalk = require('chalk');
const ora = require('ora');

module.exports = async function install(options) {
  console.log(chalk.blue.bold('🚀 MCF Framework Installer'));
  console.log();
  
  const spinner = ora('Downloading and installing MCF...').start();
  
  try {
    // Get the install.sh script path from the repository root
    const repoRoot = path.resolve(__dirname, '../../');
    const installScript = path.join(repoRoot, 'install.sh');
    
    // Build install command arguments
    const args = [];
    if (options.yes) {
      args.push('--yes');
    }
    
    // Execute the install.sh script
    const installProcess = spawn('bash', [installScript, ...args], {
      stdio: 'pipe',
      cwd: process.cwd()
    });
    
    let output = '';
    let errorOutput = '';
    
    installProcess.stdout.on('data', (data) => {
      output += data.toString();
      // Show real-time output for important messages
      const lines = data.toString().split('\n');
      lines.forEach(line => {
        if (line.trim() && (line.includes('INFO:') || line.includes('SUCCESS:') || line.includes('WARN:'))) {
          spinner.text = line.replace(/.*INFO:\s*/, '').replace(/.*SUCCESS:\s*/, '').replace(/.*WARN:\s*/, '');
        }
      });
    });
    
    installProcess.stderr.on('data', (data) => {
      errorOutput += data.toString();
    });
    
    installProcess.on('close', (code) => {
      if (code === 0) {
        spinner.succeed('MCF installation completed successfully!');
        console.log();
        console.log(chalk.green('✅ MCF is now installed and ready to use.'));
        console.log();
        console.log(chalk.blue('Next steps:'));
        console.log('  • Run', chalk.yellow('mcf setup'), 'to configure MCF');
        console.log('  • Run', chalk.yellow('mcf run'), 'to start a MCF session');
        console.log('  • Run', chalk.yellow('mcf status'), 'to check installation status');
      } else {
        spinner.fail('MCF installation failed');
        console.error(chalk.red('❌ Installation failed with exit code:'), code);
        if (errorOutput) {
          console.error(chalk.red('\nError output:'));
          console.error(errorOutput);
        }
        process.exit(code);
      }
    });
    
  } catch (error) {
    spinner.fail('Installation failed');
    console.error(chalk.red('❌ Error during installation:'), error.message);
    process.exit(1);
  }
};