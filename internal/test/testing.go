package test

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/DelineaXPM/dsv-k8s/v2/pkg/config"
)

const (
	ConfigEnvVar      = "DSV_K8S_TEST_CONFIG"
	ConfigFileEnvVar  = "DSV_K8S_TEST_CONFIG_FILE"
	DefaultConfigFile = "../../configs/credentials.json"
	DefaultSecretPath = "/test/secret"
	SecretPathEnvVar  = "DSV_K8S_TEST_SECRET_PATH"
)

// Ensure log output doesn't pollute tests.
func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
}

// SecretPath returns the secret path for testing
func SecretPath() string {
	if v := os.Getenv(SecretPathEnvVar); v != "" {
		return v
	}
	return DefaultSecretPath
}

func credentialsFromFilePath(credentialsFilePath string) config.Credentials {
	if credentials, err := config.GetCredentials(credentialsFilePath); err == nil {
		return *credentials
	} else {
		panic(err)
	}
}

// Credentials returns the credentials for testing
func Credentials() config.Credentials {
	if v := os.Getenv(ConfigEnvVar); v != "" {
		if credentials, err := config.MakeCredentials([]byte(v)); err == nil {
			return *credentials
		} else {
			panic(err) // FIXME: avoid using panic and error out gracefully
		}
	} else if v := os.Getenv(ConfigFileEnvVar); v != "" {
		return credentialsFromFilePath(v)
	} else {
		return credentialsFromFilePath(DefaultConfigFile)
	}
}
