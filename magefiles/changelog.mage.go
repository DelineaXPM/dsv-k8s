package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
)

// Changelog contains the changelog tasks that batch up the changelog commands and allow triggering a release with an explicit version.
type Changelog mg.Namespace

// getVersion returns the version and path for the changefile to use for the semver and release notes.
func getVersion() (releaseVersion, cleanPath string, err error) { //nolint:unparam // leaving as optional parameter for future release tasks.
	releaseVersion, err = sh.Output("changie", "latest")
	if err != nil {
		pterm.Error.Printfln("changie pulling latest release note version failure: %v", err)
		return "", "", err
	}
	cleanVersion := strings.TrimSpace(releaseVersion)
	cleanPath = filepath.Join(".changes", cleanVersion+".md")
	if os.Getenv("GITHUB_WORKSPACE") != "" {
		cleanPath = filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".changes", cleanVersion+".md")
	}
	return cleanVersion, cleanPath, nil
}

// 📦 Bump the application as an interactive command, prompting for semver change type, merging changelog, and running format and git add.
func (Changelog) Bump() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("(Changelog) Bump()")
	bumpType, _ := pterm.DefaultInteractiveSelect.
		WithOptions([]string{"patch", "minor", "major"}).
		Show()
	pterm.Info.Printfln("bumping by: %s", bumpType)
	if bumpType == "major" {
		pterm.Warning.Printfln("🔥major bumping should be done with care as this signifies large breaking changes")
		pterm.Warning.Println("pick the value against to signify you get this 😀")
		bumpType, _ = pterm.DefaultInteractiveSelect.
			WithOptions([]string{"patch", "minor", "major"}).
			Show()
	}
	if err := sh.RunV("changie", "batch", bumpType); err != nil {
		pterm.Warning.Printf("changie batch failure (non-terminating as might be repeating batch command): %v", err)
	}
	if err := sh.RunV("changie", "merge"); err != nil {
		return err
	}
	if err := sh.RunV("trunk", "fmt"); err != nil {
		return err
	}
	if err := sh.RunV("trunk", "check", "--ci"); err != nil {
		pterm.Warning.Printfln("trunk check failure. This is non-terminating for the mage task, but you should check it before merging")
	}
	if err := sh.RunV("git", "add", ".changes/*"); err != nil {
		return err
	}
	if err := sh.RunV("git", "add", "charts/**/Chart.yaml"); err != nil {
		return err
	}
	if err := sh.RunV("git", "add", "CHANGELOG.md"); err != nil {
		return err
	}

	releaseVersion, _, err := getVersion()
	if err != nil {
		return err
	}
	pterm.Info.Println(" Are you ready to create a commit with these changes?")
	confirm, err := pterm.DefaultInteractiveConfirm.
		WithDefaultValue(false).
		WithRejectText("no").
		WithConfirmText("yes").
		WithDefaultValue(false).Show()
	if err != nil {
		return err
	}
	if !confirm {
		pterm.Warning.Println("someone changed their mind")
		return nil
	}
	if err := sh.RunV("git", "commit", "-m", fmt.Sprintf("feat: 🚀 create release %s", releaseVersion)); err != nil {
		return err
	}
	return nil
}

// 📦 Merge updates the changelog without bumping the version.
// This is useful for when you are picking up after the changie batch has already completed, but need to re-run the changie merge.
func (Changelog) Merge() error {
	magetoolsutils.CheckPtermDebug()
	pterm.DefaultSection.Println("(Changelog) Merge()")
	if err := sh.RunV("changie", "merge"); err != nil {
		return err
	}
	if err := sh.RunV("trunk", "fmt"); err != nil {
		return err
	}
	if err := sh.RunV("trunk", "check", "--ci"); err != nil {
		pterm.Warning.Printfln("trunk check failure. This is non-terminating for the mage task, but you should check it before merging")
	}
	if err := sh.RunV("git", "add", ".changes/*"); err != nil {
		return err
	}
	if err := sh.RunV("git", "add", "CHANGELOG.md"); err != nil {
		return err
	}
	releaseVersion, _, err := getVersion()
	if err != nil {
		return err
	}
	pterm.Info.Println(" Are you ready to create a commit with these changes?")
	confirm, err := pterm.DefaultInteractiveConfirm.
		WithDefaultValue(false).
		WithRejectText("no").
		WithConfirmText("yes").
		WithDefaultValue(false).Show()
	if err != nil {
		return err
	}
	if !confirm {
		pterm.Warning.Println("someone changed their mind")
		return nil
	}
	if err := sh.RunV("git", "commit", "-m", fmt.Sprintf("feat: 🚀 create release %s", releaseVersion)); err != nil {
		return err
	}
	return nil
}
