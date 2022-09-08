package cmd

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage local configuration",
	Long:  "Manages configuration information in the encrypted SECRETS confiruation file ($HOME/.secrets/config)",
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func getConfigFileName() string {
	// set path and file names
	d := fmt.Sprintf("%s/.secrets/", os.Getenv("HOME"))
	f := "config"

	// create dir if not there already
	if !fileExists(d) {
		os.Mkdir(d, 0755)
	}

	// return full path to file
	return fmt.Sprintf("%s/%s", d, f)

}

func setConfigDefaults() error {

	// set default values
	sc.VaultAddress = "http://127.0.0.1:9000"
	sc.AuthToken = ""
	sc.Project = ""

	// save
	e := saveConfig()
	if e != nil {
		return e
	}

	// done
	return nil

}

func getConfig() error {

	// read in config
	b, e := os.ReadFile(getConfigFileName())
	if e != nil {
		return e
	}

	// decrypt config
	sc, e = decConfig(b)
	if e != nil {
		return e
	}

	// done
	return nil

}

func saveConfig() error {
	// encrypt
	b, e := encConfig(sc)
	if e != nil {
		return e
	}

	// write out config
	os.WriteFile(getConfigFileName(), b, 0600)
	if e != nil {
		return e
	}

	// done
	return nil
}

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
