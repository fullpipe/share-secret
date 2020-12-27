package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"log"
	"strings"
	"time"

	sss "github.com/dsprenkels/sss-go"
	"github.com/google/uuid"
)

type Secret struct {
	UnsealPools map[string]*Pool
}

type CreateInput struct {
	Name, Secret                string
	KeepersNum, KeepersRequired int
}

type UnsealInput struct {
	Pool, Share string
}

type Share string
type Shares []Share

func NewSecret() *Secret {
	s := &Secret{UnsealPools: make(map[string]*Pool)}

	return s

}

func (e *Secret) StartUnseal(_, reply *string) error {
	if len(e.UnsealPools) >= 10000 {
		return errors.New("to_many_unseals")
	}

	id := uuid.New().String()
	*reply = id

	e.UnsealPools[id] = NewPool()

	go func() {
		<-time.After(10 * time.Minute)
		delete(e.UnsealPools, id)
	}()

	return nil
}

func (e *Secret) Unseal(in *UnsealInput, reply *string) error {
	pool, ok := e.UnsealPools[in.Pool]
	if !ok {
		return errors.New("no_unseal")
	}

	pool.Add(in.Share)

	restored, err := sss.CombineShares(pool.Shares)
	if err != nil {
		return errors.New("restoration_failed")
	}

	secret := string(bytes.Trim(restored, "\x00"))

	if !strings.HasPrefix(secret, "~!") || !strings.HasSuffix(secret, "!~") {
		return errors.New("need_more_shares")
	}

	secret = strings.TrimPrefix(secret, "~!")
	secret = strings.TrimSuffix(secret, "!~")
	secret = strings.TrimSpace(secret)

	*reply = secret

	return nil
}

func (e *Secret) Seal(in *CreateInput, reply *Shares) error {
	in.Secret = strings.TrimSpace(in.Secret)
	if len([]byte(in.Secret)) > 60 {
		return errors.New("Secret is too long")
	}

	secret := []byte("~!" + in.Secret + "!~")

	if in.KeepersNum > 255 {
		return errors.New("Too many secret keepers")
	}

	if in.KeepersRequired > in.KeepersNum {
		return errors.New("The number of keepers required must be less than their total number")
	}

	data := make([]byte, 64)
	copy(data, secret)

	shares, err := sss.CreateShares(data, in.KeepersNum, in.KeepersRequired)
	if err != nil {
		log.Fatalln(err)
	}

	for _, s := range shares {
		*reply = append(*reply, Share(base64.StdEncoding.EncodeToString(s)))
	}

	return nil
}
