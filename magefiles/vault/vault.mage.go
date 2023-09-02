package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/constants"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// DSV is the namespace for mage tasks related to DSV, such as client credential creation.
type (
	DSV               mg.Namespace
	ClientCredentials struct {
		ClientID string `json:"clientId"`     //nolint:tagliatelle // json tag required as is
		Secret   string `json:"clientSecret"` //nolint:tagliatelle // json tag required as is
	}
)

var (
	dsvprofilename  = os.Getenv("DSV_PROFILE_NAME")
	rolename        = "dsv-k8s-tests"
	secretpath      = fmt.Sprintf("tests:%s", "dsv-k8s")
	policyname      = fmt.Sprintf("secrets:%s", secretpath)
	policysubjects  = fmt.Sprintf("roles:%s", rolename)
	policyresources = fmt.Sprintf("secrets:%s:<.*>", secretpath)

	// secretpathclient = fmt.Sprintf("clients:%s", secretpath)
	desc           = "a secret for testing operation of with dsv-k8s"
	clientcredfile = filepath.Join(constants.CacheDirectory, fmt.Sprintf("%s.json", rolename))
	clientcredname = rolename
	secretkey      = "food" // just simple test placeholder for now
	testsecretkey  = fmt.Sprintf("secrets:%s:%s", secretpath, secretkey)
	//nolint:gosec // test value, so fine to leave hard coded
	testsecretvalue = `
{
	"taco":"burrito",
	"username": "tacoeater",
	"domain": "tacoeater.com"
}
` //  placeholder for testing, not sensitive, and ok to leave for now
)

// checkDSVProfileName checks if the DSV_PROFILE_NAME is set and returns an error if not.
func checkDSVProfileName() error {
	if dsvprofilename == "" {
		pterm.Error.Println(
			"DSV_PROFILE_NAME is not set and this is required to ensure the correct dsv tenant for testing is used",
		)
		return fmt.Errorf("DSV_PROFILE_NAME is required")
	}
	return nil
}

// ➕ SetupDSV creates the policy, role, and client credentials.
func (DSV) SetupDSV() error {
	magetoolsutils.CheckPtermDebug()
	if err := checkDSVProfileName(); err != nil {
		pterm.Error.Println("DSV_PROFILE_NAME is not set and this is required to automate the setup of the test credentials")
		return fmt.Errorf("DSV_PROFILE_NAME is required: %w", err)
	}
	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelInfo).WithCaller(true)
	if _, err := os.Stat(clientcredfile); err == nil {
		logger.Error("client credentials already exist, please run `mage vault:destroy` to remove and try again")
		return fmt.Errorf("tear down existing test credentials before recreate to avoid conflicts")
	}
	// dsv role create
	logger.Info("creating role", logger.Args("rolename", rolename))

	if err := sh.RunV("dsv", "role", "create", "--name", rolename, "--profile", dsvprofilename); err != nil {
		logger.Warn("unable to create role", logger.Args("rolename", rolename))
	}
	logger.Info("created role", logger.Args("rolename", rolename))

	// dsv policy create
	if err := sh.RunV("dsv", "policy", "create",
		"--path", policyname,
		"--actions", "read,list",
		"--effect", "allow",
		"--subjects", policysubjects,
		"--desc", fmt.Sprintf("scoped access for %s by %s", secretpath, rolename),
		"--resources", policyresources,
		"--profile", dsvprofilename,
	); err != nil {
		logger.Warn("unable to create policy", logger.Args("policyname", rolename))
	}
	logger.Info("created policy", logger.Args("policyname", rolename))

	logger.Info("creating client credentials", logger.Args("clientcredname", clientcredname))
	err := sh.RunV(
		"dsv",
		"client",
		"create",
		"--role", rolename,
		"--plain",
		"--profile", dsvprofilename,
		"--out", fmt.Sprintf("file:%s", clientcredfile),
	)
	if err != nil {
		logger.Warn("unable to create client credentials", logger.Args("clientcredname", clientcredname))
	}
	logger.Info("created client credentials", logger.Args("clientcredname", clientcredname))

	b, err := os.ReadFile(clientcredfile)
	if err != nil {
		logger.Error(
			"unable to read client credentials file",
			logger.Args("clientcredfile", clientcredfile, "error", err),
		)
		return err
	}
	var clientcred ClientCredentials
	err = json.Unmarshal(b, &clientcred)
	if err != nil {
		logger.Error(
			"unable to unmarshal client credentials file",
			logger.Args("clientcredfile", clientcredfile, "error", err),
		)
		return err
	}

	logger.Info("Put in .cache/charts/dsv-k8s/values.yaml", logger.Args(
		"clientID", clientcred.ClientID,
		"clientSecret", clientcred.Secret,
	))

	return nil
}

