package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
	"github.com/sheldonhull/magetools/pkg/req"
)

func checkEnvVar(envVar string, required bool) (string, error) {
	envVarValue := os.Getenv(envVar)
	if envVarValue == "" && required {
		pterm.Error.Printfln(
			"%s is required and unable to proceed without this being provided. terminating task.",
			envVar,
		)
		return "", fmt.Errorf("%s is required", envVar)
	}
	if envVarValue == "" {
		pterm.Debug.Printfln(
			"checkEnvVar() found no value for: %q, however this is marked as optional, so not exiting task",
			envVar,
		)
	}
	pterm.Debug.Printfln("checkEnvVar() found value: %q=%q", envVar, envVarValue)
	return envVarValue, nil
}

// 🔨 Build builds the project for the current platform.
func Build() error {
	magetoolsutils.CheckPtermDebug()
	binary, err := req.ResolveBinaryByInstall("goreleaser", "github.com/goreleaser/goreleaser@latest")
	if err != nil {
		return err
	}

	releaserArgs := []string{
		"build",
		"--clean",
		"--snapshot",
		"--single-target",
	}
	pterm.Debug.Printfln("goreleaser: %+v", releaserArgs)

	return sh.RunV(binary, releaserArgs...) // "--skip-announce",.
}

// 🔨 BuildAll builds all the binaries defined in the project, for all platforms. This includes Docker image generation but skips publish.
// If there is no additional platforms configured in the task, then basically this will just be the same as `mage build`.
func BuildAll() error {
	magetoolsutils.CheckPtermDebug()
	binary, err := req.ResolveBinaryByInstall("goreleaser", "github.com/goreleaser/goreleaser@latest")
	if err != nil {
		return err
	}

	releaserArgs := []string{
		"release",
		"--snapshot",
		"--clean",
		"--skip-publish",
	}
	pterm.Debug.Printfln("goreleaser: %+v", releaserArgs)
	return sh.RunV(binary, releaserArgs...)
	// To pass in explicit version mapping, you can do this. I'm not using at this time.
	// Return sh.RunWithV(map[string]string{
	// 	"GORELEASER_CURRENT_TAG": "latest",
	// }, binary, releaserArgs...)
}

// 🔨 Release generates a release for the current platform.
func Release() error {
	magetoolsutils.CheckPtermDebug()
	binary, err := req.ResolveBinaryByInstall("goreleaser", "github.com/goreleaser/goreleaser@latest")
	if err != nil {
		return err
	}

	if _, err = checkEnvVar("DOCKER_ORG", true); err != nil {
		return err
	}

	changieBinary, err := req.ResolveBinaryByInstall("changie", "github.com/miniscruff/changie@latest")
	if err != nil {
		pterm.Error.Println("unable to install changelog binary")
		return err
	}
	releaseVersion, err := sh.Output(changieBinary, "latest")
	if err != nil {
		pterm.Warning.Printfln("changie pulling latest release note version failure: %v", err)
	}
	cleanVersion := strings.TrimSpace(releaseVersion)
	cleanpath := filepath.Join(".changes", cleanVersion+".md")
	if os.Getenv("GITHUB_WORKSPACE") != "" {
		cleanpath = filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".changes", cleanVersion+".md")
	}

	releaserArgs := []string{
		"release",
		"--clean",
		"--skip-validate",
		fmt.Sprintf("--release-notes=%s", cleanpath),
	}
	pterm.Debug.Printfln("goreleaser: %+v", releaserArgs)

	return sh.RunWithV(map[string]string{
		"GORELEASER_CURRENT_TAG": cleanVersion,
	},
		binary,
		releaserArgs...,
	)
}
