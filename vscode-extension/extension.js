const vscode = require('vscode');
const { execFile } = require('child_process');
const path = require('path');
const fs = require('fs');
const os = require('os');

let diagnosticCollection;
let outputChannel;

function getOutputChannel() {
  if (!outputChannel) {
    outputChannel = vscode.window.createOutputChannel('Pseugo Pseudocode');
  }
  return outputChannel;
}

function findPseugoBinary(activeDocPath) {
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

function escapeRegExp(string) {
  return string.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}

/**
 * Find exact start and end column of a word inside a line
 */
function findWordInLine(rawLine, word, startIndex = 0, endIndex = rawLine.length) {
  const searchArea = rawLine.substring(startIndex, endIndex);
  const regex = new RegExp(`\\b${escapeRegExp(word)}\\b`);
  const match = regex.exec(searchArea);
  if (match) {
    const startCol = startIndex + match.index;
    return {
      startCol,
      endCol: startCol + word.length
    };
  }
  const idx = searchArea.indexOf(word);
  if (idx !== -1) {
    const startCol = startIndex + idx;
    return {
      startCol,
      endCol: startCol + word.length
    };
  }
  return { startCol: 0, endCol: word.length };
}

/**
 * Run strict real-time diagnostics on a pseudocode document
 */
function runDiagnostics(document) {
  if (!diagnosticCollection || document.languageId !== 'pseudocode') {
    return;
  }

  const bin = findPseugoBinary(document.fileName);
  const text = document.getText();

  // Write unsaved buffer to temporary file for live as-you-type validation
  const tmpDir = os.tmpdir();
  const tmpFile = path.join(tmpDir, `pseugo_lint_${Date.now()}_${process.pid}.pseudo`);

  fs.writeFile(tmpFile, text, 'utf8', (writeErr) => {
    if (writeErr) return;

    execFile(bin, ['check', tmpFile], (err, stdout, stderr) => {
      // Clean up temp file
      fs.unlink(tmpFile, () => {});

      if (!err) {
        // No errors! Clear diagnostics
        diagnosticCollection.delete(document.uri);
        return;
      }

      const output = (stderr || '') + '\n' + (stdout || '');
      const diagnostics = [];

      // Regex matches errors like: "29:14: expected end of line" or "13:5: undeclared variable \"x\""
      const errRegex = /(?:^|\n)(\d+):(\d+):\s*(.+)/g;
      let match;
      while ((match = errRegex.exec(output)) !== null) {
        const line = Math.max(0, parseInt(match[1], 10) - 1);
        const col = Math.max(0, parseInt(match[2], 10) - 1);
        const msg = match[3].trim();

        let range;
        if (line < document.lineCount) {
          const lineText = document.lineAt(line).text;
          const wordRange = document.getWordRangeAtPosition(new vscode.Position(line, col));
          if (wordRange) {
            range = wordRange;
          } else {
            const startCol = Math.min(col, lineText.length);
            const endCol = Math.min(startCol + 1, lineText.length);
            range = new vscode.Range(line, startCol, line, Math.max(endCol, startCol + 1));
          }
        } else {
          range = new vscode.Range(0, 0, 0, 1);
        }

        const diag = new vscode.Diagnostic(range, msg, vscode.DiagnosticSeverity.Error);
        diag.source = 'pseugo (strict)';
        diagnostics.push(diag);
      }

      diagnosticCollection.set(document.uri, diagnostics);
    });
  });
}

/**
 * Parse document symbols (types, vars, constants, programs) with exact column ranges
 */
function parseDocumentSymbols(document) {
  const text = document.getText();
  const lines = text.split(/\r?\n/);

  const programs = [];
  const types = {}; // name -> { line, startCol, endCol, fields: { name: { type, line, startCol, endCol } } }
  const vars = {};  // name -> { type, line, startCol, endCol, aliasOf }
  const constants = {}; // name -> { type, val, line, startCol, endCol }
  const aliases = {}; // name -> { of, line, startCol, endCol }

  let inKamus = false;
  let inType = null;

  for (let i = 0; i < lines.length; i++) {
    const rawLine = lines[i];
    // Strip comments
    const line = rawLine.replace(/\{[^}]*\}|\/\/.*$/g, '').trim();
    if (!line) continue;

    const progMatch = line.match(/^program\s+([a-zA-Z_][a-zA-Z0-9_]*)/i);
    if (progMatch) {
      const progName = progMatch[1];
      const pos = findWordInLine(rawLine, progName, 0);
      programs.push({ name: progName, line: i, startCol: pos.startCol, endCol: pos.endCol });
      continue;
    }

    if (/^kamus\b/i.test(line)) {
      inKamus = true;
      continue;
    }
    if (/^algoritma\b/i.test(line) || /^endprogram\b/i.test(line)) {
      inKamus = false;
      inType = null;
      continue;
    }

    if (inKamus) {
      // Check type definition start: type NamaRecord <
      const typeStartMatch = line.match(/^type\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*</i);
      if (typeStartMatch) {
        inType = typeStartMatch[1];
        const pos = findWordInLine(rawLine, inType, 0);
        types[inType] = { line: i, startCol: pos.startCol, endCol: pos.endCol, fields: {} };
        continue;
      }
      if (inType) {
        if (line.includes('>')) {
          inType = null;
          continue;
        }
        // Field definition: f1, f2 : integer
        const colonIdx = rawLine.indexOf(':');
        const fMatch = line.match(/^([a-zA-Z0-9_,\s]+):\s*([a-zA-Z_][a-zA-Z0-9_]*)/i);
        if (fMatch && types[inType] && colonIdx !== -1) {
          const fNames = fMatch[1].split(',').map(s => s.trim()).filter(Boolean);
          const fType = fMatch[2].trim();
          for (const fn of fNames) {
            const pos = findWordInLine(rawLine, fn, 0, colonIdx);
            types[inType].fields[fn] = {
              type: fType,
              line: i,
              startCol: pos.startCol,
              endCol: pos.endCol
            };
          }
        }
        continue;
      }

      // Constant definition: constant MAX : integer <- 100
      const colonIdx = rawLine.indexOf(':');
      const constMatch = line.match(/^constant\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*:\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:<-|=)\s*(.+)/i);
      if (constMatch) {
        const cName = constMatch[1];
        const pos = findWordInLine(rawLine, cName, 0, colonIdx !== -1 ? colonIdx : rawLine.length);
        constants[cName] = {
          type: constMatch[2],
          val: constMatch[3].trim(),
          line: i,
          startCol: pos.startCol,
          endCol: pos.endCol
        };
        continue;
      }

      // Var or Alias definition: v1, v2 : type / alias
      const varMatch = line.match(/^([a-zA-Z0-9_,\s]+):\s*([a-zA-Z_][a-zA-Z0-9_]*)/i);
      if (varMatch && colonIdx !== -1) {
        const vNames = varMatch[1].split(',').map(s => s.trim()).filter(Boolean);
        const targetType = varMatch[2].trim();

        if (vars[targetType]) {
          // It's an alias of an existing var!
          for (const vn of vNames) {
            const pos = findWordInLine(rawLine, vn, 0, colonIdx);
            aliases[vn] = {
              of: targetType,
              line: i,
              startCol: pos.startCol,
              endCol: pos.endCol
            };
            vars[vn] = {
              type: vars[targetType].type,
              line: i,
              startCol: pos.startCol,
              endCol: pos.endCol,
              aliasOf: targetType
            };
          }
        } else {
          for (const vn of vNames) {
            const pos = findWordInLine(rawLine, vn, 0, colonIdx);
            vars[vn] = {
              type: targetType,
              line: i,
              startCol: pos.startCol,
              endCol: pos.endCol
            };
          }
        }
        continue;
      }
    }
  }

  return { programs, types, vars, constants, aliases };
}