// 🔐 CreateSecret creates a secret for usage with this specific client, policy, and role setup.
// This probably needs refactoring to allow input via pterm or via file.
// At time of creation (2023-04) it's a draft task to help with better test setup for developers wanting to test and have isolated
// permissions for just this specific secret path, role, client. It's all hard coded but can improve in the future.
func (DSV) CreateSecret() error {
	magetoolsutils.CheckPtermDebug()
	if err := checkDSVProfileName(); err != nil {
		pterm.Error.Println("DSV_PROFILE_NAME is not set and this is required to automate the setup of the test credentials")
		return fmt.Errorf("DSV_PROFILE_NAME is required: %w", err)
	}

	logger := pterm.DefaultLogger.WithLevel(pterm.LogLevelInfo).WithCaller(true)
	logger.Info("creating secret for DSV client")
	secretkey := "food"
	if err := sh.RunV(
		"dsv",
		"secret",
		"create",
		"--path", testsecretkey,
		"--data", testsecretvalue,
		"--desc", desc,
		"--profile", dsvprofilename,
	); err != nil {
		logger.Error("unable to create secret", logger.Args("secretkey", secretkey, "error", err))
		return err
	}
	logger.Info("created secret for DSV client", logger.Args("secretkey", secretkey))
	return nil
}

// ConvertClientToCredentials reads the client credentials created in .cache and converts to the format the helm chart/injector expect.
func (DSV) ConvertClientToCredentials() error {
	if os.Getenv("DSV_TENANT_NAME") == "" {
		return fmt.Errorf("DSV_TENANT_NAME is required, make sure you've set in .env and run `direnv allow`")
	}
	// Read the input JSON file
	input, err := os.ReadFile(clientcredfile)
	if err != nil {
		return err
	}

	// Unmarshal the input JSON into a struct
	var data struct {
		ClientID     string `json:"clientId"`     //nolint:tagliatelle // json tag required as is
		ClientSecret string `json:"clientSecret"` //nolint:tagliatelle // json tag required as is
	}
	if err := json.Unmarshal(input, &data); err != nil {
		return err
	}

	// Create the output JSON struct
	output := struct {
		Default struct {
			Credentials struct {
				ClientID     string `json:"clientId"`     //nolint:tagliatelle // json tag required as is
				ClientSecret string `json:"clientSecret"` //nolint:tagliatelle // json tag required as is
			} `json:"credentials"`
			Tenant string `json:"tenant"`
		} `json:"default"`
	}{
		Default: struct {
			Credentials struct {
				ClientID     string `json:"clientId"`
				ClientSecret string `json:"clientSecret"`
			} `json:"credentials"`
			Tenant string `json:"tenant"`
		}{
			Credentials: struct {
				ClientID     string `json:"clientId"`
				ClientSecret string `json:"clientSecret"`
			}{
				ClientID:     data.ClientID,
				ClientSecret: data.ClientSecret,
			},
			Tenant: os.Getenv("DSV_TENANT_NAME"),
		},
	}

	// Marshal the output JSON
	outputJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	// Write the output JSON to a file
	if err := os.WriteFile(constants.CacheCredentialFile, outputJSON, constants.PermissionUserReadWriteExecute); err != nil {
		return err
	}

	return nil
}

// 🗑️ Destroy tears down the test credentials hard setup by this automation only
func (DSV) Destroy() {
	// if .cache/credentials.json exists, then read and get the clientid. use this as input for the client delete
	var clientcredid string
	if _, err := os.Stat(clientcredfile); err == nil {
		// file exists
		b, err := os.ReadFile(clientcredfile)
		if err != nil {
			pterm.Warning.Println("unable to read client credentials file, this may be expected if already deleted")
		}
		var clientcred ClientCredentials
		err = json.Unmarshal(b, &clientcred)
		if err != nil {
			pterm.Warning.Println("unable to unmarshal client credentials file, this may be expected if already deleted")
		}
		clientcredid = clientcred.ClientID
		// run dsv client delete against the test credentials, warn only on error
		if err := sh.RunV("dsv", "client", "delete", "--client-id", clientcredid, "--profile", dsvprofilename, "--force"); err != nil {
			pterm.Warning.Println("unable to delete client credentials, this may be expected if already deleted")
		}
	}

	// run dsv policy delete against the test role, warn only on error
	if err := sh.RunV("dsv", "policy", "delete", "--path", policyname, "--profile", dsvprofilename, "--force"); err != nil {
		pterm.Warning.Println("unable to delete policy, this may be expected if already deleted")
	}
	// run dsv role delete against the test role, warn only on error
	if err := sh.RunV("dsv", "role", "delete", "--name", rolename, "--profile", dsvprofilename, "--force"); err != nil {
		pterm.Warning.Println("unable to delete role, this may be expected if already deleted")
	}

	// remove the test secret with force, and warn only on error
	if err := sh.RunV("dsv", "secret", "delete", "--path", testsecretkey, "--profile", dsvprofilename, "--force"); err != nil {
		pterm.Warning.Println("unable to delete secret, this may be expected if already deleted")
	}

	_ = sh.Rm(clientcredfile)
}
