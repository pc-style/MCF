import fs from 'fs';
import path from 'path';
import os from 'os';
import chalk from 'chalk';

export default async function status() {
  console.log(chalk.blue.bold('📊 MCF Status Check'));
  console.log();
  
  const mcfDir = path.join(os.homedir(), 'mcf');
  const settingsFile = path.join(mcfDir, '.claude', 'settings.json');
  const claudeMcfScript = path.join(mcfDir, 'claude-mcf.sh');
  const templatesDir = path.join(mcfDir, 'templates');
  const hooksDir = path.join(mcfDir, '.claude', 'hooks');
  const scriptsDir = path.join(mcfDir, 'scripts');
  
  let allGood = true;
  
  // Check MCF installation
  console.log(chalk.blue('🔍 Installation Status'));
  if (fs.existsSync(mcfDir)) {
    console.log(chalk.green('  ✅ MCF directory found'), chalk.gray(`(${mcfDir})`));
  } else {
    console.log(chalk.red('  ❌ MCF directory not found'));
    console.log(chalk.blue('     Run'), chalk.yellow('mcf install'), chalk.blue('to install MCF'));
    allGood = false;
  }
  
  // Check core files
  console.log(chalk.blue('🔍 Core Files'));
  const coreFiles = [
    { path: claudeMcfScript, name: 'Main runner script' },
    { path: settingsFile, name: 'Settings configuration' },
    { path: path.join(scriptsDir, 'template-engine.py'), name: 'Template engine' }
  ];
  
  coreFiles.forEach(({ path: filePath, name }) => {
    if (fs.existsSync(filePath)) {
      console.log(chalk.green(`  ✅ ${name} found`));
    } else {
      console.log(chalk.red(`  ❌ ${name} missing`));
      allGood = false;
    }
  });
  
  // Check directories
  console.log(chalk.blue('🔍 Directory Structure'));
  const directories = [
    { path: templatesDir, name: 'Templates directory' },
    { path: hooksDir, name: 'Hooks directory' },
    { path: scriptsDir, name: 'Scripts directory' }
  ];
  
  directories.forEach(({ path: dirPath, name }) => {
    if (fs.existsSync(dirPath)) {
      const items = fs.readdirSync(dirPath);
      console.log(chalk.green(`  ✅ ${name} found`), chalk.gray(`(${items.length} items)`));
    } else {
      console.log(chalk.red(`  ❌ ${name} missing`));
      allGood = false;
    }
  });
  
  // Check configuration
  console.log(chalk.blue('🔍 Configuration'));
  if (fs.existsSync(settingsFile)) {
    try {
      const settings = JSON.parse(fs.readFileSync(settingsFile, 'utf8'));
      
      // Check hooks configuration
      if (settings.hooks) {
        const hookTypes = Object.keys(settings.hooks);
        console.log(chalk.green('  ✅ Hooks system configured'), chalk.gray(`(${hookTypes.length} hook types)`));
      } else {
        console.log(chalk.yellow('  ⚠️  Hooks system not configured'));
      }
      
      // Check status line
      if (settings.statusLine) {
        console.log(chalk.green('  ✅ Status line enabled'));
      } else {
        console.log(chalk.yellow('  ⚠️  Status line disabled'));
      }
      
      // Check output style
      if (settings.outputStyle) {
        console.log(chalk.green(`  ✅ Output style: ${settings.outputStyle}`));
      } else {
        console.log(chalk.yellow('  ⚠️  Output style not set'));
      }
      
    } catch (error) {
      console.log(chalk.red('  ❌ Settings file corrupted'));
      allGood = false;
    }
  } else {
    console.log(chalk.yellow('  ⚠️  No configuration found'));
    console.log(chalk.blue('     Run'), chalk.yellow('mcf setup'), chalk.blue('to configure MCF'));
  }
  
  // Check templates
  console.log(chalk.blue('🔍 Templates'));
  if (fs.existsSync(templatesDir)) {
    const templates = fs.readdirSync(templatesDir).filter(f => f.endsWith('.json'));
    if (templates.length > 0) {
      console.log(chalk.green(`  ✅ ${templates.length} templates available`));
    } else {
      console.log(chalk.yellow('  ⚠️  No templates found'));
    }
  }
  
  // Check shell integration
  console.log(chalk.blue('🔍 Shell Integration'));
  const localBinDir = path.join(os.homedir(), '.local', 'bin');
  const claudeMcfLink = path.join(localBinDir, 'claude-mcf');
  
  if (fs.existsSync(claudeMcfLink)) {
    try {
      const linkTarget = fs.readlinkSync(claudeMcfLink);
      if (linkTarget === claudeMcfScript) {
        console.log(chalk.green('  ✅ Shell integration configured'));
      } else {
        console.log(chalk.yellow('  ⚠️  Shell integration points to wrong target'));
      }
    } catch {
      console.log(chalk.yellow('  ⚠️  Shell integration link broken'));
    }
  } else {
    console.log(chalk.yellow('  ⚠️  Shell integration not configured'));
    console.log(chalk.blue('     Run'), chalk.yellow('mcf setup'), chalk.blue('to configure shell integration'));
  }
  
  // Overall status
  console.log();
  if (allGood) {
    console.log(chalk.green.bold('✅ MCF is fully operational!'));
    console.log();
    console.log(chalk.blue('Ready to use:'));
    console.log('  • Run', chalk.yellow('mcf run'), 'to start a MCF session');
    console.log('  • Run', chalk.yellow('mcf templates'), 'to manage templates');
  } else {
    console.log(chalk.red.bold('❌ MCF has some issues that need attention'));
    console.log();
    console.log(chalk.blue('Recommended actions:'));
    if (!fs.existsSync(mcfDir)) {
      console.log('  • Run', chalk.yellow('mcf install'), 'to install MCF');
    } else {
      console.log('  • Run', chalk.yellow('mcf install'), 'to reinstall/repair MCF');
    }
    if (!fs.existsSync(settingsFile)) {
      console.log('  • Run', chalk.yellow('mcf setup'), 'to configure MCF');
    }
  }
};