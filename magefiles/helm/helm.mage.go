// Helm package provides render, installation, and other automation commands for Helm charts.
package helm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DelineaXPM/dsv-k8s/v2/magefiles/constants"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	"github.com/pterm/pterm"
	mtu "github.com/sheldonhull/magetools/pkg/magetoolsutils"
	"github.com/sheldonhull/magetools/pkg/req"

	helmclient "github.com/mittwald/go-helm-client"
)

type Helm mg.Namespace

// ⚙️ Init sets up the required files to allow for local editing/overriding from CacheDirectory.
//
// It does this by using the HelmChartlist and copying the default values.yaml to the CacheDirectory.
func (Helm) Init() error {
	pterm.DefaultSection.Println("(Helm) Init()")
	for _, chart := range constants.HelmChartsList {
		pterm.DefaultSection.Printfln(
			"Copy values.yaml for: %s to CacheDirectory: %s",
			chart.ReleaseName,
			constants.CacheDirectory,
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
		targetFile := filepath.Join(constants.CacheDirectory, ln)
		targetDir, _ := filepath.Split(filepath.Join(constants.CacheDirectory, ln))
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

// newClient returns a helm client, and allows passing a kubeconfig path.
// By default it just uses what's set in the environment.
// If kubeconfig is provided then by default this should be using the local kind cluster setup.
//nolint:unparam // Allow initially. I placed in case I want to allow running against a target instance other than kind to allow overriding. Add an env check for KUBECONFIG and allow override then.
func newClient( //nolint:ireturn // Ignore for this helm project
	namespace, kubeconfig string,
) (helmclient.Client, error) {
	if kubeconfig == "" {
		opt := &helmclient.Options{
			Namespace: namespace,
			Debug:     true,
			Linting:   true,

			DebugLog: func(format string, v ...interface{}) {
				pterm.Debug.Printfln("[helmclient] "+format, v...)
				// Change this to your own logger. Default is 'log.Printf(format, v...)'.
			},
			// DebugLog:         func(format string, v ...interface{}) {},.
		}

		helmClient, err := helmclient.New(opt)
		if err != nil {
			return nil, fmt.Errorf("helm client failed: %w", err)
		}
		return helmClient, nil
	}
	if kubeconfig != "" {
		if _, err := os.Stat(constants.Kubeconfig); os.IsNotExist(err) {
			return nil, fmt.Errorf("kubeconfig file does not exist: %w", err)
		}
		// Read kubeconfig.
		kubeconfigData, err := os.ReadFile(constants.Kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to read kubeconfig file: %w", err)
		}

		opt := &helmclient.KubeConfClientOptions{
			Options: &helmclient.Options{
				Namespace: namespace,
				Debug:     true,
				Linting:   true,
				DebugLog: func(format string, v ...interface{}) {
					pterm.Debug.Printfln("[helmclient-kubeconfig] "+format, v...)
					// Change this to your own logger. Default is 'log.Printf(format, v...)'.
				},
			},
			KubeContext: constants.KindContextName,
			KubeConfig:  kubeconfigData,
		}

		helmClient, err := helmclient.NewClientFromKubeConf(opt)
		if err != nil {
			return nil, fmt.Errorf("helm client using kubeconfig failed: %w", err)
		}
		return helmClient, nil
	}
	return nil, nil
}

// getValuesYaml returns the values file as as string for consuming.
func getValuesYaml(helmValuesFile string) (string, error) {
	mtu.CheckPtermDebug()
	if _, err := os.Stat(helmValuesFile); os.IsNotExist(err) {
		pterm.Warning.Printfln(
			"%s file does not exist\n\nMake sure to run the following task to get working copies of the template files\tmage helm:init",
			helmValuesFile,
		)
		return "", fmt.Errorf(
			"getValuesYaml: %s file does not exist, probably need to run helm:init first ",
			helmValuesFile,
		)
	}

	// Since files don't exist, we should go ahead and copy the helm template files to get this started.
	vf, err := os.ReadFile(helmValuesFile)
	if err != nil {
		return "", fmt.Errorf("(Helm) Render() - Error reading values file: %w", err)
	}
	return string(vf), nil
}

// 💾 Render uses Helm to output rendered yaml for testing helm integration.
func (Helm) Render() error {
	mtu.CheckPtermDebug()
	pterm.DefaultHeader.Println("(Helm) Render()")
	for _, chart := range constants.HelmChartsList {
		pterm.DefaultSection.Printfln("New Client: %s", chart.ReleaseName)
		client, err := newClient(chart.Namespace, constants.Kubeconfig)
		if err != nil {
			return err
		}
		valuesYaml, err := getValuesYaml(chart.Values)
		if err != nil || len(valuesYaml) == 0 {
			return fmt.Errorf("failed to get values.yaml: %w", err)
		}
		chartSpec := helmclient.ChartSpec{
			ReleaseName:  chart.ReleaseName,
			ChartName:    chart.ChartPath,
			Namespace:    chart.Namespace,
			UpgradeCRDs:  true,
			Wait:         true,
			ValuesYaml:   valuesYaml,
			GenerateName: true,
			DryRun:       true,
		}

		out, err := client.TemplateChart(&chartSpec)
		if err != nil {
			return fmt.Errorf("helm template failed: %w", err)
		}

		outFile := filepath.Join(constants.ArtifactDirectory, chart.ReleaseName+".yml")
		if err := os.WriteFile(outFile, out, constants.PermissionUserReadWriteExecute); err != nil {
			return fmt.Errorf("failed to write templated output for helm chart: %w", err)
		}
		if err := sh.Run("yamlfmt", outFile, "-w"); err != nil {
			pterm.Warning.Printfln("yamlfmt failed: %s", err)
		}
		pterm.Success.Printfln("(Helm) Render() - Successfully rendered chart to %s", outFile)
	}
	return nil
}

// 🔍 Lint uses Helm to lint the chart for issues.
func (Helm) Lint() error {
	mtu.CheckPtermDebug()
	pterm.DefaultHeader.Println("(Helm) Lint()")
	for _, chart := range constants.HelmChartsList {
		pterm.DefaultSection.Printfln("New Client: %s", chart.ReleaseName)
		client, err := newClient(chart.Namespace, constants.Kubeconfig)
		if err != nil {
			return err
		}
		valuesYaml, err := getValuesYaml(chart.Values)
		if err != nil {
			return err
		}
		chartSpec := helmclient.ChartSpec{
			ReleaseName: chart.ReleaseName,
			ChartName:   chart.ChartPath,
			Namespace:   chart.Namespace,
			UpgradeCRDs: true,
			Wait:        true,
			ValuesYaml:  valuesYaml,
			DryRun:      true,
		}

		err = client.LintChart(&chartSpec)
		if err != nil {
			pterm.Warning.Printfln("(Helm) Lint() - Error linting chart: %s", err)
			return fmt.Errorf("helm lint failed: %w", err)
		}
	}
	return nil
}

// 🚀 Install uses Helm to install the chart.
func (Helm) Install() error {
	mtu.CheckPtermDebug()
	pterm.DefaultHeader.Println("(Helm) Install()")

	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultHelmTimeoutMinutes)
	defer cancel()
	for _, chart := range constants.HelmChartsList {
		pterm.DefaultSection.Printfln("New Client: %s", chart.ReleaseName)
		client, err := newClient(chart.Namespace, constants.Kubeconfig)
		if err != nil {
			return err
		}
		valuesYaml, err := getValuesYaml(chart.Values)
		if err != nil {
			return err
		}
		chartSpec := helmclient.ChartSpec{
			ReleaseName:      chart.ReleaseName,
			ChartName:        chart.ChartPath,
			Namespace:        chart.Namespace,
			Wait:             true,
			ValuesYaml:       valuesYaml,
			DryRun:           false,
			CreateNamespace:  true,
			Atomic:           true,
			Replace:          true,
			CleanupOnFail:    true,
			Force:            true,
			WaitForJobs:      true,
			Recreate:         true,
			DependencyUpdate: true,
			Timeout:          constants.DefaultHelmTimeoutMinutes,
		}

		rel, err := client.InstallOrUpgradeChart(ctx, &chartSpec, &helmclient.GenericHelmOptions{})
		if err != nil {
			return fmt.Errorf("helm install failed: %w", err)
		}
		pterm.Info.Printfln("releasestatus: %v", rel.Info.Status)
	}
	return nil
}

// 🚀 Uninstall uses Helm to uninstall the chart.
func (Helm) Uninstall() error {
	mtu.CheckPtermDebug()
	pterm.DefaultHeader.Println("(Helm) Uninstall()")

	for _, chart := range constants.HelmChartsList {
		pterm.DefaultSection.Printfln("New Client: %s", chart.ReleaseName)
		client, err := newClient(chart.Namespace, constants.Kubeconfig)
		if err != nil {
			return err
		}
		err = client.UninstallReleaseByName(chart.ReleaseName)
		if err != nil {
			pterm.Warning.Printfln("helm uninstall failed: %v", err)
			continue
		}
		pterm.Success.Printfln("uninstall of %q successful", chart.ReleaseName)
	}
	pterm.Success.Println("(Helm) Uninstall() - Successfully uninstalled charts")
	return nil
}

// Docs generates helm documentation using `helm-doc` tool.
func (Helm) Docs() error {
	binary, err := req.ResolveBinaryByInstall("helm-docs", "github.com/norwoodj/helm-docs/cmd/helm-docs@latest")
	if err != nil {
		return err
	}
	for _, chart := range constants.HelmChartsList {
		pterm.DefaultSection.Printfln("Generating docs for %s", chart.ReleaseName)
		err := sh.Run(binary,
			"--chart-search-root", chart.ChartPath,
			"--output-file", "README.md",
		)
		if err != nil {
			return fmt.Errorf("helm-docs failed: %w", err)
		}
		pterm.Success.Printfln("generated file: %s", filepath.Join(chart.ChartPath, "README.md"))
	}
	pterm.Success.Println("(Helm) Docs() - Successfully generated readmes for charts")

	return nil
}
