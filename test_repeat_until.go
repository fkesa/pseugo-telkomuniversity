package main

import "fmt"

func hitung_ulang() {
	var N, j int

	fmt.Scan(&N)
	j = 1
	for {
		fmt.Println(j)
		j = (j + 1)
		if j > N {
			break
		}
	}
}
func main() {
	hitung_ulang()
}
