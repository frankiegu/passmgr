// Copyright (c) 2017, David Url
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package filestore

import (
	"testing"
)

func TestNonceRotates(t *testing.T) {
	gcm := aesGcm{nonce: []byte{0xAB, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}}
	gcm.incrementNonce()
	assertEqual(t,
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		gcm.nonce[4:],
		"")
}

func TestNonceIncrements(t *testing.T) {
	gcm := aesGcm{
		nonce: []byte{0xAB, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		AEAD:  &mockAead{},
	}
	_, err := gcm.Encrypt([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t,
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		gcm.nonce[4:],
		"")
}

func TestNonceRandomizedPart(t *testing.T) {
	gcm := aesGcm{
		nonce: []byte{0xAB, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		AEAD:  &mockAead{},
	}
	_, err := gcm.Encrypt([]byte{})
	if err != nil {
		t.Fatal(err)
	}
	assertNotEqual(t,
		[]byte{0xAB, 0x00, 0x00, 0x00},
		gcm.nonce[:4],
		"")
}

func TestNonceDoesNotChange(t *testing.T) {
	gcm := aesGcm{
		nonce: []byte{0xAB, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF},
		AEAD:  &mockAead{},
	}
	_, err := gcm.Decrypt([]byte{0xAB, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xAB})
	if err != nil {
		t.Fatal(err)
	}
	assertEqual(t,
		[]byte{0xAB, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF},
		gcm.nonce,
		"")
}

func TestCiphertextTooShort(t *testing.T) {
	gcm := aesGcm{
		nonce: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		AEAD:  &mockAead{},
	}
	_, err := gcm.Decrypt([]byte{0xAB, 0xAB})
	if err == nil {
		t.Error("len(ciphertext)<nonceSize should result in error")
	}
}

type mockAead struct {
}

func (a *mockAead) NonceSize() int {
	return 12
}

func (a *mockAead) Overhead() int {
	return 0
}

func (a *mockAead) Seal(dst []byte, nonce []byte, plaintext []byte, additionalData []byte) []byte {
	return dst
}

func (a *mockAead) Open(dst []byte, nonce []byte, ciphertext []byte, additionalData []byte) ([]byte, error) {
	return dst, nil
}
