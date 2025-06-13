package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"golang.org/x/crypto/pbkdf2"
	"io"
	"math/big"
)

const (
	KeyLength   = 32
	SaltLength  = 16
	NonceLength = 12
	PBKDF2Iter  = 100_000
)

var charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

func GeneratePassphrase(length int) (string, error) {
	pass := make([]byte, length)

	for i := range pass {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))
		if err != nil {
			return "Error: ", err
		}
		pass[i] = charSet[n.Int64()]
	}
	return string(pass), nil
}

func deriveKey(passphrase string, salt []byte) []byte {
	return pbkdf2.Key([]byte(passphrase), salt, PBKDF2Iter, KeyLength, sha256.New)
}

func Encrypt(data []byte, passphrase string) (ciphertext, nonce, salt []byte, err error) {
	salt = make([]byte, SaltLength)
	if _, err = io.ReadFull(rand.Reader, salt); err != nil {
		return
	}
	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	nonce = make([]byte, NonceLength)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	ciphertext = aesgcm.Seal(nil, nonce, data, nil)
	return
}

func Decrypt(ciphertext, nonce, salt []byte, passphrase string) ([]byte, error) {
	key := deriveKey(passphrase, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("wrong passphrase or corrupted data")
	}

	return plaintext, nil
}