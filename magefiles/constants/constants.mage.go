package constants

// Since we are dealing with builds, having a constants file until using a config input makes it easy.

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build.

const (
	// ArtifactDirectory is a directory containing artifacts for the project and shouldn't be committed to source.
	ArtifactDirectory = ".artifacts"

	// PermissionUserReadWriteExecute is the permissions for the artifact directory.
	PermissionUserReadWriteExecute = 0o0700

	// ConfigsDirectory is where the credentials used for testing is placed.
	// This isn't committed in git.
	ConfigsDirectory = "configs"

	// CacheDirectory is where the cache for the project is placed, ie artifacts that don't need to be rebuilt often.
	CacheDirectory = ".cache"

	// ExamplesDirectory is the directory where the kubernetes manifests are stored.
	ExamplesDirectory = "examples"

	// CacheChartsDirectory is the directory where the cached helm values file is copied to.
	CacheChartsDirectory = ".cache/charts"

	// CacheCredentialFile is the path to the credential file for the project, which is cached locally.
	CacheCredentialFile = ".cache/credentials.json" //nolint:gosec // this is a test project and this directory is excluded from source
)

const (
	// KindClusterName is the name of the kind cluster.
	KindClusterName = "dsvtest"
	// KindClusterName is the name of the kind cluster.
	KindContextName = "dsvtest"
	// KubeconfigPath is the path to the kubeconfig file for this project, which is cached locally.
	Kubeconfig = ".cache/config"
	// KubectlNamespace is the namespace used for all kubectl commands, so that they don't operate in default or other namespace by accident.
	KubectlNamespace = "dsv"

	// DockerImageQualified is the qualified path of the image in Docker Hub.
	DockerImageQualified = "docker.io/delineaxpm/dsv-k8s"
	// DockerImageNameLocal is the name of the built image to run locally and load with minikube/kind.
	DockerImageNameLocal = "dsv-k8s"
)

const (

	// HelmTimeout is the timeout for helm commands using the CLI.
	HelmTimeout = "5m"
	// ChartsDirectory is the directory where the helm charts are placed, in sub directories.
	ChartsDirectory = "charts"
	// SternFilter is the filter for dsv-filter-name for streaming logs easily.
	SternFilter = "dsv-syncer, dsv-injector"
)

const (
	// MinikubeCPU is the CPU count for minikube.
	MinikubeCPU = "2"
	// MinikubeMemory is the memory for minikube.
	MinikubeMemory = "2048"
)
