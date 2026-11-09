package main

import "fmt"

type Mahasiswa struct {
	nama  string
	nilai int
	lulus bool
}

func main() {
	var mhs Mahasiswa
	var batasLulus int = 70

	fmt.Scan(&mhs.nama, &mhs.nilai)

	if mhs.nilai >= batasLulus {
		mhs.lulus = true
		fmt.Println(mhs.nama, "LULUS dengan nilai", mhs.nilai)
	} else {
		mhs.lulus = false
		fmt.Println(mhs.nama, "TIDAK LULUS dengan nilai", mhs.nilai)
	}

	for i := 1; i <= 3; i++ {
		fmt.Println("Selesai diproses.")
	}
}
