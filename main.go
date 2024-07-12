package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/ovh/configstore"
	"github.com/ovh/symmecrypt/keyloader"
)

const EncryptionKeyCtx = "encryption-key"

func generateSecureKey() string {
	b := make([]byte, 32)
	// Read 32 random bytes from the cryptographically secure random number generator
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error generating random bytes:", err)
		return ""
	}
	// Print the random bytes in hexadecimal format
	return fmt.Sprintf("%x", b)
}

func main() {
	url := "host=localhost user=parrot password=parrot123 dbname=parrot port=5432 sslmode=disable"

	secretKey := generateSecureKey()
	kk := keyloader.KeyConfig{
		Identifier: EncryptionKeyCtx,
		Cipher:     "aes-gcm",
		Timestamp:  time.Now().UnixMilli(),
		Sealed:     false,
		Key:        string(secretKey),
	}

	keyList := []configstore.Item{}
	keyList = append(keyList, configstore.NewItem(EncryptionKeyCtx, kk.String(), 1))
	configstore.RegisterProvider("env", func() (configstore.ItemList, error) {
		return configstore.ItemList{
			Items: keyList,
		}, nil
	})

	fmt.Printf("Secret Key: %s\n", secretKey)
	HowToEncrypt()
	EncryptUsingGorm(url)
	EncryptUsingGorm(url)
}
