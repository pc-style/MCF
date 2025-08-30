const fs = require('fs');
const path = require('path');
const os = require('os');
const chalk = require('chalk');
const inquirer = require('inquirer');
const ora = require('ora');

module.exports = async function setup() {
  console.log(chalk.blue.bold('🔧 MCF Setup & Configuration'));
  console.log();
  
  const mcfDir = path.join(os.homedir(), 'mcf');
  const settingsFile = path.join(mcfDir, '.claude', 'settings.json');
  
  // Check if MCF is installed
  if (!fs.existsSync(mcfDir)) {
    console.log(chalk.red('❌ MCF is not installed.'));
    console.log(chalk.blue('Run'), chalk.yellow('mcf install'), chalk.blue('first.'));
    return;
  }
  
  console.log(chalk.green('✅ MCF installation found'));
  console.log();
  
  // Check current configuration
  let currentSettings = {};
  if (fs.existsSync(settingsFile)) {
    try {
      currentSettings = JSON.parse(fs.readFileSync(settingsFile, 'utf8'));
      console.log(chalk.blue('Current configuration found:'));
      console.log(chalk.gray(`  Settings file: ${settingsFile}`));
      
      // Show some key settings
      if (currentSettings.hooks) {
        console.log(chalk.gray(`  Hooks configured: ${Object.keys(currentSettings.hooks).length} types`));
      }
      if (currentSettings.statusLine) {
        console.log(chalk.gray('  Status line: enabled'));
      }
      console.log();
    } catch (error) {
      console.log(chalk.yellow('⚠️  Warning: Could not read current settings'));
    }
  }
  
  // Interactive setup questions
  const answers = await inquirer.prompt([
    {
      type: 'confirm',
      name: 'enableHooks',
      message: 'Enable MCF intelligent hooks system?',
      default: true
    },
    {
      type: 'confirm',
      name: 'enableStatusLine',
      message: 'Enable enhanced status line?',
      default: true
    },
    {
      type: 'list',
      name: 'outputStyle',
      message: 'Choose output style:',
      choices: [
        { name: 'Explanatory (recommended)', value: 'explanatory' },
        { name: 'Concise', value: 'concise' },
        { name: 'Minimal', value: 'minimal' }
      ],
      default: 'explanatory'
    },
    {
      type: 'confirm',
      name: 'setupShellIntegration',
      message: 'Add MCF to your shell PATH?',
      default: true
    }
  ]);
  
  const spinner = ora('Configuring MCF...').start();
  
  try {
    // Update settings
    const newSettings = {
      ...currentSettings,
      outputStyle: answers.outputStyle
    };
    
    // Configure hooks
    if (answers.enableHooks && currentSettings.hooks) {
      newSettings.hooks = currentSettings.hooks;
      spinner.text = 'Hooks system enabled';
    } else if (!answers.enableHooks) {
      delete newSettings.hooks;
      spinner.text = 'Hooks system disabled';
    }
    
    // Configure status line
    if (answers.enableStatusLine && currentSettings.statusLine) {
      newSettings.statusLine = currentSettings.statusLine;
      spinner.text = 'Status line enabled';
    } else if (!answers.enableStatusLine) {
      delete newSettings.statusLine;
      spinner.text = 'Status line disabled';
    }
    
    // Write updated settings
    fs.writeFileSync(settingsFile, JSON.stringify(newSettings, null, 2));
    
    // Set up shell integration if requested
    if (answers.setupShellIntegration) {
      const localBinDir = path.join(os.homedir(), '.local', 'bin');
      const claudeMcfScript = path.join(mcfDir, 'claude-mcf.sh');
      const linkPath = path.join(localBinDir, 'claude-mcf');
      
      // Create ~/.local/bin if it doesn't exist
      if (!fs.existsSync(localBinDir)) {
        fs.mkdirSync(localBinDir, { recursive: true });
      }
      
      // Create symlink if script exists
      if (fs.existsSync(claudeMcfScript)) {
        try {
          if (fs.existsSync(linkPath)) {
            fs.unlinkSync(linkPath);
          }
          fs.symlinkSync(claudeMcfScript, linkPath);
          spinner.text = 'Shell integration configured';
        } catch (error) {
          console.log(chalk.yellow('⚠️  Could not set up shell integration'));
        }
      }
    }
    
    spinner.succeed('MCF configuration completed!');
    console.log();
    console.log(chalk.green('✅ MCF is now configured and ready to use.'));
    console.log();
    console.log(chalk.blue('You can now:'));
    console.log('  • Run', chalk.yellow('mcf run'), 'to start a MCF session');
    console.log('  • Run', chalk.yellow('mcf templates'), 'to manage templates');
    console.log('  • Run', chalk.yellow('mcf status'), 'to check system status');
    
    if (answers.setupShellIntegration) {
      console.log();
      console.log(chalk.blue('Shell integration:'));
      console.log('  • Make sure', chalk.yellow('~/.local/bin'), 'is in your PATH');
      console.log('  • You can run', chalk.yellow('claude-mcf'), 'directly from any directory');
    }
    
  } catch (error) {
    spinner.fail('Configuration failed');
    console.error(chalk.red('❌ Error during setup:'), error.message);
    process.exit(1);
  }
};