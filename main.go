package main

import (
	"fmt"
	"time"

	"github.com/subsavage/ShadowSend/internal/storage"
)

func main() {
	db, err := storage.NewDB("shadowsend.db")
	if err != nil {
		panic(err)
	}

	id := "abc123"
	meta := storage.FileMetadata{
		Filename:  "secret.txt",
		Data:      []byte("Top secret bytes here..."),
		Nonce:     []byte("123456789012"), // dummy nonce
		Salt:      []byte("salt1234abcd5678"),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := db.SaveFile(id, meta); err != nil {
		panic(err)
	}

	loaded, err := db.GetFile(id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Retrieved filename:", loaded.Filename)
	fmt.Println("Content:", string(loaded.Data))
}
