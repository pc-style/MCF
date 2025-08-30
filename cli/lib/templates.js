const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');
const os = require('os');
const chalk = require('chalk');
const ora = require('ora');

module.exports = async function templates(action, name) {
  console.log(chalk.blue.bold('📚 MCF Templates'));
  console.log();
  
  const mcfDir = path.join(os.homedir(), 'mcf');
  const scriptsDir = path.join(mcfDir, 'scripts');
  const templateEngine = path.join(scriptsDir, 'template-engine.py');
  
  // Check if MCF is installed
  if (!fs.existsSync(mcfDir)) {
    console.log(chalk.red('❌ MCF is not installed.'));
    console.log(chalk.blue('Run'), chalk.yellow('mcf install'), chalk.blue('first.'));
    return;
  }
  
  if (!fs.existsSync(templateEngine)) {
    console.log(chalk.red('❌ Template engine not found.'));
    console.log(chalk.blue('Try running'), chalk.yellow('mcf install'), chalk.blue('to reinstall.'));
    return;
  }
  
  // Default to 'list' if no action specified
  if (!action) {
    action = 'list';
  }
  
  // Validate action
  const validActions = ['list', 'init', 'info'];
  if (!validActions.includes(action)) {
    console.log(chalk.red(`❌ Invalid action: ${action}`));
    console.log(chalk.blue('Valid actions:'), validActions.join(', '));
    return;
  }
  
  // Validate name for actions that require it
  if (['init', 'info'].includes(action) && !name) {
    console.log(chalk.red(`❌ Template name is required for '${action}' action`));
    console.log(chalk.blue('Usage:'), chalk.yellow(`mcf templates ${action} <template-name>`));
    return;
  }
  
  const spinner = ora(`${action === 'list' ? 'Listing' : action === 'init' ? 'Initializing from' : 'Getting info for'} templates...`).start();
  
  try {
    // Build command arguments
    const args = [templateEngine, action];
    if (name) {
      args.push(name);
    }
    
    // Execute the template-engine.py script
    const templateProcess = spawn('python3', args, {
      stdio: 'pipe',
      cwd: process.cwd()
    });
    
    let output = '';
    let errorOutput = '';
    
    templateProcess.stdout.on('data', (data) => {
      output += data.toString();
    });
    
    templateProcess.stderr.on('data', (data) => {
      errorOutput += data.toString();
    });
    
    templateProcess.on('close', (code) => {
      spinner.stop();
      
      if (code === 0) {
        // Success - show the output
        if (output.trim()) {
          console.log(output.trim());
        }
        
        // Add helpful next steps
        if (action === 'list') {
          console.log();
          console.log(chalk.blue('💡 Next steps:'));
          console.log('  • Run', chalk.yellow('mcf templates info <name>'), 'to see template details');
          console.log('  • Run', chalk.yellow('mcf templates init <name>'), 'to initialize from template');
        } else if (action === 'init') {
          console.log();
          console.log(chalk.green(`✅ Template '${name}' initialized successfully!`));
        }
        
      } else {
        console.log(chalk.red(`❌ Template operation failed with exit code: ${code}`));
        
        if (errorOutput.trim()) {
          console.error(chalk.red('Error:'), errorOutput.trim());
        } else if (output.trim()) {
          console.log(output.trim());
        }
        
        // Provide helpful suggestions
        if (action === 'init' || action === 'info') {
          console.log();
          console.log(chalk.blue('💡 Suggestions:'));
          console.log('  • Run', chalk.yellow('mcf templates'), 'to see available templates');
          console.log('  • Check that the template name is spelled correctly');
        }
      }
    });
    
    templateProcess.on('error', (error) => {
      spinner.fail('Failed to run template engine');
      console.error(chalk.red('❌ Error:'), error.message);
      
      if (error.code === 'ENOENT') {
        console.log();
        console.log(chalk.blue('Possible solutions:'));
        console.log('  • Make sure Python 3 is installed and in your PATH');
        console.log('  • Try running', chalk.yellow('mcf install'), 'to reinstall MCF');
      }
      
      process.exit(1);
    });
    
  } catch (error) {
    spinner.fail('Failed to manage templates');
    console.error(chalk.red('❌ Error:'), error.message);
    process.exit(1);
  }
};