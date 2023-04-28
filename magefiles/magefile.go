// ⚡ Core Mage Tasks.
package main

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/constants"
	"github.com/bitfield/script"

	// mage:import
	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/helm"
	// mage:import
	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/k8s"
	// mage:import
	_ "github.com/DelineaXPM/dsv-k8s/v2/magefiles/kind"
	// mage:import
	_ "github.com/DelineaXPM/dsv-k8s/v2/magefiles/minikube"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/ci"
	"github.com/sheldonhull/magetools/fancy"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"

	// mage:import
	"github.com/sheldonhull/magetools/gotools"
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
	)

	_, err := exec.LookPath("aqua")
	if err != nil && os.IsNotExist(err) {
		pterm.Error.Printfln("unable to resolve aqua cli tool, please install for automated project tooling setup: https://aquaproj.github.io/docs/tutorial-basics/quick-start#install-aqua")
		return err
	}

	if ci.IsCI() {
		installArgs := []string{}

		if mg.Verbose() {
			installArgs = append(installArgs, "--log-level")
			installArgs = append(installArgs, "debug")
		}
		installArgs = append(installArgs, "install")
		installArgs = append(installArgs, "aqua")
		if err := sh.RunV("aqua", installArgs...); err != nil {
			pterm.Error.Printfln("aqua-ci%v", err)
			return err
		}
		pterm.Debug.Println("CI detected, done with init")
		return nil
	}

	pterm.DefaultSection.Println("Aqua install (any first packages)")
	if err := sh.RunV("aqua", "install", "--tags", "first"); err != nil {
		pterm.Warning.Printfln("aqua install failed, continuing: %v", err)
	}
	pterm.Success.Println("aqua install --tags first complete")
	pterm.DefaultSection.Println("Aqua install remaining tools")
	if err := sh.RunV("aqua", "install"); err != nil {
		pterm.Warning.Printfln("aqua install failed, continuing: %v", err)
	}
	pterm.Success.Println("aqua install complete")
	// These can run in parallel as different toolchains.
	mg.SerialDeps(
		(InstallTrunk),
		(TrunkInit),
	)
	pterm.Info.Printfln("Initializing .cache/ directory with copies of Kubernetes YAML + Helm charts.\nUse this to edit your local configurations for minikube based testing")
	mg.Deps(
		k8s.K8s{}.Init,
		helm.Helm{}.Init,
	)
	return nil
}

// InstallTrunk installs trunk.io tooling if it isn't already found.
func InstallTrunk() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("InstallTrunk()")
	if runtime.GOOS == "windows" {
		pterm.Warning.Println("InstallTrunk() trunk.io not supported on windows, skipping")
		return nil
	}
	_, err := exec.LookPath("trunk")
	if err != nil {
		// if os.IsNotExist(err) {
		pterm.Warning.Printfln("unable to resolve aqua cli tool, please install for automated project tooling setup: https://aquaproj.github.io/docs/tutorial-basics/quick-start#install-aqua")
		_, err := script.Exec("curl https://get.trunk.io -fsSL").Exec("bash -s -- -y").Stdout()
		if err != nil {
			return err
		}
		// }
	} else {
		pterm.Success.Printfln("trunk.io already installed, skipping")
	}
	return nil
}

// TrunkInit ensures the required runtimes are installed.
func TrunkInit() error {
	return sh.RunV("trunk", "install")
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
