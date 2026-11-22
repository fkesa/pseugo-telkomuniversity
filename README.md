# Pseudocode Compiler & Go Transpiler

Compiler dan transpiler dua arah untuk subset pseudocode pada materi perkuliahan Algoritma dan Pemrograman (Alpro) Telkom University.

Alur pemrosesan:
- **Pseudocode ke Go**: `pseudocode source -> lexer -> parser -> AST -> semantic/type checker -> Go code generator`
- **Go ke Pseudocode**: `Go source -> Go parser/AST -> decompiler/transpiler -> pseudocode source`

Untuk dokumentasi lengkap logika pseudocode, tipe data, operator, dan padanan sintaks, silakan baca [docs/pseudocode_logic.md](docs/pseudocode_logic.md).

---

## Persyaratan

- Go 1.22 atau lebih baru
- Tanpa dependensi eksternal (menggunakan Go standard library murni)

---

## Cara Menjalankan (`pseugo`)

Sebelum menjalankan perintah, kompilasi binary `pseugo` terlebih dahulu di root repositori:
```bash
# Build binary di root repositori
go build -o pseugo ./cmd/pseugo       # Linux / macOS
go build -o pseugo.exe ./cmd/pseugo   # Windows (PowerShell: .\pseugo)
```

### 1. Jalankan Pseudocode Langsung
```bash
# Compile sementara dan langsung eksekusi pseudocode (mendukung input interaktif)
./pseugo run examples/div_mod.pseudo
```

### 2. Konversi Pseudocode ke Golang
```bash
# Simpan hasil Go ke file (posisi: <input> [output])
./pseugo convert -to-golang examples/div_mod.pseudo build/div_mod.go

# Atau langsung tanpa flag mode (otomatis terdeteksi dari ekstensi .pseudo)
./pseugo convert examples/div_mod.pseudo build/div_mod.go
./pseugo examples/div_mod.pseudo build/div_mod.go

# Atau cetak langsung ke stdout (tanpa argumen output)
./pseugo convert -to-golang examples/div_mod.pseudo

# Jalankan kode Go yang dihasilkan
go run ./build/div_mod.go
```

### 3. Konversi Golang ke Pseudocode
```bash
# Simpan hasil pseudocode ke file (posisi: <input> [output])
./pseugo convert -to-pseu build/div_mod.go div_mod.pseudo

# Atau langsung tanpa flag mode (otomatis terdeteksi dari ekstensi .go)
./pseugo convert build/div_mod.go div_mod.pseudo
./pseugo build/div_mod.go div_mod.pseudo

# Atau cetak langsung ke stdout (tanpa argumen output)
./pseugo convert -to-pseu build/div_mod.go
```

### 4. Typechecker Pseudocode (Validasi Semantik Tanpa Kompilasi)
```bash
./pseugo check examples/div_mod.pseudo
```

### 5. Tampilkan Token Lexer
```bash
./pseugo tokens examples/div_mod.pseudo
```

---

## Sintaks & Logika yang Didukung

```text
program NamaProgram
kamus
  nilai, total : integer
  a, b : real
  constant batas : integer <- 10
  type waktu <
    jam : integer
    menit : integer
  >
  wt : waktu
algoritma
  nilai <- 2 + 3 * 4
  input(a)
  output("hasil:", nilai)

  // Percabangan
  if nilai >= batas then
    output("Di atas batas")
  elseif nilai = 0 then
    output("Nol")
  else
    output("Di bawah batas")
  endif

  // Percabangan Seleksi Nilai (Switch-Case)
  depend on nilai
  10:
    output("Sepuluh")
  20, 30:
    output("Dua puluh atau tiga puluh")
  default:
    output("Lainnya")
  end

  // Perulangan For (Maju dan Mundur)
  for i <- 1 to 5 do
    nilai <- nilai + i
  endfor

  // Perulangan While
  while nilai > 0 do
    nilai <- nilai - 1
  endwhile

  // Perulangan Repeat-Until
  repeat
    nilai <- nilai + 2
  until nilai >= 10
endprogram
```

### Rangkuman Logika & Fitur:
- **Komentar**: Pascal `{ ... }`, baris tunggal `//`, dan blok `/* ... */`.
- **Tipe data dasar**: `integer`/`int`, `real`/`float`, `boolean`/`bool`, `string`, dan `char`/`rune`.
- **Tipe bentukan (Record)**: `type Nama < ... >` dan akses field record (`wt.jam`).
- **Konstanta**: `constant nama : tipe <- nilai` atau `constant nama : tipe = nilai`.
- **Alias variabel**: `alias : original` (merujuk ke alokasi memori yang sama).
- **Assignment**: `<-` dan `=`.
- **Postfix Increment / Decrement**: `x++` dan `x--`.
- **Input / Output**: `input(...)` dan `output(...)`.
- **Percabangan**:
  - `if ... then ... elseif ... then ... else ... endif`
  - `depend on <ekspresi> ... end` dan `depend on ... end` tanpa ekspresi (seleksi banyak nilai / kondisi boolean, ekuivalen `switch-case`; label case sejajar dengan `depend on`, aksi 1 indent, ditutup `end`)
- **Perulangan**:
  - `for ... to ... do ... endfor` (inklusif maju)
  - `for ... downto ... do ... endfor` (inklusif mundur)
  - `while ... do ... endwhile`
  - `repeat ... until ...`
- **Konversi Tipe (Type Casting)**: `int(pi)`, `real(n)`, `integer(...)`, `float(...)`.
- **Operator aritmetika**: `+`, `-`, `*`, `/` (pembagian riil), `DIV` (pembagian bulat), `MOD` (modulo integer).
- **Operator logika**: `NOT`, `AND` (`&&`), `OR` (`||`), `XOR`.
- **Operator perbandingan**: `=`, `==`, `/=`, `!=`, `<>`, `<`, `<=`, `>`, `>=`.
- **Dukungan multi-program & pemanggilan subprogram** dalam satu file (misal `jeniskendaraan()` di dalam perulangan `for`).

---

## Struktur Proyek

```text
cmd/pseugo/              CLI compiler & decompiler
docs/                    Dokumentasi lengkap logika & sintaks pseudocode
internal/cli/            Implementasi perintah & argumen CLI
internal/compiler/       Lexer, parser, AST, type checker, code generator, decompiler
test/                    Seluruh file pengujian (unit, integration, dan round-trip tests)
examples/                Contoh source pseudocode & tugas praktikum
vscode-extension/        Ekstensi resmi VS Code untuk pseudocode (syntax highlight, snippets, commands)
```

## Repositori Materi Rujukan

- https://github.com/fkesa/algoritma-pemrograman
- https://github.com/fkesa/algoritma-pemrograman-praktikum

# For Development
```bash
git clone https://github.com/fkesa/pseugo-telkomuniversity
cd pseugo-telkomuniversity && go build -trimpath -ldflags="-s -w" -o pseugo ./cmd/pseugo
```
