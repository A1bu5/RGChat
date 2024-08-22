// encryption.go

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// encryption.go
func decrypt(key []byte, data []byte) ([]byte, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("encrypted data is too short")
	}

	nonce, ciphertext := data[:12], data[12:]

	fmt.Printf("Received Nonce: %x\n", nonce)
	fmt.Printf("Received Ciphertext: %x\n", ciphertext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption error: %v", err)
	}

	return plaintext, nil
}
