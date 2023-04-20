package main

import (
	"crypto/tls"
	"encoding/json"
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

	env "github.com/caarlos0/env/v6"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// ExitFailure is exit code sent for failed task.
	exitFailure = 1
	// ExitSuccess is exit code sent for running without any error.
	exitSuccess = 0
)

//nolint:gochecknoglobals // ok for providing as version output
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// main is the entry point for the injector; creates an HTTPS listener and listing for v1.AdmissionReview requests
func main() {
	if err := Run(os.Args, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFailure)
	}
}

// InitLogger sets up the logger magic
// By default this is only configured to do pretty console output.
// JSON structured logs are also possible, but not in my default template layout at this time.
func InitLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log.Info().Msg("logger initialized")

}

// Run contains the actual invocation code for the injector and is public to allow running integration tests with it.
func Run(args []string, stdout io.Writer) error {

	type Config struct {
		CertFile            string `env:"DSV_CERT"  envDefault:"${HOME}/tls/cert.pem" envExpand:"true"`                       // Cert is the path to the public certificate file in PEM format.
		KeyFile             string `env:"DSV_KEY" envDefault:"${HOME}/tls/key.pem" envExpand:"true"`                          // Key is the path to the private key file in PEM format.
		CredentialsJsonFile string `env:"DSV_CREDENTIALS_JSON" envDefault:"${HOME}/credentials/config.json" envExpand:"true"` // CredentialsJsonFile is the path to the JSON formatted credentials file that is mounted as a secret.
		ServerAddress       string `env:"DSV_SERVER_ADDRESS" envDefault:":18543"`                                             // ServerAddress is the address to listen on, e.g., 'localhost:8080' or ':8443'
		Debug               bool   `env:"DSV_DEBUG" envDefault:"false"`                                                       // Debug enables debug logging.
	}

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse environment variables")
	}
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("debug logging enabled")
	}

	credentials, err := config.GetCredentials(cfg.CredentialsJsonFile)
	if err != nil {
		log.Error().Err(err).Str("credential-json", cfg.CredentialsJsonFile).Msg("unable to process credentials file")
		return fmt.Errorf("unable to process credentials file '%s': %s", cfg.CredentialsJsonFile, err)
	}
	log.Info().
		Str("credential_names", strings.Join(credentials.Names(), ", ")).
		Str("credential_file", cfg.CredentialsJsonFile).
		Msg("credentials loaded from JSON file")
	var tlsConfig *tls.Config
	if cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile); err == nil {
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		log.Info().Str("cert", cfg.CertFile).Str("key", cfg.KeyFile).Msg("LoadX509KeyPair")
	} else {
		log.Error().Err(err).Msgf("unable to load keypair for TLS: %s", err)
	}
	log.Info().Msgf("success loading keypair for TLS: [public: '%s', private: '%s']", cfg.CertFile, cfg.KeyFile)

	server := http.Server{
		Addr:      cfg.ServerAddress,
		TLSConfig: tlsConfig, // optional
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
							log.Info().Msgf("[ERROR] %s", message)
						}

						var secret corev1.Secret

						if err := json.Unmarshal(review.Request.Object.Raw, &secret); err != nil {
							log.Error().Err(err).Msg("unable to unmarshal the Secret from the v1.AdmissionReview")
							fail("unable to unmarshal the Secret from the v1.AdmissionReview", metav1.StatusReasonBadRequest, err)
						} else if review.Response, err = injector.Inject(secret, review.Request.UID, *credentials); err != nil {
							log.Error().Err(err).Msg("calling injector.Inject")
							fail("calling injector.Inject", metav1.StatusReasonInvalid, err)
						}

						if response, err := json.Marshal(review); err != nil {
							errorOut("unable to marshal v1.AdmissionReview response")
						} else {
							w.WriteHeader(http.StatusOK)
							w.Write(response)
						}
						log.Info().
							Str("secretname", secret.Name).
							Msgf("took %s", time.Since(start))
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
