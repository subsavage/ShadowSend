package main

import (
	"fmt"
	"github.com/subsavage/ShadowSend/internal/crypto"
)

func main() {
	pass, _ := crypto.GeneratePassphrase(12)
	fmt.Println("Generated passphrase:", pass)

	original := []byte("This is my top secret file content.")

	encrypted, nonce, salt, _ := crypto.Encrypt(original, pass)
	fmt.Println("Encrypted:", encrypted)

	decrypted, err := crypto.Decrypt(encrypted, nonce, salt, pass)
	if err != nil {
		panic(err)
	}

	fmt.Println("Decrypted:", string(decrypted))
}
