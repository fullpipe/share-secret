package main

import (
	"bytes"
	"encoding/base64"
)

type Pool struct {
	Shares [][]byte
}

func NewPool() *Pool {
	return &Pool{Shares: [][]byte{}}
}

func (p *Pool) Add(share string) {
	decoded, err := base64.StdEncoding.DecodeString(share)
	if err != nil {
		// on error we add invalid share to fail unseal
		decoded = []byte("fake")
	}

	for _, s := range p.Shares {
		if bytes.Compare(s, decoded) == 0 {
			return
		}
	}

	p.Shares = append(p.Shares, decoded)
}
