package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
)

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

func getConfigFileName() string {
	d := fmt.Sprintf("%s/.secrets/", os.Getenv("HOME"))
	f := "config"

	if !fileExists(d) {
		os.Mkdir(d, 0755)
	}

	return fmt.Sprintf("%s/%s", d, f)

}

func setConfigDefaults() (secretsConfig, error) {

	var sc secretsConfig
	sc.VaultAddress = "http://127.0.0.1:9000"
	sc.AuthToken = ""

	// encrypt
	b, e := encConfig(sc)
	if e != nil {
		log.Fatalf("Failed to encrypt configuration: %s", e)
	}

	// write out config
	os.WriteFile(getConfigFileName(), b, 0600)
	if e != nil {
		log.Fatalf("Failed to write encrypted configuration file: %s", e)
	}

	// done
	return sc, nil

}
