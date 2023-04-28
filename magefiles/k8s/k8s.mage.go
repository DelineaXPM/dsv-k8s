// K8s contains commands for kubectl and other kubernetes related commands.
package k8s

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/constants"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// k8s contains commands for kubectl and other kubernetes related commands.
type K8s mg.Namespace

// Init copies the k8 yaml manifest files from the examples directory to the cache directory for editing and linking in integration testing.
func (K8s) Init() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultHeader.Println("(K8s) Init()")
	// Create the cache directory if it doesn't exist.
	if _, err := os.Stat(constants.CacheManifestDirectory); os.IsNotExist(err) {
		if err := os.MkdirAll(constants.CacheManifestDirectory, constants.PermissionUserReadWriteExecute); err != nil {
			return fmt.Errorf("os.MkdirAll(): %w", err)
		}
	}
	// For each file in the examples directory, create a copy in the CacheManifestDirectory.
	de, err := os.ReadDir(constants.ExamplesDirectory)
	if err != nil {
		return err
	}
	for _, file := range de {
		originalFile := filepath.Join(constants.ExamplesDirectory, file.Name())
		targetFile := filepath.Join(constants.CacheManifestDirectory, file.Name())
		// If the file doesn't exist in the manifest directory, read it and copy it to the manifest directory.
		if _, err := os.Stat(targetFile); os.IsNotExist(err) {
			// Read the original file.
			original, err := os.ReadFile(originalFile)
			if err != nil {
				return fmt.Errorf("unable to read original file: %s, os.ReadFile(): %w", original, err)
			}
			// Create the new file from the contents of the original file.
			if err := os.WriteFile(targetFile, original, constants.PermissionUserReadWriteExecute); err != nil {
				return fmt.Errorf("unable to write new file: %s, os.WriteFile(): %w", targetFile, err)
			}
			pterm.Success.Printfln("copied starter example (edit and apply to use): %s", targetFile)
		}
	}
	pterm.Success.Println("(K8s) Init()")
	return nil
}

// Apply applies a kubernetes manifest.
func (K8s) Apply(manifest string) error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultHeader.Println("(K8s) Apply()")
	return sh.Run(
		"kubectl",
		"apply",
		"--kubeconfig", constants.Kubeconfig,
		"--context", constants.KindContextName,
		"--namespace", constants.KubectlNamespace,
		"--cluster", constants.KindContextName,
		"--wait=true",
		"--overwrite=true",
		"-f", manifest,
	)
}

// Apply applies a kubernetes manifest.
func (K8s) Delete(manifest string) {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultHeader.Println("(K8s) Delete()")
	if err := sh.Run(
		"kubectl",
		"delete",
		"--kubeconfig", constants.Kubeconfig,
		"--context", constants.KindContextName,
		"--namespace", constants.KubectlNamespace,
		"--cluster", constants.KindContextName,
		"-f", manifest,
	); err != nil {
		pterm.Warning.Printfln("(K8s) Delete() error [non-terminating]: %s", err)
	}
}

// Logs streams logs until canceled for the dsv syncing jobs, based on the label `dsv.delinea.com: syncer`.
func (K8s) Logs() error {
	magetoolsutils.CheckPtermDebug()
	if _, err := exec.LookPath("stat"); err != nil {
		pterm.Error.Printfln("install stern tool manually (see .devcontainer/Dockerfile for install command) to run this")
		return errors.New("stern tool not installed yet")
	}
	pterm.DefaultHeader.Println("(K8s) Logs()")
	pterm.Info.Printfln("if you run into log output issues, just try running:\n\n\t\tkubectl logs  --context %s --namespace %s  --selector 'dsv-filter-name in (dsv-syncer, dsv-injector)' --follow --prefix\n", constants.KindContextName, constants.KubectlNamespace)
	pterm.Info.Println("🔍 query without selector:\n\n\tstern --kubeconfig .cache/config --namespace dsv  --timestamps . ")
	pterm.Info.Println(
		"🔍 Manually run stern with the following:\n\n\t",
		"stern",
		"--namespace", constants.KubectlNamespace,
		"--timestamps",
		"--selector", "dsv-filter-name in (dsv-syncer, dsv-injector)",
	)

	pterm.Info.Println(
		"🔍 Manually run stern againt entire cluster with following:\n\n\t",
		"stern",
		"--all-namespaces",
		"--timestamps",
		".",
	)
	pterm.DefaultHeader.Println("kubectl output first")
	_ = sh.RunV("kubectl",
		"logs",
		"--kubeconfig", constants.Kubeconfig,
		"--context", constants.KindContextName,
		"--namespace", constants.KubectlNamespace,
		"--cluster", constants.KindContextName,
		"--selector", "dsv-filter-name in (dsv-syncer, dsv-injector)",
		// "--follow",
		"--since=5m",
		"--prefix",
	)
	pterm.DefaultHeader.Println("stern streaming output")
	return sh.RunV(
		"stern",
		"--namespace", constants.KubectlNamespace,
		"--timestamps",
		"--selector", "dsv-filter-name in (dsv-syncer, dsv-injector)",
	)
}

// 🔍 OutputSecret outputs the base64 decoded values for local minikube style testing.
func (K8s) OutputSecret() {
	magetoolsutils.CheckPtermDebug()
	for _, secretname := range []string{"user-domain-pass", "user-domain", "pass-domain"} {
		response, err := sh.Output(
			"kubectl",
			"--kubeconfig", constants.Kubeconfig,
			"--context", constants.KindContextName,
			"--namespace", constants.KubectlNamespace,
			"--cluster", constants.KindContextName,
			"get",
			"secret", secretname,
			//"-o", "jsonpath='{.data.password}'",
			"-o", `go-template={{.data.password}}`,
			"--ignore-not-found",
		)
		if err != nil {
			pterm.Warning.Printfln("not able to find this %q: %v", secretname, err)
		} else {
			cleanedoutput := strings.TrimSpace(response)
			b, err := base64.StdEncoding.DecodeString(cleanedoutput)
			if err != nil {
				pterm.Warning.Printfln("issue decoding string: %v", err)
				pterm.Debug.Printfln(
					"kubectl --kubeconfig %s --context %s --namespace %s --cluster %s get secret %s -o go-template='{{.data.password}}' --ignore-not-found",
					constants.Kubeconfig,
					constants.KindContextName,
					constants.KubectlNamespace,
					constants.KindContextName,
					secretname,
				)
			} else {
				pterm.Info.Printfln("🔑 [only for local testing] %q: %q", secretname, string(b)) // ♥️ nested if statements 😀
			}
		}
	}
}
