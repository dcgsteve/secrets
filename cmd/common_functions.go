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
	"syscall"

	"github.com/hashicorp/vault/api"
	"golang.org/x/term"
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

func stop(m ...string) {
	fmt.Println(m)
	os.Exit(0)
}

func stopFatal(m ...string) {
	fmt.Println(m)
	os.Exit(1)
}

func cliGetClient() *api.Client {

	c, e := getClient()
	if e != nil {
		stopFatal("Failed to create Vault client: ", e.Error())
	}

	return c
}

func getClient() (*api.Client, error) {

	// init client
	client, e := api.NewClient(&api.Config{Address: sc.VaultAddress, HttpClient: httpClient})
	if e != nil {
		return nil, e
	}

	// set token
	client.SetToken(sc.AuthToken)

	_, e = client.Auth().Token().LookupSelf()

	if e != nil {
		switch getStatusCode(e) {
		case 403:
			if t, e := getNewToken(client); e != nil {
				return nil, e
			} else {
				sc.AuthToken = t
				client.SetToken(sc.AuthToken)
				if e := saveConfig(); e != nil {
					return nil, e
				}
			}
		default:
			return nil, e
		}
	}

	// return valid client
	return client, nil
}

func getNewToken(c *api.Client) (string, error) {

	var u, p string
	var e error

	// get username (note: this must match what is in Vault)
	if sc.Username == "" {
		fmt.Printf("Vault Username: ")
		fmt.Scanln(&u)
		fmt.Println()
		if e != nil {
			return "", errors.New("failed to autheticate against Vault with username and password")
		}
	} else {
		u = sc.Username
	}

	// get password from user
	if sc.Password == "" {
		fmt.Printf("Vault password for user %q (will be hidden): ", u)
		b, e := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if e != nil {
			return "", errors.New("failed to autheticate against Vault with username and password")
		}
		p = string(b)
	} else {
		p = sc.Password
	}

	// enforce userpass authentication
	path := fmt.Sprintf("auth/userpass/login/%s", u)

	// auth
	s, e := c.Logical().Write(path, map[string]interface{}{"password": string(p)})
	if e != nil {
		return "", errors.New("failed to autheticate against Vault with username and password")
	}

	// return valid token
	return s.Auth.ClientToken, nil
}

func getStatusCode(e error) int {

	return e.(*api.ResponseError).StatusCode
}

func checkConfig() {

	var ok = true

	if sc.Password == "" {
		ok = false
	}
	if sc.Project == "" {
		ok = false
	}
	if sc.Store == "" {
		ok = false
	}
	if sc.Username == "" {
		ok = false
	}
	if sc.VaultAddress == "" {
		ok = false
	}

	if !ok {
		stop("All configuration details have not yet been entered!")
	}
}
