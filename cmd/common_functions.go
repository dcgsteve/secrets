package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/vault/api"
)

// Proper key injected at build time
var EncryptionKey string = "dummy"

func encBytes(data []byte, key string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(key)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decBytes(data []byte, key string) ([]byte, error) {
	keyHash := []byte(createHash(key))
	block, err := aes.NewCipher(keyHash)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt")
	}
	return plaintext, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func fileExists(f string) bool {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func fileNotExists(f string) bool {
	return !fileExists(f)
}

func stop(m string) {
	fmt.Println(m)
	os.Exit(0)
}

func cliGetClient() *api.Client {

	c, e := getClient()
	if e != nil {
		stop(fmt.Sprintf("Failed to create Vault client: %s", e))
	}

	return c
}

func getClient() (*api.Client, error) {

	// init client
	client, e := api.NewClient(&api.Config{Address: sc.VaultAddress, HttpClient: httpClient})
	if e != nil {
		return nil, e
	}
	client.SetToken(sc.AuthToken)

	return client, nil
}
