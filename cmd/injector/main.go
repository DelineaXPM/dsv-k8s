package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/DelineaXPM/dsv-k8s/v2/pkg/config"
	"github.com/DelineaXPM/dsv-k8s/v2/pkg/injector"

	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// main is the entry point for the injector; creates an HTTPS listener and listing for v1.AdmissionReview requests
func main() {
	var certFile, keyFile, credentialsFile string

	flag.StringVar(&certFile, "cert", "tls/cert.pem", "the path of the public certificate file in PEM format")
	flag.StringVar(&keyFile, "key", "tls/key.pem", "the path of the private key file in PEM format")
	flag.StringVar(
		&credentialsFile,
		"credentials",
		"credentials/config.json",
		"the path of JSON formatted credentials file",
	)

	server := new(http.Server)

	flag.StringVar(&server.Addr, "address", ":18543", "the address to listen on, e.g., 'localhost:8080' or ':8443'")
	flag.Parse()

	credentials, err := config.GetCredentials(credentialsFile)
	if err != nil {
		log.Fatalf("unable to process credentials file '%s': %s", credentialsFile, err)
	}
	log.Printf("[INFO] success loading %d credential sets: [%s] from '%s'",
		len(*credentials),
		strings.Join(credentials.Names(), ", "),
		credentialsFile,
	)

	if cert, err := tls.LoadX509KeyPair(certFile, keyFile); err == nil {
		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	} else {
		log.Fatalf("unable to load keypair for TLS: %s", err)
	}
	log.Printf("[INFO] success loading keypair for TLS: [public: '%s', private: '%s']", certFile, keyFile)

	server.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, request *http.Request) {
			defer request.Body.Close()

			errorOut := func(message string) {
				log.Printf("[ERROR] %s: %s", message, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if body, err := ioutil.ReadAll(request.Body); err != nil {
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
						log.Printf("[ERROR] %s", message)
					}

					var secret corev1.Secret

					if err := json.Unmarshal(review.Request.Object.Raw, &secret); err != nil {
						fail("unable to unmarshal the Secret from the v1.AdmissionReview", metav1.StatusReasonBadRequest, err)
					} else if review.Response, err = injector.Inject(secret, review.Request.UID, *credentials); err != nil {
						fail("calling injector.Inject", metav1.StatusReasonInvalid, err)
					}

					if response, err := json.Marshal(review); err != nil {
						errorOut("unable to marshal v1.AdmissionReview response")
					} else {
						w.WriteHeader(http.StatusOK)
						w.Write(response)
					}
					log.Printf("[INFO] processing Secret '%s' took %s", secret.Name, time.Since(start))
				}
			}
		},
	)

	log.Printf("[INFO] listening for v1.AdmissionReview requests on '%s'", server.Addr)
	log.Fatalf("[FATAL] %s", server.ListenAndServeTLS("", ""))
}
