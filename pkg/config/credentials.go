package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/DelineaXPM/dsv-sdk-go/v2/vault"
)

// Credentials is a mapping of credentialName to dsv-sdk-go/vault/Configuration objects
type Credentials map[string]struct {
	vault.Configuration
}

// GetCredentials opens the credentialsFile and calls GetCredentialsFromFile on the resulting file
func GetCredentials(credentialsFilePath string) (*Credentials, error) {
	if credentialsFile, err := os.Open(credentialsFilePath); err != nil {
		return nil, fmt.Errorf("unable to open configuration file '%s': %s", credentialsFilePath, err)
	} else {
		defer credentialsFile.Close()
		return GetCredentialsFromFile(credentialsFile)
	}
}

// GetCredentialsFromFile parses the credentialsFile and returns the resulting Credentials object
func GetCredentialsFromFile(credentialsFile *os.File) (*Credentials, error) {
	if contents, err := ioutil.ReadAll(credentialsFile); err != nil {
		return nil, fmt.Errorf("unable to read configuration file '%s': %s", credentialsFile.Name(), err)
	} else {
		return MakeCredentials(contents)
	}
}

func MakeCredentials(credentialJson []byte) (*Credentials, error) {
	credentials := new(Credentials)

	if err := json.Unmarshal(credentialJson, credentials); err != nil {
		return nil, fmt.Errorf("unable to unmarhal configuration: %s", err)
	} else {
		return credentials, nil
	}
}

// Names returns the list of credential names
func (credentials Credentials) Names() []string {
	names := make([]string, 0, len(credentials))

	for name := range credentials {
		names = append(names, name)
	}
	return names
}
