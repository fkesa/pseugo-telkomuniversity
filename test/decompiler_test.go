package test

import (
	"os"
	"strings"
	"testing"

	"pseudocode-compiler/internal/compiler"
)

func TestDecompileSimpleGo(t *testing.T) {
	goSrc := `package main

import "fmt"

func main() {
	var a, b int
	fmt.Scan(&a, &b)
	c := a + b
	fmt.Println("Hasil:", c)
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("DecompileGo failed: %v", err)
	}

	if !strings.Contains(pseudo, "program main") {
		t.Errorf("expected program main, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "input(a, b)") {
		t.Errorf("expected input(a, b), got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "c <- a + b") {
		t.Errorf("expected c <- a + b, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "output(\"Hasil:\", c)") {
		t.Errorf("expected output(\"Hasil:\", c), got:\n%s", pseudo)
	}

	// Verify that the decompiled pseudocode can be compiled back into Go!
	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("Compile decompiled pseudocode failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestDecompileDivMod(t *testing.T) {
	goSrc := `package main

import "fmt"

func Pembagian() {
	var bil1, bil2 int
	var hasil_bagi int
	var sisa_bagi int

	bil1 = 100
	bil2 = 8
	hasil_bagi = (bil1 / bil2)
	sisa_bagi = (bil1 % bil2)
	fmt.Println(hasil_bagi, sisa_bagi)
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("DecompileGo failed: %v", err)
	}

	if !strings.Contains(pseudo, "bil1 DIV bil2") {
		t.Errorf("expected DIV, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "bil1 MOD bil2") {
		t.Errorf("expected MOD, got:\n%s", pseudo)
	}

	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("Compile decompiled pseudocode failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestDecompileIfElse(t *testing.T) {
	goSrc := `package main

import "fmt"

func main() {
	var x int
	fmt.Scan(&x)
	if x > 10 {
		fmt.Println("Lebih dari 10")
	} else if x == 10 {
		fmt.Println("Sama dengan 10")
	} else {
		fmt.Println("Kurang dari 10")
	}
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("DecompileGo failed: %v", err)
	}

	if !strings.Contains(pseudo, "if x > 10 then") {
		t.Errorf("expected if x > 10 then, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "elseif x = 10 then") {
		t.Errorf("expected elseif x = 10 then, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "else") || !strings.Contains(pseudo, "endif") {
		t.Errorf("expected else and endif, got:\n%s", pseudo)
	}

	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("Compile decompiled pseudocode failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestDecompileLoops(t *testing.T) {
	goSrc := `package main

import "fmt"

func main() {
	var n, total int
	fmt.Scan(&n)
	total = 0
	for i := 1; i <= n; i++ {
		total = total + i
	}
	for j := 5; j >= 1; j-- {
		fmt.Println(j)
	}
	for total > 100 {
		total = total - 10
	}
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("DecompileGo failed: %v", err)
	}

	if !strings.Contains(pseudo, "for i <- 1 to n do") {
		t.Errorf("expected for i <- 1 to n do, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "for j <- 5 downto 1 do") {
		t.Errorf("expected for j <- 5 downto 1 do, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "while total > 100 do") {
		t.Errorf("expected while total > 100 do, got:\n%s", pseudo)
	}

	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("Compile decompiled pseudocode failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestDecompileRepeatUntil(t *testing.T) {
	goSrc := `package main

import "fmt"

func main() {
	var x int
	x = 0
	for {
		x = x + 1
		if x >= 10 {
			break
		}
	}
	fmt.Println(x)
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("DecompileGo failed: %v", err)
	}

	if !strings.Contains(pseudo, "repeat") || !strings.Contains(pseudo, "until x >= 10") {
		t.Errorf("expected repeat ... until x >= 10, got:\n%s", pseudo)
	}

	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("Compile decompiled pseudocode failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestDecompileRecord(t *testing.T) {
	goSrc := `package main

import "fmt"

type prisma struct {
	panjang float64
	lebar   float64
	tinggi  float64
	volume  float64
}

func main() {
	var pr prisma
	fmt.Scan(&pr.panjang, &pr.lebar, &pr.tinggi)
	pr.volume = pr.panjang * pr.lebar * pr.tinggi
	fmt.Println(pr.volume)
}
`
	pseudo, err := compiler.DecompileGo(goSrc)
	if err != nil {
		t.Fatalf("DecompileGo failed: %v", err)
	}

	if !strings.Contains(pseudo, "type prisma <") {
		t.Errorf("expected type prisma <, got:\n%s", pseudo)
	}
	if !strings.Contains(pseudo, "pr.volume <- ((pr.panjang * pr.lebar) * pr.tinggi)") &&
		!strings.Contains(pseudo, "pr.volume <- pr.panjang * pr.lebar * pr.tinggi") {
		t.Errorf("expected pr.volume assignment, got:\n%s", pseudo)
	}

	if _, err := compiler.Compile(pseudo); err != nil {
		t.Fatalf("Compile decompiled pseudocode failed: %v\nPseudocode:\n%s", err, pseudo)
	}
}

func TestRoundTripCompilation(t *testing.T) {
	origPseudo := `program Lingkaran
kamus
    r : real
    luas : real
    constant pi : real <- 3.14
algoritma
    input(r)
    luas <- pi * r * r
    output("Luas:", luas)
endprogram`

	res, err := compiler.Compile(origPseudo)
	if err != nil {
		t.Fatalf("original compile failed: %v", err)
	}

	decompiled, err := compiler.DecompileGo(res.GoSource)
	if err != nil {
		t.Fatalf("decompile failed: %v", err)
	}

	if _, err := compiler.Compile(decompiled); err != nil {
		t.Fatalf("compiling decompiled pseudocode failed: %v\nDecompiled:\n%s", err, decompiled)
	}
}

func TestDecompileModul09Tugas(t *testing.T) {
	origPseudo, err := os.ReadFile("../examples/modul09_tugas.pseudo")
	if err != nil {
		t.Skip("skipping if examples not found relative:", err)
	}

	res, err := compiler.Compile(string(origPseudo))
	if err != nil {
		t.Fatalf("Compile modul09_tugas.pseudo failed: %v", err)
	}

	decompiled, err := compiler.DecompileGo(res.GoSource)
	if err != nil {
		t.Fatalf("DecompileGo failed on modul09_tugas: %v", err)
	}

	if _, err := compiler.Compile(decompiled); err != nil {
		t.Fatalf("re-compiling decompiled modul09_tugas failed: %v\nDecompiled:\n%s", err, decompiled)
	}
}
