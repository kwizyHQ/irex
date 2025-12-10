// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
const vscode = require('vscode');
const { registerFormattingProviders } = require('./src/formatter');
const { registerCommands } = require('./src/commands');

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed

/**
 * @param {vscode.ExtensionContext} context
 */
function activate(context) {

	// Register formatting providers from src/formatter
	registerFormattingProviders(context);

	// Register CLI commands
	registerCommands(context);
}

// This method is called when your extension is deactivated
function deactivate() {}

module.exports = {
	activate,
	deactivate
}
