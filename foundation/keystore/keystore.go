// Package keystore implements the auth.KeyStore interface. This implements
// an in-memory keystore for JWT support.
package keystore

import (
	"crypto/rsa"
	"sync"
)

// KeyStore represents an in memory store implementation of the KeyStorer
// interface for use with the auth package.
type KeyStore struct {
	mu    sync.RWMutex
	store map[string]*rsa.PrivateKey
}

// New construct an empty KeyStore ready for use.
func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]*rsa.PrivateKey),
	}
}

func NewMap(store map[string]*rsa.PrivateKey) *KeyStore {
	return &KeyStore{store: store}
}
