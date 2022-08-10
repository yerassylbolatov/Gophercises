package main

import "fmt"

func caesarCipher(r rune, delta int32) rune {
	if r >= 'a' && r <= 'z' {
		return cipher(r, delta, 'a')
	}
	if r >= 'A' && r <= 'Z' {
		return cipher(r, delta, 'A')
	}
	return r
}

func cipher(r rune, delta int32, base rune) rune {
	tmp := r - base
	tmp = (tmp + delta) % 26
	return tmp + base
}

func main() {
	var res []rune
	input := "middle-Outz"
	var delta int32 = 2
	for _, r := range input {
		res = append(res, caesarCipher(r, delta))
	}
	for _, r := range string(res) {
		fmt.Println(string(r))
	}
}
