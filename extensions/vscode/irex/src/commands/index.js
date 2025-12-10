const vscode = require('vscode');
const { execFile } = require('child_process');
const util = require('util');
const execFileAsync = util.promisify(execFile);

const fs = require('fs').promises;
const path = require('path');

async function findExecutableOnPath(cmdName) {
    const paths = (process.env.PATH || '').split(path.delimiter);
    const exts = process.platform === 'win32'
        ? (process.env.PATHEXT || '.EXE;.CMD;.BAT;.PS1').split(';')
        : [''];

    for (const dir of paths) {
        if (!dir) continue;
        for (const ext of exts) {
            const full = path.join(dir, cmdName + ext);
            try {
                await fs.access(full);
                return full;
            } catch (e) {
                // not found here
            }
        }
    }
    return null;
}

async function isCliAvailable(cmdName) {
    const exe = await findExecutableOnPath(cmdName);
    return !!exe;
}

async function runIrexCommand(args, cwd) {
    // args: array of strings to pass to irex
    try {
        const proc = await execFileAsync('irex', args, { cwd, timeout: 0, maxBuffer: 1024 * 1024 * 5 });
        return { stdout: proc.stdout, stderr: proc.stderr };
    } catch (err) {
        return { error: err };
    }
}

async function ensureCliInstalled() {
    const ok = await isCliAvailable('irex');
    if (ok) return true;
    const install = 'Install CLI';
    const selection = await vscode.window.showInformationMessage('Irex CLI is not installed. Install it?', install);
    if (selection !== install) return false;
    // Offer link or instructions — open external URL to project's install page (placeholder)
    vscode.env.openExternal(vscode.Uri.parse('https://example.com/irex-install'));
    return false;
}

function registerCommands(context) {
    const disposables = [];

    // irex.init: Initialize Irex in workspace (generic)
    disposables.push(vscode.commands.registerCommand('irex.init', async () => {
        if (!await ensureCliInstalled()) return;
        const cwd = vscode.workspace.workspaceFolders && vscode.workspace.workspaceFolders[0] ? vscode.workspace.workspaceFolders[0].uri.fsPath : undefined;
        const res = await runIrexCommand(['init'], cwd);
        if (res.error) {
            vscode.window.showErrorMessage(`irex init failed: ${res.error.message}`);
            return;
        }
        vscode.window.showInformationMessage('Irex initialized.');
    }));

    // irex.initNodeTs: interactive init for node-ts (choose flags or open terminal)
    disposables.push(vscode.commands.registerCommand('irex.initNodeTs', async () => {
        if (!await ensureCliInstalled()) return;
        const cwd = vscode.workspace.workspaceFolders && vscode.workspace.workspaceFolders[0] ? vscode.workspace.workspaceFolders[0].uri.fsPath : undefined;

        const pick = await vscode.window.showQuickPick([
            { label: 'Run with flags', description: 'Provide flags/options non-interactively' },
            { label: 'Run in Terminal', description: 'Open an integrated terminal for interactive prompts' }
        ], { placeHolder: 'How would you like to run `irex init node-ts`?' });

        if (!pick) return; // cancelled

        if (pick.label === 'Run in Terminal') {
            const term = vscode.window.createTerminal({ name: 'Irex: init node-ts', cwd });
            term.show(true);
            term.sendText('irex init node-ts', true);
            vscode.window.showInformationMessage('Opened terminal to run `irex init node-ts`.');
            return;
        }

        // Run with flags: ask user for flags string, then execute
        const flagsInput = await vscode.window.showInputBox({ prompt: 'Enter flags for `irex init node-ts` (e.g. --flag value)', placeHolder: '--example-flag value' });
        if (typeof flagsInput === 'undefined') return; // cancelled

        // split args respecting simple quoting
        const extraArgs = flagsInput.trim().length > 0 ? flagsInput.match(/(?:\".*?\"|[^\s]+)+/g).map(s => s.replace(/^"|"$/g, '')) : [];
        const args = ['init', 'node-ts', ...extraArgs];

        const res = await runIrexCommand(args, cwd);
        if (res.error) {
            vscode.window.showErrorMessage(`irex init node-ts failed: ${res.error.message}`);
            return;
        }

        const stdout = res.stdout ? String(res.stdout).trim() : '';
        vscode.window.showInformationMessage(stdout.length ? `Irex: ${stdout}` : 'Irex: Initialized node TypeScript project.');
    }));

    // irex.watch: start development server
    disposables.push(vscode.commands.registerCommand('irex.watch', async () => {
        if (!await ensureCliInstalled()) return;
        const cwd = vscode.workspace.workspaceFolders && vscode.workspace.workspaceFolders[0] ? vscode.workspace.workspaceFolders[0].uri.fsPath : undefined;
        const res = await runIrexCommand(['watch'], cwd);
        if (res.error) {
            vscode.window.showErrorMessage(`irex watch failed: ${res.error.message}`);
            return;
        }
        vscode.window.showInformationMessage('Irex: Development server started.');
    }));

    disposables.forEach(d => context.subscriptions.push(d));
}

module.exports = { registerCommands, runIrexCommand, ensureCliInstalled };
