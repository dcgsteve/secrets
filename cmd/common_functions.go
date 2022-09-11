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
	"strings"
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

func cliGetClient() *api.Client {

	c, e := getClient()
	if e != nil {
		stop("Failed to create Vault client: ", e.Error())
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

	// Force check against Vault
	_, e = client.Auth().Token().LookupSelf()

	if e != nil {
		// TODO check with Andres - better way of checking directly against StatusCode by casting error ?
		if strings.Contains(e.Error(), "permission denied") {
			fmt.Println("Your Vault token has expired - will attempt to request a new one ...")
			t, e := getNewToken(client)
			if e != nil {
				stop("Failed to create new token: ", e.Error())
			}
			sc.AuthToken = t
			e = saveConfig()
			if e != nil {
				return nil, errors.New("failed to save config after requesting new token from Vault: " + e.Error())
			}
		} else {
			return nil, errors.New("failed to check access to Vault: " + e.Error())
		}
	}

	return client, nil

}

func getNewToken(c *api.Client) (string, error) {

	var u string
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
	fmt.Printf("Vault password for user %q (will be hidden): ", u)
	b, e := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if e != nil {
		return "", errors.New("failed to autheticate against Vault with username and password")
	}
	p := map[string]interface{}{
		"password": string(b),
	}

	// enforce userpass authentication
	path := fmt.Sprintf("auth/userpass/login/%s", u)

	// auth
	s, e := c.Logical().Write(path, p)
	if e != nil {
		return "", errors.New("failed to autheticate against Vault with username and password")
	}

	// return valid token
	return s.Auth.ClientToken, nil
}

// func (h *CLIHandler) Auth(c *api.Client, m map[string]string) (string, error) {
// 	mount, ok := m["mount"]
// 	if !ok {
// 		mount = "github"
// 	}

// 	token, ok := m["token"]
// 	if !ok {
// 		if token = os.Getenv("VAULT_AUTH_GITHUB_TOKEN"); token == "" {
// 			return "", fmt.Errorf("GitHub token should be provided either as 'value' for 'token' key,\nor via an env var VAULT_AUTH_GITHUB_TOKEN")
// 		}
// 	}

// 	path := fmt.Sprintf("auth/%s/login", mount)
// 	secret, err := c.Logical().Write(path, map[string]interface{}{
// 		"token": token,
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	if secret == nil {
// 		return "", fmt.Errorf("empty response from credential provider")
// 	}

// 	return secret.Auth.ClientToken, nil
// }
