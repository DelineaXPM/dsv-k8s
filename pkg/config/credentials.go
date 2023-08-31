package config

import (
	"encoding/json"
	"fmt"
	"io"
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
		return nil, fmt.Errorf("unable to open configuration file '%s': %w", credentialsFilePath, err)
	} else {
		defer credentialsFile.Close()
		return GetCredentialsFromFile(credentialsFile)
	}
}

// GetCredentialsFromFile parses the credentialsFile and returns the resulting Credentials object
func GetCredentialsFromFile(credentialsFile *os.File) (*Credentials, error) {
	if contents, err := io.ReadAll(credentialsFile); err != nil {
		return nil, fmt.Errorf("unable to read configuration file '%s': %w", credentialsFile.Name(), err)
	} else {
		return MakeCredentials(contents)
	}
}

func MakeCredentials(credentialJSON []byte) (*Credentials, error) {
	credentials := new(Credentials)

	if err := json.Unmarshal(credentialJSON, credentials); err != nil {
		return nil, fmt.Errorf("unable to unmarshal configuration: %w", err)
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
