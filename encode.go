package main

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strings"
)

var base64symbols string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
var codeSymbols string = "abcdefghijklmnopqrstuvwxyz0123456789"

func encode(raw string, symbols string) *big.Int {
	s := base64.StdEncoding.EncodeToString([]byte(raw))
	B := big.NewInt(int64(len(symbols)))
	num := big.NewInt(0)

	for _, ch := range s {
		num.Mul(num, B)
		num.Add(num, big.NewInt(int64(strings.Index(symbols, string(ch)))))
	}

	return num
}

func decode(num *big.Int, symbols string) string {
	B := big.NewInt(int64(len(symbols)))
	sb := ""

	fmt.Println(num.Cmp(big.NewInt(0)))
	mod := big.NewInt(0)
	for num.Cmp(big.NewInt(0)) > 0 {
		mod.Set(num)
		mod.Mod(mod, B)
		sb += string(symbols[mod.Int64()])

		num.Div(num, B)
	}

	dec, _ := base64.StdEncoding.DecodeString(reverse(sb))

	return string(dec)
}

func reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}

	return string(rns)
}
