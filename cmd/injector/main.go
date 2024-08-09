package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DelineaXPM/dsv-k8s/v2/pkg/config"
	"github.com/DelineaXPM/dsv-k8s/v2/pkg/injector"

	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caarlos0/env/v6"

	"github.com/DelineaXPM/dsv-k8s/v2/internal/logger"
)

const (
	// ExitFailure is exit code sent for failed task.
	exitFailure = 1
	// ExitSuccess is exit code sent for running without any error.
	exitSuccess = 0
)

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

// main is the entry point for the injector. It creates an HTTPS listener and listing for v1.AdmissionReview requests
func main() {
	if err := Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFailure)
	}
	os.Exit(exitSuccess) // shouldn't hit this if run is invoked correctly
}

// Run contains the actual invocation code for the injector and is public to allow running integration tests with it.
func Run(args []string) error { //nolint:funlen,cyclop // ok for Run
	log := logger.New()
	log.Info().
		Str("version", version).
		Str("commit", commit).
		Str("date", date).
		Str("buildName", buildName).
		Msg("injector version information")

	// Config is the configuration for the injector.
	// This is provided by environment variables.
	type Config struct {
		CertFile            string `env:"DSV_CERT"  envDefault:"${HOME}/tls/cert.pem" envExpand:"true"`                       // Cert is the path to the public certificate file in PEM format.
		KeyFile             string `env:"DSV_KEY" envDefault:"${HOME}/tls/key.pem" envExpand:"true"`                          // Key is the path to the private key file in PEM format.
		CredentialsJSONFile string `env:"DSV_CREDENTIALS_JSON" envDefault:"${HOME}/credentials/config.json" envExpand:"true"` // CredentialsJSONFile is the path to the JSON formatted credentials file that is mounted as a secret.
		ServerAddress       string `env:"DSV_SERVER_ADDRESS" envDefault:":18543"`                                             // ServerAddress is the address to listen on, e.g., 'localhost:8080' or ':8443'
		Debug               bool   `env:"DSV_DEBUG" envDefault:"false"`                                                       // Debug enables debug logging.
	}

	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse environment variables")
		return fmt.Errorf("fatal issue, as unabale to parse the required environment variables: %w", err)
	}
	if cfg.Debug {
		logger.EnableDebug()
		log.Info().Msg("debug logging enabled")
	}
	log.Info().Strs("args", args).Msg("starting injector, args passed, but not used, as environment variables are used instead")

	credentials, err := config.GetCredentials(cfg.CredentialsJSONFile)
	if err != nil {
		log.Error().Err(err).Str("credential-json", cfg.CredentialsJSONFile).Msg("unable to process credentials file")
		return fmt.Errorf("unable to process credentials file %q: %w", cfg.CredentialsJSONFile, err)
	}
	log.Info().
		Str("credential_names", strings.Join(credentials.Names(), ", ")).
		Str("credential_file", cfg.CredentialsJSONFile).
		Msg("credentials loaded from JSON file")
	var tlsConfig *tls.Config
	if cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile); err == nil {
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		log.Info().Str("cert", cfg.CertFile).Str("key", cfg.KeyFile).Msg("LoadX509KeyPair")

		// Parse the certificate to get the expiration date
		certData, err := os.ReadFile(cfg.CertFile)
		if err != nil {
			log.Error().Err(err).Msg("unable to read certificate file")
			return fmt.Errorf("unable to read certificate file: %w", err)
		}
		block, _ := pem.Decode(certData)
		if block == nil {
			log.Error().Msg("failed to parse certificate PEM")
			return fmt.Errorf("failed to parse certificate PEM")
		}
		parsedCert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse certificate")
			return fmt.Errorf("failed to parse certificate: %w", err)
		}

		// Calculate the number of days until the certificate expires
		daysUntilExpiry := int(time.Until(parsedCert.NotAfter).Hours() / 24)

		log.Info().
			Str("cert", cfg.CertFile).
			Str("key", cfg.KeyFile).
			Int("days_until_expiry", daysUntilExpiry).
			Msg("LoadX509KeyPair")
	} else {
		log.Error().Err(err).Msgf("unable to load keypair for TLS: %s", err)
		return fmt.Errorf("unable to load keypair for TLS: %w", err)
	}
	log.Info().Msgf("success loading keypair for TLS: [public: '%s', private: '%s']", cfg.CertFile, cfg.KeyFile)
	server := http.Server{
		Addr:              cfg.ServerAddress,
		TLSConfig:         tlsConfig, // optional
		ReadHeaderTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, request *http.Request) {
				defer request.Body.Close()

				errorOut := func(message string) {
					log.Error().Err(err).Msg(message)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				if body, err := io.ReadAll(request.Body); err != nil {
					errorOut("error reading v1.AdmissionReview request")
				} else {
					review := new(v1.AdmissionReview)
					start := time.Now()

					if err := json.Unmarshal(body, review); err != nil {
						errorOut("unable to unmarshal v1.AdmissionReview")
					} else {
						fail := func(action string, reason metav1.StatusReason, err error) {
							message := fmt.Sprintf("%s: %s", action, err)
							review.Response = &v1.AdmissionResponse{
								UID:     review.Request.UID,
								Allowed: true,
								Result: &metav1.Status{
									Message: message,
									Reason:  reason,
									Status:  metav1.StatusFailure,
								},
							}
							log.Error().
								Err(err).
								Str("action", action).
								Str("reason", string(reason)).
								Msg("failure")
						}

						var secret corev1.Secret

						if err := json.Unmarshal(review.Request.Object.Raw, &secret); err != nil {
							log.Error().Err(err).Msg("unable to unmarshal the Secret from the v1.AdmissionReview")
							fail("unable to unmarshal the Secret from the v1.AdmissionReview", metav1.StatusReasonBadRequest, err)
						} else if review.Response, err = injector.Inject(secret, review.Request.UID, *credentials, log); err != nil {
							log.Error().Err(err).Msg("calling injector.Inject")
							fail("calling injector.Inject", metav1.StatusReasonInvalid, err)
						}

						if response, err := json.Marshal(review); err != nil {
							errorOut("unable to marshal v1.AdmissionReview response")
						} else {
							w.WriteHeader(http.StatusOK)
							_, err := w.Write(response)
							if err != nil {
								fail("unable to write v1.AdmissionReview response", metav1.StatusReasonInternalError, err)
							}
						}
						log.Info().
							Str("secretname", secret.Name).
							Dur("duration", time.Since(start)).
							Msg("injection complete")
					}
				}
			},
		),
	}

	log.Info().Msg("listening for v1.AdmissionReview requests")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Error().Err(err).Msg("failure with ListenAndServeTLS")
		return fmt.Errorf("failure with ListenAndServeTLS: %w", err)
	}
	return nil
}
