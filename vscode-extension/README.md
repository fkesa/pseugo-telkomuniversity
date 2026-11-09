# Ekstensi Pseudocode (Telkom University / pseugo)

Ekstensi resmi Visual Studio Code untuk bahasa Pseudocode pada mata kuliah Algoritma dan Pemrograman (Alpro) Telkom University.

---

## Fitur Utama

1. **Syntax Highlighting Lengkap**:
   - Struktur program: `program`, `endprogram`, `kamus`, `algoritma`
   - Percabangan: `if`, `then`, `else`, `elseif`, `else if`, `endif`, `depend on`, `end`, `default`
   - Perulangan: `for ... to ... do ... endfor`, `for ... downto ... do ... endfor`, `while ... do ... endwhile`, `repeat ... until ...`
   - Tipe data: `integer`, `int`, `real`, `float`, `boolean`, `bool`, `string`, `char`, `rune`, serta tipe record `type < ... >`
   - Komentar: Pascal `{ ... }`, baris `//`, dan blok `/* ... */`
   - Operator: `<-`, `=`, `++`, `--`, `DIV`, `MOD`, `AND`, `OR`, `XOR`, `NOT`, `&&`, `||`, `!`, dsb.
2. **Snippets Cepat**:
   - Ketik `prog` -> Membuat kerangka lengkap program
   - Ketik `if`, `ifelse`, `ifelseif` -> Template percabangan
   - Ketik `depend` -> Template `depend on ... end`
   - Ketik `for`, `fordown` -> Template perulangan for maju/mundur
   - Ketik `while` -> Template while-do
   - Ketik `repeat` -> Template repeat-until
   - Ketik `type` -> Template record/struct
3. **Integrasi Eksekusi & Validasi `pseugo`**:
   - Klik kanan pada editor file `.pseudo` atau buka Command Palette (`Ctrl+Shift+P`):
     - **Pseugo: Check Pseudocode (Typecheck)**: Memeriksa kebenaran sintaks dan tipe data secara instan.
     - **Pseugo: Convert to Golang**: Mengonversi file pseudocode menjadi file `.go` dan membukanya berdampingan.
     - **Pseugo: Convert & Run with Go**: Mengonversi ke Go lalu langsung menjalankannya di terminal terintegrasi VS Code.

---

## File yang Didukung

- `.pseudo`
- `.pseudocode`

---

## Cara Instalasi

Ekstensi ini dapat dipasang secara langsung dengan menyalin folder `vscode-extension` ke folder ekstensi VS Code Anda:

```bash
# Windows
xcopy /E /I /Y "vscode-extension" "%USERPROFILE%\.vscode\extensions\pseugo-pseudocode"

# Linux / macOS
cp -r vscode-extension ~/.vscode/extensions/pseugo-pseudocode
```

Setelah disalin, restart atau buka ulang Visual Studio Code.
