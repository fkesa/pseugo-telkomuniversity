const vscode = require('vscode');
const { execFile, exec } = require('child_process');
const path = require('path');
const fs = require('fs');

let outputChannel;

function getOutputChannel() {
  if (!outputChannel) {
    outputChannel = vscode.window.createOutputChannel('Pseugo Pseudocode');
  }
  return outputChannel;
}

function findPseugoBinary(activeDocPath) {
  // 1. Check workspace folders
  const wsFolders = vscode.workspace.workspaceFolders;
  if (wsFolders) {
    for (const folder of wsFolders) {
      const candidates = [
        path.join(folder.uri.fsPath, 'pseugo.exe'),
        path.join(folder.uri.fsPath, 'pseugo'),
        path.join(folder.uri.fsPath, 'bin', 'pseugo.exe'),
        path.join(folder.uri.fsPath, 'bin', 'pseugo')
      ];
      for (const c of candidates) {
        if (fs.existsSync(c)) {
          return c;
        }
      }
    }
  }

  // 2. Check active file parent
  if (activeDocPath) {
    let dir = path.dirname(activeDocPath);
    for (let i = 0; i < 4; i++) {
      const candidates = [
        path.join(dir, 'pseugo.exe'),
        path.join(dir, 'pseugo'),
        path.join(dir, 'bin', 'pseugo.exe'),
        path.join(dir, 'bin', 'pseugo')
      ];
      for (const c of candidates) {
        if (fs.existsSync(c)) {
          return c;
        }
      }
      const parent = path.dirname(dir);
      if (parent === dir) break;
      dir = parent;
    }
  }

  return 'pseugo';
}

function activate(context) {
  const channel = getOutputChannel();

  // Command: Check / Typecheck
  const checkCmd = vscode.commands.registerCommand('pseugo.check', () => {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
      vscode.window.showWarningMessage('Buka file pseudocode terlebih dahulu.');
      return;
    }

    const doc = editor.document;
    doc.save().then(() => {
      const filePath = doc.fileName;
      const bin = findPseugoBinary(filePath);

      channel.clear();
      channel.show(true);
      channel.appendLine(`[Pseugo] Memeriksa ${filePath}...`);

      execFile(bin, ['check', filePath], (err, stdout, stderr) => {
        if (err) {
          channel.appendLine(stderr || stdout || err.message);
          vscode.window.showErrorMessage(`[Pseugo] Validasi gagal! Periksa tab Output 'Pseugo Pseudocode'.`);
        } else {
          channel.appendLine(stdout || '[OK] Typecheck berhasil!');
          vscode.window.showInformationMessage(`[Pseugo] ${path.basename(filePath)} valid!`);
        }
      });
    });
  });

  // Command: Convert to Golang
  const convertCmd = vscode.commands.registerCommand('pseugo.convertToGolang', () => {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
      vscode.window.showWarningMessage('Buka file pseudocode terlebih dahulu.');
      return;
    }

    const doc = editor.document;
    doc.save().then(() => {
      const filePath = doc.fileName;
      const defaultOut = filePath.replace(/\.(pseudo|pseudocode|txt)$/i, '') + '.go';
      const bin = findPseugoBinary(filePath);

      channel.clear();
      channel.show(true);
      channel.appendLine(`[Pseugo] Mengonversi ${filePath} ke ${defaultOut}...`);

      execFile(bin, ['convert', '-to-golang', '-i', filePath, '-o', defaultOut], (err, stdout, stderr) => {
        if (err) {
          channel.appendLine(stderr || stdout || err.message);
          vscode.window.showErrorMessage(`[Pseugo] Konversi gagal! Periksa tab Output.`);
        } else {
          channel.appendLine(stdout || `Berhasil dikonversi ke ${defaultOut}`);
          vscode.window.showInformationMessage(`[Pseugo] Berhasil dikonversi ke ${path.basename(defaultOut)}`);
          vscode.workspace.openTextDocument(defaultOut).then(newDoc => {
            vscode.window.showTextDocument(newDoc, vscode.ViewColumn.Beside);
          });
        }
      });
    });
  });

  // Command: Convert & Run with Go
  const runCmd = vscode.commands.registerCommand('pseugo.runGolang', () => {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
      vscode.window.showWarningMessage('Buka file pseudocode terlebih dahulu.');
      return;
    }

    const doc = editor.document;
    doc.save().then(() => {
      const filePath = doc.fileName;
      const defaultOut = filePath.replace(/\.(pseudo|pseudocode|txt)$/i, '') + '.go';
      const bin = findPseugoBinary(filePath);

      execFile(bin, ['convert', '-to-golang', '-i', filePath, '-o', defaultOut], (err, stdout, stderr) => {
        if (err) {
          const ch = getOutputChannel();
          ch.show(true);
          ch.appendLine(stderr || stdout || err.message);
          vscode.window.showErrorMessage(`[Pseugo] Konversi gagal sebelum menjalankan kode Go.`);
          return;
        }

        let terminal = vscode.window.terminals.find(t => t.name === 'Pseugo Runner');
        if (!terminal) {
          terminal = vscode.window.createTerminal('Pseugo Runner');
        }
        terminal.show();
        terminal.sendText(`go run "${defaultOut}"`);
      });
    });
  });

  context.subscriptions.push(checkCmd, convertCmd, runCmd);
}

function deactivate() {}

module.exports = {
  activate,
  deactivate
};
