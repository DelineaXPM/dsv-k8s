// K8s contains commands for kubectl and other kubernetes related commands.
package k8s

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/constants"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// k8s contains commands for kubectl and other kubernetes related commands.
type K8s mg.Namespace

// Init copies the k8 yaml manifest files from the examples directory to the cache directory for editing and linking in integration testing.
func (K8s) Init() error {
	mtu.CheckPtermDebug()
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
	mtu.CheckPtermDebug()
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
	mtu.CheckPtermDebug()
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
	mtu.CheckPtermDebug()
	if _, err := exec.LookPath("stat"); err != nil {
		pterm.Error.Printfln("install stern tool manually (see .devcontainer/Dockerfile for install command) to run this")
		return errors.New("stern tool not installed yet")
	}
	pterm.DefaultHeader.Println("(K8s) Logs()")
	pterm.Warning.Printfln("if you run into log output issues, just try running:\n\n\t\tkubectl logs  --context %s --namespace %s  --selector 'dsv-filter-name in (dsv-syncer, dsv-injector)' --follow --prefix\n", constants.KindContextName, constants.KubectlNamespace)
	pterm.Warning.Println("query without selector:\n\n\tstern --kubeconfig .cache/config --namespace dsv  --timestamps . ")
	pterm.Debug.Println(
		"Manually run stern with the following:\n\n\t",
		"stern",
		"--namespace", constants.KubectlNamespace,
		"--timestamps",
		"--selector", "dsv-filter-name in (dsv-syncer, dsv-injector)",
	)
	return sh.RunV(
		"stern",
		"--namespace", constants.KubectlNamespace,
		"--timestamps",
		"--selector", "dsv-filter-name in (dsv-syncer, dsv-injector)",
	)
}
