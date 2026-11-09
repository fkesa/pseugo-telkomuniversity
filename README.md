# Pseudocode Compiler (Go)

Compiler pembelajaran untuk subset pseudocode pada materi Algoritma dan Pemrograman. Pipeline:

`source -> lexer -> parser -> AST -> semantic/type checker -> Go code generator`

Compiler menghasilkan source Go; compiler ini tidak mengeksekusi pseudocode secara langsung.

## Persyaratan

- Go 1.22 atau lebih baru
- Tidak memakai dependency eksternal

## Menjalankan

Dari direktori proyek:

```bash
go test ./...
go run ./cmd/pseudoc -in examples/div_mod.pseudo -out build/div_mod.go
go run ./build/div_mod.go
```

Tampilkan token lexer:

```bash
go run ./cmd/pseudoc -in examples/div_mod.pseudo -tokens
```

Tanpa `-out`, hasil Go dicetak ke stdout:

```bash
go run ./cmd/pseudoc -in examples/div_mod.pseudo
```

## Sintaks yang didukung

```text
program NamaProgram
kamus
  nilai : integer
  a, b : real
  constant batas : integer <- 10
  type waktu <
    jam : integer
    menit : integer
  >
algoritma
  nilai <- 2 + 3 * 4
  input(a)
  output("hasil:", nilai)
  for i <- 1 to 5 do
    nilai <- nilai + i
  endfor
endprogram
```

- Tipe dasar: `boolean`, `string`, `char`, `integer`/`int`, dan `real`.
- Deklarasi variabel satu atau beberapa nama dalam satu baris.
- Konstanta dengan `constant nama : tipe <- ekspresi`.
- Deklarasi record `type Nama < ... >` dan variabel dengan tipe record tersebut.
- Alias variabel dengan `nama_alias : nama_variabel_yang_sudah_dideklarasikan`; alias menunjuk ke penyimpanan variabel yang sama.
- Assignment dengan `<-`.
- `input(...)` dan `output(...)`. Input membaca token yang dipisahkan whitespace; untuk `char`, token pertama dibaca sebagai satu rune.
- Perulangan inklusif `for i <- awal to akhir do ... endfor` dan `downto`.
- Operator aritmetika `+`, `-`, `*`, `DIV`, `MOD`.
- Operator logika `NOT`, `AND`, `OR`, `XOR`.
- Perbandingan `=`, `==`, `/=`, `!=`, `<`, `<=`, `>`, `>=`.
- Tanda kurung, literal integer/real/string/char, dan boolean `true`/`false` atau `benar`/`salah`.
- Pemisah statement/deklarasi berupa baris baru atau `;`.

## Aturan tipe

- Pemeriksaan tipe dilakukan statis sebelum code generation.
- Assignment harus memiliki tipe yang sama; tidak ada konversi integer-real otomatis.
- `DIV` dan `MOD` hanya menerima dua operand `integer` dan menghasilkan `integer`.
- `AND`, `OR`, `XOR`, dan `NOT` hanya menerima boolean.
- Batas `for` harus bertipe `integer`.
- Konstanta tidak boleh diubah atau menjadi target `input`.

## Batasan dan asumsi yang perlu diverifikasi dengan dosen/modul

Materi README sumber menyebut `IF ELSE` belum diajarkan, jadi konstruksi tersebut sengaja belum dimasukkan. Materi yang tersedia juga belum menetapkan tata bahasa komentar, aturan case-sensitive, sintaks perbandingan secara formal, literal boolean/char, akses field record (`wt.jam`), maupun perilaku `DIV`/`MOD` untuk bilangan negatif. Implementasi ini mengambil asumsi berikut:

1. Keyword dan operator kata (`DIV`, `MOD`, `AND`, dan seterusnya) case-insensitive; nama variabel dan nama tipe buatan sendiri case-sensitive.
2. `true`/`false` serta `benar`/`salah` diterima sebagai literal boolean.
3. String memakai tanda kutip ganda dan char memakai tanda kutip tunggal dengan escape ala Go.
4. `for` berjalan sampai batas akhir secara inklusif.
5. `DIV` dan `MOD` mengikuti pembagian dan remainder integer Go. Perilaku pada operand negatif perlu dicocokkan dengan definisi mata kuliah.
6. Record bisa dideklarasikan sebagai tipe dan variabel, tetapi akses/assignment field belum didukung karena sintaksnya belum terlihat di materi yang dirujuk.
7. Komentar belum didukung. Jangan menambahkan sintaks komentar sebelum aturan resminya diketahui.

Lexer dan parser menghentikan proses pada kesalahan pertama karena token/struktur setelahnya mungkin tidak dapat dipercaya. Type checker mengumpulkan beberapa diagnostic semantik dalam satu proses.

## Struktur

```text
cmd/pseudoc/             CLI
internal/compiler/       lexer, parser, AST, type checker, Go code generator
examples/                contoh source pseudocode
```

## Repositori materi rujukan

- https://github.com/fkesa/algoritma-pemrograman
- https://github.com/fkesa/algoritma-pemrograman-praktikum

## Contoh yang disertakan

- `examples/hello.pseudo`: input string dan integer serta output.
- `examples/div_mod.pseudo`: pembagian integer dan sisa bagi.
- `examples/loop.pseudo`: loop `to` dan `downto`.
- `examples/record.pseudo`: record, konstanta, dan tipe custom.
- `examples/alias.pseudo`: dua nama untuk storage variabel yang sama.
- `examples/literals.pseudo`: real, char, boolean, dan operator logika.
- `examples/input_char.pseudo`: pembacaan karakter.
