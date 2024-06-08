// Helm package provides render, installation, and other automation commands for Helm charts.
package helm

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/constants"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/pterm/pterm"
	"github.com/sheldonhull/magetools/pkg/magetoolsutils"
	"github.com/sheldonhull/magetools/pkg/req"
)

// Helm contains Mage tasks for invoking Helm cli.
type Helm mg.Namespace

// invokeHelm is a wrapper for running the helm binary.
func invokeHelm(args ...string) error {
	binary, err := req.ResolveBinaryByInstall("helm", "helm.sh/helm/v3@latest")
	if err != nil {
		return err
	}
	return sh.Run(binary, args...)
}

// ⚙️ Init sets up the required files to allow for local editing/overriding from CacheDirectory.
//
// It does this by using the HelmChartlist and copying the default values.yaml to the CacheDirectory.
func (Helm) Init() error {
	pterm.DefaultSection.Println("(Helm) Init()")
	magetoolsutils.CheckPtermDebug()
	for _, chart := range constants.HelmChartsList {
		pterm.DefaultSection.Printfln(
			"Copy values.yaml for: %s to CacheDirectory: %s",
			chart.ReleaseName,
			constants.CacheChartsDirectory,
		)

		// FileToCopy is the helm Values file without the parent directory path.
		fileToCopyList := strings.Split(chart.Values, string(filepath.Separator))
		pterm.Debug.Printfln("[getValuesYaml] fileToCopyList: %+v", fileToCopyList)
		if len(fileToCopyList) <= 1 {
			return fmt.Errorf("failed to get file name from %+v", fileToCopyList)
		}
		ln := filepath.Join(fileToCopyList[1:]...)
		fileToCopy := filepath.Join(constants.ChartsDirectory, ln)

		// Since the file doesn't exist let's read the contents and update an equivalent in the CacheDirectory for local editing and tweaking.
		targetFile := filepath.Join(constants.CacheChartsDirectory, ln)
		targetDir, _ := filepath.Split(filepath.Join(constants.CacheChartsDirectory, ln))
		if _, err := os.Stat(targetFile); !os.IsNotExist(err) {
			pterm.Info.Printfln("file: %s already exists in target, bypassing", targetFile)
			continue
		}
		// CopyPlaceholderValueshelmValuesFile copies the placeholder values.yaml file to the cache directory to get started with customizing local tests.
		pterm.Info.Printfln(
			"Init() %s file doesn't exist, so copying original helm values.yaml file from charts directory to bootstrap this",
			fileToCopy,
		)

		// Create an equivalent path in the .cache directory to work with.
		if err := os.MkdirAll(targetDir, constants.PermissionUserReadWriteExecute); err != nil {
			pterm.Error.Printfln("unable to create the required directory: %v", err)
		}
		data, err := os.ReadFile(fileToCopy)
		pterm.Info.Printfln("copying: %s to targetDir: %s", fileToCopy, targetDir)
		if err != nil {
			pterm.Error.Printfln("unable to read the file: %v", err)
		}
		if err := os.WriteFile(targetFile, data, constants.PermissionUserReadWriteExecute); err != nil {
			pterm.Error.Printfln("unable to write the file: %v", err)
		}
	}
	pterm.Success.Printfln("(Helm) Init()")
	return nil
}

