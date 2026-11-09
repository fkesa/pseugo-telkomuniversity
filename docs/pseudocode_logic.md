# Dokumentasi Sintaks & Logika Pseudocode

Dokumen ini berisi spesifikasi lengkap mengenai sintaks, logika, tipe data, operator, serta panduan fitur dua arah (Pseudocode ke Go dan Go ke Pseudocode) pada compiler pembelajaran Algoritma dan Pemrograman (Alpro) Telkom University.

---

## Daftar Isi
1. [Struktur Dasar Program](#1-struktur-dasar-program)
2. [Tipe Data](#2-tipe-data)
3. [Deklarasi (`kamus`)](#3-deklarasi-kamus)
4. [Logika & Kontrol Alur (`algoritma`)](#4-logika--kontrol-alur-algoritma)
   - [Assignment (Penugasan)](#assignment-penugasan)
   - [Input dan Output](#input-dan-output)
   - [Percabangan (If - Then - Elseif - Else - Endif)](#percabangan-if---then---elseif---else---endif)
   - [Perulangan (For, While, Repeat-Until)](#perulangan-looping)
5. [Operator yang Didukung](#5-operator-yang-didukung)
   - [Aritmetika](#operator-aritmetika)
   - [Perbandingan / Relasional](#operator-perbandingan--relasional)
   - [Logika](#operator-logika)
6. [Fitur Konversi Golang ke Pseudocode](#6-fitur-konversi-golang-ke-pseudocode)
   - [Pemetaan Sintaks Go ke Pseudocode](#pemetaan-sintaks-go-ke-pseudocode)
   - [Rekonstruksi Otomatis Kamus](#rekonstruksi-otomatis-kamus)
7. [Panduan Penggunaan CLI (`pseugo`)](#7-panduan-penggunaan-cli-pseugo)
8. [Contoh Kasus Komprehensif](#8-contoh-kasus-komprehensif)

---

## 1. Struktur Dasar Program

Setiap program pseudocode memiliki kerangka standar yang terdiri dari 3 bagian: nama program, bagian deklarasi (`kamus`), dan bagian eksekusi logika (`algoritma`).

```text
program NamaProgram
kamus
    // deklarasi variabel, konstanta, dan tipe data bentukan
algoritma
    // instruksi, input/output, percabangan, dan perulangan
endprogram
```

### Multi-Program dalam Satu File
Compiler mendukung penulisan beberapa program dalam satu file `.pseudo`. Pada saat dikompilasi ke Go, setiap program akan dibuat sebagai fungsi terpisah dan dipanggil berurutan di dalam fungsi `main()`.

---

## 2. Tipe Data

Pseudocode mendukung tipe data primitif dan tipe data bentukan (record):

| Tipe Pseudocode | Padanan di Go | Keterangan & Contoh Literal |
| :--- | :--- | :--- |
| `integer` atau `int` | `int` | Bilangan bulat. Contoh: `0`, `42`, `-10` |
| `real` atau `float` | `float64` | Bilangan riil/desimal. Contoh: `3.14`, `0.0`, `-2.5` |
| `boolean` atau `bool` | `bool` | Nilai kebenaran: `true`/`false` atau `benar`/`salah` |
| `string` | `string` | Untaian teks berkutip ganda. Contoh: `"Halo Dunia"` |
| `char` atau `rune` | `rune` | Karakter tunggal berkutip tunggal. Contoh: `'A'`, `'9'`, `'\n'` |
| Tipe Bentukan (`type`) | `struct` | Record yang menggabungkan beberapa field |

---

## 3. Deklarasi (`kamus`)

Semua variabel, konstanta, tipe record, dan alias dideklarasikan di dalam blok `kamus` sebelum bagian `algoritma`.

### A. Deklarasi Variabel
Dapat dideklarasikan satu per satu atau sekaligus dipisahkan koma:
```text
kamus
    x : integer
    a, b, c : real
    nama : string
    sudah_selesai : boolean
```

### B. Deklarasi Konstanta
Konstanta bernilai tetap dan nilainya tidak dapat diubah di bagian `algoritma`.
```text
kamus
    constant pi : real <- 3.14
    constant batas_maks : integer <- 100
    constant status_aktif : boolean <- true
```

### C. Deklarasi Tipe Bentukan (Record / Struct)
Record dideklarasikan menggunakan blok `type <Nama> < ... >`:
```text
kamus
    type Titik <
        x : real
        y : real
    >
    type Mahasiswa <
        nim : string
        nilai : integer
        lulus : boolean
    >
    p1 : Titik
    mhs : Mahasiswa
```

### D. Deklarasi Alias Variabel
Alias membuat nama baru yang menunjuk ke lokasi memori variabel yang sudah dideklarasikan sebelumnya:
```text
kamus
    original : integer
    alias : original
```

---

## 4. Logika & Kontrol Alur (`algoritma`)

### Assignment (Penugasan)
Penugasan menggunakan operator panah kiri `<-`:
```text
algoritma
    x <- 10
    total <- total + x
    p1.x <- 5.0
    mhs.lulus <- true
```

### Input dan Output
- `input(...)`: Membaca masukan dari pengguna (dipisahkan whitespace). Untuk `char`, membaca karakter pertama sebagai satu rune.
- `output(...)`: Menampilkan nilai, ekspresi, atau teks ke layar.

```text
algoritma
    input(x, y)
    input(p1.x, p1.y)
    output("Hasil penjumlahan:", x + y)
```

### Percabangan (If - Then - Elseif - Else - Endif)

Mendukung 3 pola percabangan:

#### 1. Percabangan Tunggal (`if - then - endif`)
```text
if nilai >= 75 then
    output("Lulus")
endif
```

#### 2. Percabangan Ganda (`if - then - else - endif`)
```text
if x MOD 2 = 0 then
    output(x, "adalah bilangan genap")
else
    output(x, "adalah bilangan ganjil")
endif
```

#### 3. Percabangan Majemuk (`if - then - elseif - else - endif`)
```text
if skor >= 80 then
    indeks <- "A"
elseif skor >= 70 then
    indeks <- "B"
elseif skor >= 60 then
    indeks <- "C"
else
    indeks <- "E"
endif
```

### Perulangan (Looping)

Mendukung 4 variasi perulangan:

#### 1. For Maju (`to`)
Perulangan dengan pencacah menaik secara inklusif (dari batas awal hingga batas akhir):
```text
for i <- 1 to 10 do
    output("Iterasi ke-", i)
endfor
```

#### 2. For Mundur (`downto`)
Perulangan dengan pencacah menurun secara inklusif:
```text
for i <- 10 downto 1 do
    output("Hitung mundur:", i)
endfor
```

#### 3. While-Do
Perulangan dengan pemeriksaan kondisi di awal. Blok perulangan hanya berjalan selama kondisi bernilai `true`:
```text
while nilai < 100 do
    nilai <- nilai * 2
    output(nilai)
endwhile
```

#### 4. Repeat-Until
Perulangan dengan pemeriksaan kondisi di akhir. Blok perulangan dijalankan minimal satu kali, lalu berhenti saat kondisi bernilai `true`:
```text
repeat
    input(x)
    total <- total + x
until x = 0
```

---

## 5. Operator yang Didukung

### Operator Aritmetika
| Operator | Fungsi | Tipe Operand | Tipe Hasil |
| :---: | :--- | :--- | :--- |
| `+` | Penjumlahan | `integer` / `real` | Sesuai operand (promosi ke `real` jika salah satu `real`) |
| `-` | Pengurangan / Negasi | `integer` / `real` | Sesuai operand |
| `*` | Perkalian | `integer` / `real` | Sesuai operand |
| `/` | Pembagian Riil | `integer` / `real` | `real` (selalu menghasilkan riil) |
| `DIV` | Pembagian Bulat | `integer` saja | `integer` |
| `MOD` | Sisa Bagi (Modulo) | `integer` saja | `integer` |

### Operator Perbandingan / Relasional
| Operator | Fungsi |
| :---: | :--- |
| `=` atau `==` | Sama dengan |
| `!=` atau `/=` atau `<>` | Tidak sama dengan |
| `<` | Kurang dari |
| `<=` | Kurang dari atau sama dengan |
| `>` | Lebih dari |
| `>=` | Lebih dari atau sama dengan |

> Hasil operasi perbandingan selalu bertipe `boolean`.

### Operator Logika
| Operator | Fungsi | Operand | Hasil |
| :---: | :--- | :--- | :--- |
| `AND` | Konjungsi logika (DAN) | `boolean` | `boolean` |
| `OR` | Disjungsi logika (ATAU) | `boolean` | `boolean` |
| `XOR` | Exclusive OR | `boolean` | `boolean` |
| `NOT` | Negasi logika (BUKAN) | `boolean` | `boolean` |

---

## 6. Fitur Konversi Golang ke Pseudocode

Compiler menyediakan decompiler Go -> Pseudocode yang mengurai AST Go (`go/parser` dan `go/ast`) dan merekonstruksinya menjadi pseudocode yang valid.

### Pemetaan Sintaks Go ke Pseudocode

| Kode Sumber Go | Hasil Pseudocode |
| :--- | :--- |
| `var x int` | `x : integer` (di `kamus`) |
| `var a, b float64` | `a, b : real` (di `kamus`) |
| `const PI = 3.14` | `constant PI : real <- 3.14` |
| `type Prisma struct { ... }` | `type Prisma < ... >` |
| `x := 10` | `x : integer` di `kamus`, `x <- 10` di `algoritma` |
| `x = 5` | `x <- 5` |
| `x += 2` | `x <- x + 2` |
| `x++` / `x--` | `x <- x + 1` / `x <- x - 1` |
| `fmt.Scan(&a, &b)` | `input(a, b)` |
| `fmt.Println("Total:", sum)` | `output("Total:", sum)` |
| `if cond { ... } else { ... }` | `if cond then ... else ... endif` |
| `for i := 1; i <= n; i++` | `for i <- 1 to n do ... endfor` |
| `for i := n; i >= 1; i--` | `for i <- n downto 1 do ... endfor` |
| `for cond { ... }` | `while cond do ... endwhile` |
| `for { ... if cond { break } }` | `repeat ... until cond` |
| `a / b` (integer) | `a DIV b` |
| `a / b` (float/real) | `a / b` |
| `a % b` | `a MOD b` |
| `&&`, `\|\|`, `!` | `AND`, `OR`, `NOT` |
| `==`, `!=` | `=`, `!=` |

### Rekonstruksi Otomatis Kamus
Pada kode Go, deklarasi variabel sering kali tersebar (misalnya dengan operator short-declaration `:=` atau `var` di tengah fungsi). Decompiler secara otomatis:
1. Menemukan semua variabel dan tipe data yang digunakan.
2. Melakukan inferensi tipe data (`integer`, `real`, `boolean`, `string`, `char`, atau nama `type`).
3. Mengelompokkan variabel berdasarkan tipe datanya di blok `kamus`.
4. Menempatkan inisialisasi nilainya di blok `algoritma` sebagai penugasan (`<-`).
5. Membersihkan prefix internal (seperti `__pseudo_user_`).

---

## 7. Panduan Penggunaan CLI (`pseugo`)

CLI menyediakan mode sub-perintah (`convert`, `check`, dan `tokens`):

### Konversi Pseudocode ke Go
```bash
# Simpan hasil Go ke file
go run ./cmd/pseugo convert -to-golang -i examples/div_mod.pseudo -o build/div_mod.go

# Cetak hasil Go ke terminal (stdout)
go run ./cmd/pseugo convert -to-golang -i examples/div_mod.pseudo

# Jalankan kode Go yang dihasilkan
go run ./build/div_mod.go
```

### Konversi Go ke Pseudocode
```bash
# Simpan hasil pseudocode ke file
go run ./cmd/pseugo convert -to-pseu -i build/div_mod.go -o div_mod.pseudo

# Cetak hasil pseudocode ke terminal (stdout)
go run ./cmd/pseugo convert -to-pseu -i build/div_mod.go
```

### Menjalankan Typechecker Saja (Validasi Tanpa Kompilasi)
Memeriksa apakah sintaks dan tipe data pseudocode valid sebelum dikonversi:

```bash
go run ./cmd/pseugo check examples/div_mod.pseudo
```

### Melihat Token Lexer
```bash
go run ./cmd/pseugo tokens examples/div_mod.pseudo
```

---

## 8. Contoh Kasus Komprehensif

### Contoh 1: Pembagian Bulat dan Modulo
```text
program PembagianDanModulo
kamus
    bil1, bil2, hasil_bagi, sisa_bagi : integer
algoritma
    input(bil1, bil2)
    hasil_bagi <- bil1 DIV bil2
    sisa_bagi <- bil1 MOD bil2
    output("Hasil Bagi:", hasil_bagi)
    output("Sisa Bagi:", sisa_bagi)
endprogram
```

### Contoh 2: Record / Struct Geometri
```text
program HitungPrisma
kamus
    type prisma <
        panjang : real
        lebar : real
        tinggi : real
        volume : real
    >
    pr : prisma
algoritma
    input(pr.panjang, pr.lebar, pr.tinggi)
    pr.volume <- pr.panjang * pr.lebar * pr.tinggi
    output("Volume:", pr.volume)
endprogram
```

### Contoh 3: Percabangan Majemuk dan Perulangan
```text
program StatistikAngka
kamus
    i, n, bil, total : integer
    rata_rata : real
algoritma
    input(n)
    total <- 0
    for i <- 1 to n do
        input(bil)
        if bil > 0 then
            total <- total + bil
        endif
    endfor
    rata_rata <- total / n
    output("Total Positif:", total)
    output("Rata-rata:", rata_rata)
endprogram
```
