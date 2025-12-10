const vscode = require('vscode');
const fs = require('fs').promises;
const { execFile } = require('child_process');
const util = require('util');
const execFileAsync = util.promisify(execFile);

async function provideFormattingEdits(document, range) {
    const filePath = document.uri.fsPath;
    try {
        await execFileAsync('irex', ['format', filePath, '-w'], { timeout: 10000 });
    } catch (err) {
        vscode.window.showErrorMessage(`irex format failed: ${err.message}`);
        return [];
    }

    let formatted;
    try {
        formatted = await fs.readFile(filePath, 'utf8');
    } catch (err) {
        vscode.window.showErrorMessage(`Failed to read formatted file: ${err.message}`);
        return [];
    }

    const fullRange = range || new vscode.Range(
        new vscode.Position(0, 0),
        document.lineAt(document.lineCount - 1).range.end
    );

    return [vscode.TextEdit.replace(fullRange, formatted)];
}

function registerFormattingProviders(context) {
    const formattingProvider = {
        provideDocumentFormattingEdits: async (document, options, token) => {
            return provideFormattingEdits(document);
        }
    };

    const rangeFormattingProvider = {
        provideDocumentRangeFormattingEdits: async (document, range, options, token) => {
            return provideFormattingEdits(document, range);
        }
    };

    const disposables = [];
    disposables.push(vscode.languages.registerDocumentFormattingEditProvider('specifications-hcl', formattingProvider));
    disposables.push(vscode.languages.registerDocumentRangeFormattingEditProvider('specifications-hcl', rangeFormattingProvider));

    disposables.forEach(d => context.subscriptions.push(d));
}

module.exports = {
    registerFormattingProviders
};