/**
 * Keyword Documentation for Hover
 */
const KEYWORD_DOCS = {
  'div': '**DIV**: Pembagian bulat (integer division).\nHanya menerima operan bertipe `integer` dan menghasilkan `integer`.\nContoh: `10 DIV 3` menghasilkan `3`.',
  'mod': '**MOD**: Sisa pembagian bulat (modulo).\nHanya menerima operan bertipe `integer` dan menghasilkan `integer`.\nContoh: `10 MOD 3` menghasilkan `1`.',
  'program': '**program <Nama>**: Titik awal pendefinisian program pseudocode.\nHarus memiliki bagian `kamus` dan `algoritma`, lalu ditutup `endprogram`.',
  'kamus': '**kamus**: Blok deklarasi variabel, konstanta, tipe record, dan alias sebelum dieksekusi di `algoritma`.',
  'algoritma': '**algoritma**: Blok eksekusi instruksi logika (percabangan, perulangan, I/O, penugasan).',
  'endprogram': '**endprogram**: Kata kunci penutup program pseudocode.',
  'input': '**input(...)**: Membaca satu atau beberapa masukan dari pengguna (dipisahkan whitespace).',
  'output': '**output(...)**: Menampilkan teks, variabel, atau hasil ekspresi ke layar pengguna.',
  'depend': '**depend on <ekspresi> ... end**: Percabangan seleksi nilai (switch-case).\nLabel kasus ditulis langsung nilainya (misal `1:`, `\'A\':`) dan ditutup `end` atau `enddepend`.',
  'repeat': '**repeat ... until <kondisi>**: Perulangan dengan evaluasi di akhir.\nDijalankan minimal 1 kali, berhenti saat kondisi bernilai `true`.',
  'while': '**while <kondisi> do ... endwhile**: Perulangan dengan evaluasi di awal.\nBerjalan selama kondisi bernilai `true`.',
  'for': '**for <var> <- <awal> to|downto <akhir> do ... endfor**: Perulangan pencacah inklusif menaik (`to`) atau menurun (`downto`).',
  'constant': '**constant <nama> : <tipe> <- <nilai>**: Mendeklarasikan konstanta tetap yang nilainya tidak dapat diubah.',
  'type': '**type <NamaRecord> < ... >**: Mendeklarasikan tipe data bentukan (record/struct) yang memiliki beberapa field.',
  'int': '**int(x)** / **integer(x)**: Konversi tipe data (type casting) dari desimal (`real`) menjadi bulat (`integer`).',
  'real': '**real(x)** / **float(x)**: Konversi tipe data (type casting) dari bilangan bulat (`integer`) menjadi desimal (`real`).'
};