// 🚀 Install uses Helm to
// 🚀 Install installs/upgrades the helm charts for charts listed in constants.HelmChartsList.
func (Helm) Install() error {
	magetoolsutils.CheckPtermDebug()
	if os.Getenv("KUBECONFIG") != constants.Kubeconfig {
		pterm.Warning.Printfln("KUBECONFIG is not set to %s. Make sure direnv/env variables loading if you want to keep the project changes from changing your user KUBECONFIG.", constants.Kubeconfig)
	}
	for _, chart := range constants.HelmChartsList {
		pterm.Info.Printfln("Installing chart: %s", chart.ReleaseName)
		sourceValuesFile := filepath.Join(constants.CacheChartsDirectory, chart.ReleaseName, "values.yaml")
		if err := Checkfile(sourceValuesFile); err != nil {
			pterm.Warning.Printfln("validation of values file threw some errors, might run into issues if this wasn't expected")
		}
		debugHelm := "--debug=false"
		if mg.Verbose() {
			pterm.Debug.Println("debug flag enabled for helm")
			debugHelm = "--debug=true" // enable verbose output
		}
		if _, err := os.Stat(constants.CacheCredentialFile); os.IsNotExist(err) {
			return fmt.Errorf("credentials file: %s doesn't exist, so skipping", constants.CacheCredentialFile)
		}
		if err := invokeHelm("upgrade",
			chart.ReleaseName,
			chart.ChartPath,
			"--namespace", constants.KubectlNamespace,
			"--install", // install if not already installed
			"--atomic",  // if set, the installation process deletes the installation on failure. The --wait flag will be set automatically if --atomic is used
			// "--replace", // re-use the given name, only if that name is a deleted release which remains in the history. This is unsafe in production
			"--wait", // waits, those atomic already runs this
			"--values", sourceValuesFile,
			"--timeout", constants.HelmTimeout,
			"--force",             // force resource updates through a replacement strategy
			"--wait-for-jobs",     // will wait until all Jobs have been completed before marking the release as successful
			"--dependency-update", // update dependencies if they are missing before installing the chart
			"--set-file", fmt.Sprintf("credentialsJson=%s", constants.CacheCredentialFile),

			debugHelm,
			// NOTE: Can pass credentials/certs etc in. NOT ADDED YET - "--set-file", "sidecar.configFile=config.yaml",
		); err != nil {
			pterm.Warning.Printfln("failed to install chart: %s, err: %v", chart.ReleaseName, err)
		} else {
			pterm.Success.Printfln("successfully installed chart: %s", chart.ReleaseName)
		}
	}
	return nil
}

// Uninstall uninstalls all the charts listed in constants.HelmChartsList.
func (Helm) Uninstall() {
	magetoolsutils.CheckPtermDebug()
	if os.Getenv("KUBECONFIG") != ".cache/config" {
		pterm.Warning.Printfln("KUBECONFIG is not set to %s. Make sure direnv/env variables loading if you want to keep the project changes from changing your user KUBECONFIG.", constants.Kubeconfig)
	}
	for _, chart := range constants.HelmChartsList {
		pterm.Info.Printfln("Uninstalling: %s", chart.ReleaseName)
		if err := invokeHelm("uninstall",
			chart.ReleaseName,
			"--wait",  // waits, those atomic already runs this
			"--debug", // enable verbose output
		); err != nil {
			pterm.Warning.Printfln("failed to uninstall: %s, err: %v", chart.ReleaseName, err)
		} else {
			pterm.Success.Printfln("Successfully uninstalled: %s", chart.ReleaseName)
		}
	}
}

// 💾 Render outputs the Kubernetes manifests from the helm template for debugging purposes.
func (Helm) Render() {
	magetoolsutils.CheckPtermDebug()
	if os.Getenv("KUBECONFIG") != constants.Kubeconfig {
		pterm.Warning.Printfln("KUBECONFIG is not set to .cache/config. Make sure direnv/env variables loading if you want to keep the project changes from changing your user KUBECONFIG.")
	}
	for _, chart := range constants.HelmChartsList {
		pterm.Info.Printfln("Rendering: %s", chart.ReleaseName)
		targetDirectory := filepath.Join(constants.CacheDirectory, "exported-template", chart.ReleaseName)
		_ = sh.Rm(targetDirectory) // no need to check for error, just clean directory if exists
		if err := os.MkdirAll(targetDirectory, constants.PermissionUserReadWriteExecute); err != nil {
			pterm.Error.Printfln("unable to create target chart directory for rendering helm template. what gives?")
			return
		}
		if err := invokeHelm("template",
			chart.ReleaseName,
			chart.ChartPath,
			"--values", filepath.Join(constants.CacheChartsDirectory, chart.ReleaseName, "values.yaml"),
			// "--create-namespace",
			// "--dependency-update",
			"--output-dir", targetDirectory,
		); err != nil {
			pterm.Warning.Printfln("failed to render template to: %s, err: %v", targetDirectory, err)
		} else {
			pterm.Success.Printfln("Successfully exported to targetDirectory: %s", targetDirectory)
		}
	}
}

