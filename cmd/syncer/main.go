package main

import (
	"fmt"
	"os"
	"time"

	"github.com/DelineaXPM/dsv-k8s/v2/internal/logger"
	"github.com/DelineaXPM/dsv-k8s/v2/pkg/config"
	"github.com/DelineaXPM/dsv-k8s/v2/pkg/syncer"
	"github.com/rs/zerolog"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	env "github.com/caarlos0/env/v6"
)

const (
	// ExitFailure is exit code sent for failed task.
	exitFailure = 1
	// ExitSuccess is exit code sent for running without any error.
	exitSuccess = 0
)

// main is the entrypoint for the syncer; parses the credentials file and calls syncer.Sync
func main() {
	if err := Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFailure)
	}
	os.Exit(exitSuccess) // shouldn't hit this if run is invoked correctly
}

//nolint:gochecknoglobals // ok for providing as version output
var (
	// Version is the descriptive version, normally the tag from which the app was built.
	// Since git tags can be changed, use Commit instead as the most accurate version.
	version = "dev"
	// Commit is the git commit hash that the build was generated from.
	commit = "none"
	// Date is the date the binary was produced.
	date = "unknown"
	// buildName is the build name for easier confirmation on local builds that a build has changed.
	buildName = "unknown"
)

// Run contains the actual invocation code for the syncer and is public to allow running integration tests with it.
func Run(args []string) error { //nolint:funlen // ok for Run
	log := logger.New()
	log.Info().
		Str("version", version).
		Str("commit", commit).
		Str("date", date).
		Str("buildName", buildName).
		Msg("syncer version information")

	// Config is the configuration for the syncer.
	// This is provided by environment variables.
	type Config struct {
		Namespace           string `env:"DSV_NAMESPACE"`                                                                      // DSV_NAMESPACE is the namespace for secrets to sync. "" (the default) by default includes all namespaces.
		Debug               bool   `env:"DSV_DEBUG" envDefault:"false"`                                                       // Debug enables debug logging.
		KubeConfig          string `env:"KUBECONFIG" envDefault:"${HOME}/.kube/config" envExpand:"true"`                      // KubeConfig is the path to the kubeconfig file and only required if running locally. By default the connection will be built as "incluster".
		CredentialsJSONFile string `env:"DSV_CREDENTIALS_JSON" envDefault:"${HOME}/credentials/config.json" envExpand:"true"` // CredentialsJSONFile is the path to the JSON formatted credentials file that is mounted as a secret.
	}

	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse environment variables")
	}
	log.Info().Strs("args", args).Msg("starting syncer, args passed, but not used, as environment variables are used instead")
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("debug logging enabled")
	}

	credentials, err := config.GetCredentials(cfg.CredentialsJSONFile)
	if err != nil {
		log.Fatal().Msgf("[ERROR] unable to process configuration file '%s': %s", cfg.CredentialsJSONFile, err)
		return fmt.Errorf("unable to process configuration file %q: %w", cfg.CredentialsJSONFile, err)
	}
	log.Info().
		Strs("credential_list", credentials.Names()).
		Int("credential_count", len(*credentials)).
		Msg("success loading credential sets")

	start := time.Now()
	if err := syncer.Sync(
		func(namespace string) (rconfig *rest.Config, err error) {
			rconfig, err = rest.InClusterConfig()
			if err != nil {
				log.Debug().Msg("unable to get InClusterConfig, falling back to KubeConfig")

				rconfig, err = clientcmd.BuildConfigFromFlags("", cfg.KubeConfig)
				if err != nil {
					log.Error().Err(err).Msg("error getting Kubernetes Client rest.Config")
					return nil, fmt.Errorf("error getting Kubernetes Client rest.Config: %w", err)
				}
				return rconfig, nil
			}
			return rconfig, nil
		},
		cfg.Namespace,
		*credentials,
		log,
	); err != nil {
		log.Fatal().Msgf("[ERROR] unable to sync Secrets: %s", err)
	}
	log.Info().Dur("duration", time.Since(start)).Msg("syncer processing complete")
	return nil
}