function activate(context) {
  diagnosticCollection = vscode.languages.createDiagnosticCollection('pseugo');
  context.subscriptions.push(diagnosticCollection);

  const channel = getOutputChannel();

  // Real-time diagnostics with debounce
  let debounceTimer = null;
  const triggerDiagnostics = (doc) => {
    if (!doc || doc.languageId !== 'pseudocode') return;
    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      runDiagnostics(doc);
    }, 300);
  };

  context.subscriptions.push(
    vscode.workspace.onDidOpenTextDocument((doc) => runDiagnostics(doc)),
    vscode.workspace.onDidSaveTextDocument((doc) => runDiagnostics(doc)),
    vscode.workspace.onDidChangeTextDocument((e) => triggerDiagnostics(e.document)),
    vscode.workspace.onDidCloseTextDocument((doc) => diagnosticCollection.delete(doc.uri))
  );

  // Run on all active pseudocode documents right away
  if (vscode.window.activeTextEditor) {
    runDiagnostics(vscode.window.activeTextEditor.document);
  }

  // 1. Hover Provider (Docstrings, Types, Variables like Pylance)
  const hoverProvider = vscode.languages.registerHoverProvider('pseudocode', {
    provideHover(document, position) {
      const range = document.getWordRangeAtPosition(position);
      if (!range) return null;

      const word = document.getText(range);
      const lower = word.toLowerCase();

      // Keyword Documentation
      if (KEYWORD_DOCS[lower]) {
        return new vscode.Hover(new vscode.MarkdownString(KEYWORD_DOCS[lower]));
      }

      // Inspect Symbol in File
      const symbols = parseDocumentSymbols(document);

      // Check Constant
      if (symbols.constants[word]) {
        const c = symbols.constants[word];
        const md = new vscode.MarkdownString();
        md.appendCodeblock(`(constant) ${word} : ${c.type} <- ${c.val}`, 'pseudocode');
        return new vscode.Hover(md);
      }

      // Check Variable
      if (symbols.vars[word]) {
        const v = symbols.vars[word];
        const md = new vscode.MarkdownString();
        if (v.aliasOf) {
          md.appendCodeblock(`(variable alias) ${word} : ${v.aliasOf} (tipe: ${v.type})`, 'pseudocode');
        } else {
          md.appendCodeblock(`(variable) ${word} : ${v.type}`, 'pseudocode');
        }
        return new vscode.Hover(md);
      }

      // Check Record Type
      if (symbols.types[word]) {
        const t = symbols.types[word];
        const md = new vscode.MarkdownString();
        let desc = `type ${word} <\n`;
        for (const [fn, fld] of Object.entries(t.fields)) {
          const ft = typeof fld === 'string' ? fld : fld.type;
          desc += `    ${fn} : ${ft}\n`;
        }
        desc += `>`;
        md.appendCodeblock(desc, 'pseudocode');
        return new vscode.Hover(md);
      }

      // Check Record Field (e.g. hovering over 'x' in 'p1.x')
      if (range.start.character > 1) {
        const prevChar = document.getText(new vscode.Range(position.line, range.start.character - 1, position.line, range.start.character));
        if (prevChar === '.') {
          const baseRange = document.getWordRangeAtPosition(new vscode.Position(position.line, range.start.character - 2));
          if (baseRange) {
            const baseVar = document.getText(baseRange);
            if (symbols.vars[baseVar]) {
              const recTypeName = symbols.vars[baseVar].type;
              if (symbols.types[recTypeName] && symbols.types[recTypeName].fields[word]) {
                const fld = symbols.types[recTypeName].fields[word];
                const fType = typeof fld === 'string' ? fld : fld.type;
                const md = new vscode.MarkdownString();
                md.appendCodeblock(`(field) ${baseVar}.${word} : ${fType}`, 'pseudocode');
                return new vscode.Hover(md);
              }
            }
          }
        }
      }

      return null;
    }
  });

  // 2. Autocompletion / IntelliSense Provider
  const completionProvider = vscode.languages.registerCompletionItemProvider(
    'pseudocode',
    {
      provideCompletionItems(document, position) {
        const linePrefix = document.lineAt(position).text.substr(0, position.character);
        const symbols = parseDocumentSymbols(document);
        const items = [];

        // Check if typing field access: var.
        const dotMatch = linePrefix.match(/([a-zA-Z_][a-zA-Z0-9_]*)\.$/);
        if (dotMatch) {
          const varName = dotMatch[1];
          if (symbols.vars[varName]) {
            const typeName = symbols.vars[varName].type;
            if (symbols.types[typeName]) {
              for (const [fn, fld] of Object.entries(symbols.types[typeName].fields)) {
                const ft = typeof fld === 'string' ? fld : fld.type;
                const ci = new vscode.CompletionItem(fn, vscode.CompletionItemKind.Field);
                ci.detail = `${fn} : ${ft}`;
                ci.documentation = `Field dari record ${typeName}`;
                items.push(ci);
              }
              return items;
            }
          }
        }

        // Global variables from kamus
        for (const [vn, v] of Object.entries(symbols.vars)) {
          const ci = new vscode.CompletionItem(vn, vscode.CompletionItemKind.Variable);
          ci.detail = `${vn} : ${v.type}`;
          ci.documentation = v.aliasOf ? `Alias dari variabel ${v.aliasOf}` : `Variabel lokal bertipe ${v.type}`;
          items.push(ci);
        }

        // Constants from kamus
        for (const [cn, c] of Object.entries(symbols.constants)) {
          const ci = new vscode.CompletionItem(cn, vscode.CompletionItemKind.Constant);
          ci.detail = `constant ${cn} : ${c.type} = ${c.val}`;
          items.push(ci);
        }

        // Record Types
        for (const [tn, t] of Object.entries(symbols.types)) {
          const ci = new vscode.CompletionItem(tn, vscode.CompletionItemKind.Struct);
          ci.detail = `type ${tn}`;
          items.push(ci);
        }

        // Standard Primitive Types
        const primTypes = ['integer', 'real', 'boolean', 'string', 'char', 'int', 'float', 'bool', 'rune'];
        for (const pt of primTypes) {
          const ci = new vscode.CompletionItem(pt, vscode.CompletionItemKind.TypeParameter);
          ci.detail = `Tipe data dasar ${pt}`;
          items.push(ci);
        }

        // Builtins & Casts
        const builtins = [
          { name: 'input', detail: 'input(var1, var2, ...)', kind: vscode.CompletionItemKind.Function },
          { name: 'output', detail: 'output(val1, val2, ...)', kind: vscode.CompletionItemKind.Function },
          { name: 'int', detail: 'int(realValue)', kind: vscode.CompletionItemKind.Function },
          { name: 'real', detail: 'real(intValue)', kind: vscode.CompletionItemKind.Function }
        ];
        for (const b of builtins) {
          const ci = new vscode.CompletionItem(b.name, b.kind);
          ci.detail = b.detail;
          items.push(ci);
        }

        // Keywords
        const keywords = [
          'program', 'kamus', 'algoritma', 'endprogram',
          'if', 'then', 'else', 'elseif', 'endif',
          'depend on', 'end', 'default',
          'for', 'to', 'downto', 'do', 'endfor',
          'while', 'endwhile', 'repeat', 'until',
          'DIV', 'MOD', 'AND', 'OR', 'NOT', 'XOR'
        ];
        for (const kw of keywords) {
          const ci = new vscode.CompletionItem(kw, vscode.CompletionItemKind.Keyword);
          items.push(ci);
        }

        return items;
      }
    },
    '.' // Trigger on dot for field completion
  );

  // 3. Go to Definition Provider (F12 / Ctrl+Click)
  const definitionProvider = vscode.languages.registerDefinitionProvider('pseudocode', {
    provideDefinition(document, position) {
      const range = document.getWordRangeAtPosition(position);
      if (!range) return null;

      const word = document.getText(range);
      const symbols = parseDocumentSymbols(document);

      // Check if clicking on field access (e.g. p1.x where cursor is on x)
      if (range.start.character > 1) {
        const prevChar = document.getText(new vscode.Range(position.line, range.start.character - 1, position.line, range.start.character));
        if (prevChar === '.') {
          const baseRange = document.getWordRangeAtPosition(new vscode.Position(position.line, range.start.character - 2));
          if (baseRange) {
            const baseVar = document.getText(baseRange);
            if (symbols.vars[baseVar]) {
              const recTypeName = symbols.vars[baseVar].type;
              if (symbols.types[recTypeName] && symbols.types[recTypeName].fields[word]) {
                const fld = symbols.types[recTypeName].fields[word];
                return new vscode.Location(
                  document.uri,
                  new vscode.Range(fld.line, fld.startCol, fld.line, fld.endCol)
                );
              }
            }
          }
        }
      }

      // Check variable / alias
      if (symbols.vars[word]) {
        const v = symbols.vars[word];
        // If this variable is an alias of a main variable:
        if (v.aliasOf && symbols.vars[v.aliasOf]) {
          const mainVar = symbols.vars[v.aliasOf];
          // Return both the alias definition and the main variable definition with exact columns!
          return [
            new vscode.Location(document.uri, new vscode.Range(v.line, v.startCol, v.line, v.endCol)),
            new vscode.Location(document.uri, new vscode.Range(mainVar.line, mainVar.startCol, mainVar.line, mainVar.endCol))
          ];
        }
        return new vscode.Location(
          document.uri,
          new vscode.Range(v.line, v.startCol, v.line, v.endCol)
        );
      }

      // Check constant
      if (symbols.constants[word]) {
        const c = symbols.constants[word];
        return new vscode.Location(
          document.uri,
          new vscode.Range(c.line, c.startCol, c.line, c.endCol)
        );
      }

      // Check record type
      if (symbols.types[word]) {
        const t = symbols.types[word];
        return new vscode.Location(
          document.uri,
          new vscode.Range(t.line, t.startCol, t.line, t.endCol)
        );
      }

      return null;
    }
  });

  // 4. Document Symbol Provider (Breadcrumbs & Outline View)
  const symbolProvider = vscode.languages.registerDocumentSymbolProvider('pseudocode', {
    provideDocumentSymbols(document) {
      const symbols = parseDocumentSymbols(document);
      const result = [];

      for (const p of symbols.programs) {
        const range = new vscode.Range(p.line, p.startCol, Math.min(document.lineCount - 1, p.line + 10), 0);
        result.push(new vscode.DocumentSymbol(
          p.name,
          'Program',
          vscode.SymbolKind.Package,
          range,
          range
        ));
      }

      for (const [tn, t] of Object.entries(symbols.types)) {
        const range = new vscode.Range(t.line, t.startCol, t.line, t.endCol);
        const sym = new vscode.DocumentSymbol(
          tn,
          'Record / Struct',
          vscode.SymbolKind.Struct,
          range,
          range
        );
        for (const [fn, fld] of Object.entries(t.fields)) {
          const ft = typeof fld === 'string' ? fld : fld.type;
          const fRange = new vscode.Range(fld.line, fld.startCol, fld.line, fld.endCol);
          sym.children.push(new vscode.DocumentSymbol(
            fn,
            ft,
            vscode.SymbolKind.Field,
            fRange,
            fRange
          ));
        }
        result.push(sym);
      }

      for (const [cn, c] of Object.entries(symbols.constants)) {
        const range = new vscode.Range(c.line, c.startCol, c.line, c.endCol);
        result.push(new vscode.DocumentSymbol(
          cn,
          `${c.type} = ${c.val}`,
          vscode.SymbolKind.Constant,
          range,
          range
        ));
      }

      for (const [vn, v] of Object.entries(symbols.vars)) {
        const range = new vscode.Range(v.line, v.startCol, v.line, v.endCol);
        result.push(new vscode.DocumentSymbol(
          vn,
          v.aliasOf ? `alias of ${v.aliasOf}` : v.type,
          vscode.SymbolKind.Variable,
          range,
          range
        ));
      }

      return result;
    }
  });

  // 5. Commands: Check, Convert, Run
  const checkCmd = vscode.commands.registerCommand('pseugo.check', () => {
    const editor = vscode.window.activeTextEditor;
    if (!editor) return;
    runDiagnostics(editor.document);
    vscode.window.showInformationMessage(`[Pseugo] Memeriksa ${path.basename(editor.document.fileName)}...`);
  });

  const convertCmd = vscode.commands.registerCommand('pseugo.convertToGolang', () => {
    const editor = vscode.window.activeTextEditor;
    if (!editor) return;

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

  const runCmd = vscode.commands.registerCommand('pseugo.runGolang', () => {
    const editor = vscode.window.activeTextEditor;
    if (!editor) return;

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
          vscode.window.showErrorMessage(`[Pseugo] Konversi gagal.`);
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

  context.subscriptions.push(
    hoverProvider,
    completionProvider,
    definitionProvider,
    symbolProvider,
    checkCmd,
    convertCmd,
    runCmd
  );
}

function deactivate() {
  if (diagnosticCollection) {
    diagnosticCollection.clear();
    diagnosticCollection.dispose();
  }
}

module.exports = {
  activate,
  deactivate
};
