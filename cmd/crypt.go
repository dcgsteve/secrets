package cmd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"io"
)

// Proper key injected at build time
var EncryptionKey string = "dummy"

func encConfig(sc secretsConfig) ([]byte, error) {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	e := enc.Encode(sc)
	if e != nil {
		return nil, e
	}

	r, e := encBytes(buf.Bytes(), EncryptionKey)
	if e != nil {
		return nil, e
	}

	return r, nil

}

func decConfig(c []byte) (secretsConfig, error) {

	var sc secretsConfig

	d, e := decBytes(c, EncryptionKey)
	if e != nil {
		return sc, e
	}

	buf := bytes.NewBuffer(d)
	dec := gob.NewDecoder(buf)

	e = dec.Decode(&sc)
	if e != nil {
		return sc, e
	}

	return sc, nil
}
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
		panic(err.Error())
	}
	return plaintext, nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
