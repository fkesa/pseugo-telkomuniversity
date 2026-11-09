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

### 1. Konversi Pseudocode ke Golang
```bash
# Simpan hasil Go ke file
go run ./cmd/pseugo convert -to-golang -i examples/div_mod.pseudo -o build/div_mod.go

# Atau cetak langsung ke stdout (tanpa -o)
go run ./cmd/pseugo convert -to-golang -i examples/div_mod.pseudo

# Jalankan kode Go yang dihasilkan
go run ./build/div_mod.go
```

### 2. Konversi Golang ke Pseudocode
```bash
# Simpan hasil pseudocode ke file
go run ./cmd/pseugo convert -to-pseu -i build/div_mod.go -o div_mod.pseudo

# Atau cetak langsung ke stdout (tanpa -o)
go run ./cmd/pseugo convert -to-pseu -i build/div_mod.go
```

### 3. Typechecker Pseudocode (Validasi Semantik Tanpa Kompilasi)
```bash
go run ./cmd/pseugo check examples/div_mod.pseudo
```

### 4. Tampilkan Token Lexer
```bash
go run ./cmd/pseugo tokens examples/div_mod.pseudo
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
- **Tipe data dasar**: `integer`/`int`, `real`/`float`, `boolean`/`bool`, `string`, dan `char`/`rune`.
- **Tipe bentukan (Record)**: `type Nama < ... >` dan akses field record (`wt.jam`).
- **Konstanta**: `constant nama : tipe <- nilai`.
- **Alias variabel**: `alias : original` (merujuk ke alokasi memori yang sama).
- **Assignment**: `<-`.
- **Input / Output**: `input(...)` dan `output(...)`.
- **Percabangan**: `if ... then ... elseif ... then ... else ... endif`.
- **Perulangan**:
  - `for ... to ... do ... endfor` (inklusif maju)
  - `for ... downto ... do ... endfor` (inklusif mundur)
  - `while ... do ... endwhile`
  - `repeat ... until ...`
- **Operator aritmetika**: `+`, `-`, `*`, `/` (pembagian riil), `DIV` (pembagian bulat), `MOD` (modulo integer).
- **Operator logika**: `NOT`, `AND`, `OR`, `XOR`.
- **Operator perbandingan**: `=`, `==`, `/=`, `!=`, `<>`, `<`, `<=`, `>`, `>=`.
- **Dukungan multi-program** dalam satu file.

---

## Struktur Proyek

```text
cmd/pseugo/              CLI compiler & decompiler
docs/                    Dokumentasi lengkap logika & sintaks pseudocode
internal/cli/            Implementasi perintah & argumen CLI
internal/compiler/       Lexer, parser, AST, type checker, code generator, decompiler
test/                    Seluruh file pengujian (unit, integration, dan round-trip tests)
examples/                Contoh source pseudocode & tugas praktikum
```

## Repositori Materi Rujukan

- https://github.com/fkesa/algoritma-pemrograman
- https://github.com/fkesa/algoritma-pemrograman-praktikum

# For Development
```bash
git clone https://github.com/fkesa/pseugo-telkomuniversity
cd pseugo-telkomuniversity && go build -o .\build\pseugo.exe .\cmd\pseugo
```
