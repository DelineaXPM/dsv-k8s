// ⚡ Core Mage Tasks.
package main

import (
	"os"

	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/constants"
	// mage:import
	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/helm"
	// mage:import
	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/k8s"
	// mage:import
	_ "github.com/DelineaXPM/dsv-k8s/v2/magefiles/kind"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/ci"
	"github.com/sheldonhull/magetools/fancy"
	"github.com/sheldonhull/magetools/tooling"

	// mage:import
	"github.com/sheldonhull/magetools/gittools"
	// mage:import
	"github.com/sheldonhull/magetools/gotools"
	// mage:import
	"github.com/sheldonhull/magetools/precommit"
	//mage:import
	_ "github.com/sheldonhull/magetools/secrets"
)

// createDirectories creates the local working directories for build artifacts and tooling.
func createDirectories() error {
	for _, dir := range []string{constants.ArtifactDirectory, constants.CacheDirectory, constants.ConfigsDirectory} {
		if err := os.MkdirAll(dir, constants.PermissionUserReadWriteExecute); err != nil {
			pterm.Error.Printf("failed to create dir: [%s] with error: %v\n", dir, err)

			return err
		}
		pterm.Success.Printf("✅ [%s] dir created\n", dir)
	}

	return nil
}

// Init runs multiple tasks to initialize all the requirements for running a project for a new contributor.
func Init() error { //nolint:deadcode // Not dead, it's alive.
	fancy.IntroScreen(ci.IsCI())
	pterm.Success.Println("running Init()...")

	mg.SerialDeps(
		Clean,
		createDirectories,
		(gotools.Go{}.Tidy),
		(gotools.Go{}.Init),
	)
	// These can run in parallel as different toolchains.
	mg.Deps(
		(gittools.Gittools{}.Init),
		(precommit.Precommit{}.Init),
	)

	if ci.IsCI() {
		pterm.Debug.Println("CI detected, done with init")
		return nil
	}

	pterm.DefaultSection.Println("Setup Project Specific Tools")
	if err := tooling.SilentInstallTools(toolList); err != nil {
		return err
	}

	if err := sh.Run("docker", "pull", "alpine:latest"); err != nil {
		return err
	}

	mg.Deps(
		k8s.K8s{}.Init,
		helm.Helm{}.Init,
	)
	return nil
}

// Clean up after yourself.
func Clean() {
	pterm.Success.Println("Cleaning...")
	for _, dir := range []string{constants.ArtifactDirectory, constants.CacheDirectory, constants.ConfigsDirectory} {
		err := os.RemoveAll(dir)
		if err != nil {
			pterm.Error.Printf("failed to removeall: [%s] with error: %v\n", dir, err)
		}
		pterm.Success.Printf("🧹 [%s] dir removed\n", dir)
	}
	mg.Deps(createDirectories)
}
