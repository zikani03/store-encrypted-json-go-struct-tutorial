package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/ovh/symmecrypt/keyloader"
)

func HowToEncrypt() {

	k, err := keyloader.LoadSingleKey()
	if err != nil {
		log.Fatalf("failed to load encryption key")
		return
	}

	encrypted, err := k.Encrypt([]byte("Hello World, Please encrypt me!"))
	if err != nil {
		log.Fatalf("failed to encrypt data %v", err)
		return
	}

	fmt.Println("Encrypted, encoded as Base64", base64.StdEncoding.EncodeToString(encrypted))
}
