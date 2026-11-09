# Ekstensi Pseudocode (Telkom University / pseugo)

Ekstensi resmi Visual Studio Code untuk bahasa Pseudocode pada mata kuliah Algoritma dan Pemrograman (Alpro) Telkom University.

Menyediakan pengalaman coding modern setara **Pylance / gopls** dengan strict typechecking real-time, autocompletion cerdas, hover documentation, dan navigasi definisi.

---

## Fitur Utama

### 1. Real-time Strict Typechecking (Seperti Pylance / gopls)
- Validasi semantik dan sintaksis berjalan otomatis secara **real-time saat mengetik** (*as-you-type*) dan saat menyimpan file (*on save*).
- Menampilkan garis bergelombang merah (*red squiggly lines*) langsung di bawah token yang bermasalah.
- Semua error dicatat di panel **Problems** (`Ctrl+Shift+M`) VS Code dengan lokasi baris dan kolom yang presisi:
  - Variabel belum dideklarasikan (*undeclared variable*)
  - Ketidakcocokan tipe penugasan (*assignment type mismatch*)
  - Operasi tidak valid (misal `DIV` pada tipe `real`)
  - Kondisi percabangan/perulangan non-boolean
  - Duplikasi nama variabel atau tipe record
  - Akses field record yang tidak terdaftar

### 2. Hover Tooltips & Type Inspection
- Mengarahkan kursor (*hover*) pada variabel atau konstanta akan menampilkan tipe datanya (misal `(variable) total : integer`, `(constant) PI : real = 3.14`).
- Mengarahkan kursor pada variabel alias menampilkan target variabel aslinya.
- Mengarahkan kursor pada tipe record menampilkan struktur field lengkapnya.
- Mengarahkan kursor pada operator/keyword (`DIV`, `MOD`, `depend on`, `repeat`, `for`, `int()`, dll) menampilkan dokumentasi dan contoh penggunaannya.

### 3. Autocompletion / IntelliSense Cerdas
- Menyarankan kata kunci pseudocode, tipe data primitif, dan fungsi bawaan (`input`, `output`, `int`, `real`).
- Secara otomatis menyarankan variabel dan konstanta yang telah dideklarasikan di bagian `kamus`.
- **Dot-completion untuk Record**: Mengetik titik setelah nama variabel record (misal `mhs.`) akan secara spesifik menyarankan field-field milik tipe record tersebut.

### 4. Navigasi Go to Definition (`F12` / `Ctrl+Click`)
- Menekan `F12` atau `Ctrl+Click` pada variabel atau tipe di bagian `algoritma` akan langsung melompat ke baris deklarasinya di bagian `kamus`.

### 5. Document Outline & Breadcrumbs
- Struktur program (`program`), tipe record (`type`), konstanta, dan variabel tampil rapi di panel Outline dan breadcrumbs navigasi VS Code.

### 6. Syntax Highlighting Lengkap
- Penyorotan kata kunci struktur (`program`, `kamus`, `algoritma`, `endprogram`).
- Percabangan (`if..then..endif`, `depend on..end`, `default`).
- Perulangan (`for..to/downto..do..endfor`, `while..endwhile`, `repeat..until`).
- Tipe data primitif dan nama record.
- Variabel & alias variabel (**Light Blue**), konstanta (**Cyan**), fungsi bawaan (**Kuning**).
- Komentar 3 gaya: Pascal `{ ... }`, baris tunggal `//`, dan blok `/* ... */`.
- Operator: `<-`, `=`, `++`, `--`, `DIV`, `MOD`, `AND`, `OR`, `XOR`, `NOT`, `&&`, `||`, `!`, dsb.

### 7. Integrasi Perintah CLI `pseugo`
- Klik kanan pada editor atau buka Command Palette (`Ctrl+Shift+P`):
  - **Pseugo: Check Pseudocode (Typecheck)**
  - **Pseugo: Convert to Golang**
  - **Pseugo: Convert & Run with Go**

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

Setelah disalin, restart atau tekan `Ctrl+Shift+P` -> **`Developer: Reload Window`** pada VS Code.
