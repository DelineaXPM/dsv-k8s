package constants

import (
	"path/filepath"
	"time"
)

type HelmCharts struct {
	ReleaseName string
	ChartPath   string
	Namespace   string
	// Values is the customized yaml file placed in the CacheDirectory to run install with.
	// Copy the values.yaml from the helm chart to start here.
	Values string
}

var HelmChartsList = []HelmCharts{
	{
		ReleaseName: "dsv-syncer",
		ChartPath:   filepath.Join(ChartsDirectory, "dsv-syncer"),
		Namespace:   "dsv",
		Values:      filepath.Join(CacheDirectory, "dsv-syncer", "values.yaml"),
	},
	{
		ReleaseName: "dsv-injector",
		ChartPath:   filepath.Join(ChartsDirectory, "dsv-injector"),
		Namespace:   "dsv",
		Values:      filepath.Join(CacheDirectory, "dsv-injector", "values.yaml"),
	},
}

// DefaultHelmTimeoutMinutes is the default timeout for helm commands.
var DefaultHelmTimeoutMinutes = time.Minute * 5

// CacheManifestDirectory is the directory where helm charts are cached for local tweaking.
// They are copied from the examples directory to allow editing without committing to source control.
var CacheManifestDirectory = filepath.Join(CacheDirectory, "manifests")
