package main

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool_Add(t *testing.T) {
	pool := NewPool()

	pool.Add(enc("asd"))
	assert.Len(t, pool.Shares, 1)

	pool.Add(enc("qwe"))
	pool.Add(enc("qwe"))
	assert.Len(t, pool.Shares, 2)
	assert.IsType(t, [][]byte{}, pool.Shares)
}

func enc(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