// Docs generates helm documentation using `helm-doc` tool.
func (Helm) Docs() error {
	magetoolsutils.CheckPtermDebug()
	binary, err := req.ResolveBinaryByInstall("helm-docs", "github.com/norwoodj/helm-docs/cmd/helm-docs@latest")
	if err != nil {
		return err
	}
	for _, chart := range constants.HelmChartsList {
		pterm.DefaultSection.Printfln("Generating docs for %s", chart.ReleaseName)
		err := sh.Run(binary,
			"--chart-search-root", chart.ChartPath,
			"--output-file", "README.md",
			// NOTE: using default layout, but can change here if we wanted.
			// "--template-files", filepath.Join("magefiles", "helm", "README.md.gotmpl"),
		)
		if err != nil {
			return fmt.Errorf("helm-docs failed: %w", err)
		}
		pterm.Success.Printfln("generated file: %s", filepath.Join(chart.ChartPath, "README.md"))
	}
	pterm.Success.Println("(Helm) Docs() - Successfully generated readmes for charts")

	return nil
}

// checkfile confirms the cached values file is correctly set to allow loading images locally.
func Checkfile(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		pterm.Error.Printfln("❌ Error reading file %q %v", file, err)
		return err
	}
	escapedRepositoryMatch := regexp.QuoteMeta(constants.DockerImageNameLocal)
	re := regexp.MustCompile(fmt.Sprintf(`repository:\s+%s`, escapedRepositoryMatch)) //nolint:varnamelen // standard prefix, can update golangcilint config
	match := re.Find(b)
	if match != nil {
		pterm.Success.Printfln("✅ %s is configured to use local image  [expected: %s]", file, constants.DockerImageNameLocal)
	} else {
		re = regexp.MustCompile(`repository:\s+[^\n]*`)
		match = re.Find(b)
		if match != nil {
			pterm.Warning.Printfln("❌ %s: not configured to use local image: %q [expected: %s] (this is fine if you are't building as a developer with changes)", file, match, escapedRepositoryMatch)
		} else {
			pterm.Warning.Printfln("❌ %s: not configured to use local image: repository not found", file)
		}
	}

	re = regexp.MustCompile(`pullPolicy:\s+Never`)
	match = re.Find(b)
	if match != nil {
		pterm.Success.Printfln("✅ %s is configured with pullPolicy of Never [expected: Never]", file)
	} else {
		re = regexp.MustCompile(`pullPolicy:\s+\w*`)
		match = re.Find(b)
		if match != nil {
			pterm.Warning.Printfln("❌ %s: not configured with pullPolicy: Never: %q (this fine if you aren't building locally and just using docker image)", file, match)
		} else {
			pterm.Warning.Printfln("❌ %s:  not configured with pullPolicy: Never: pullPolicy not found", file)
		}
	}
	re = regexp.MustCompile(`tag:\s+[']?latest[']?`)
	match = re.Find(b)
	if match != nil {
		pterm.Success.Printfln("✅ %s is configured with tag of %s [expected: latest]", file, match)
	} else {
		re = regexp.MustCompile(`tag:\s+[']?.*[']?`)
		match = re.Find(b)
		if match != nil {
			pterm.Warning.Printfln("❌ %s: not configured with tag: latest (this is fine if using docker image in cloud) %q", file, match)
		} else {
			pterm.Warning.Printfln("❌ %s: not configured with tag: Never, tag not found", file)
		}
	}
	return nil
}
