package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/thycotic/dsv-k8s/pkg/config"
	"github.com/thycotic/dsv-k8s/pkg/syncer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// main is the entrypoint for the syncer; parses the credentials file and calls syncer.Sync
func main() {
	kubeConfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	var namespace, credentialsFile string

	flag.StringVar(&kubeConfig, "kubeConfig", kubeConfig,
		"the Kubernetes Client API configuration file; ignored when running in-cluster",
	)
	flag.StringVar(&namespace, "namespace", "",
		"the Kubernetes namespace containing the Secrets to sync; \"\" (the default) for all namespaces",
	)
	flag.StringVar(&credentialsFile, "credentials", "credentials/config.json",
		"the path of JSON formatted credentials file",
	)
	flag.Parse()

	if credentials, err := config.GetCredentials(credentialsFile); err != nil {
		log.Fatalf("[ERROR] unable to process configuration file '%s': %s", credentialsFile, err)
	} else {
		log.Printf("[INFO] success loading %d credential sets: [%s]", len(*credentials), strings.Join(credentials.Names(), ", "))

		start := time.Now()

		if err := syncer.Sync(
			func(namespace string) (*rest.Config, error) {
				// First assume we are running inside a cluster
				if config, err := rest.InClusterConfig(); err != nil {
					// Failing that, try loading the kubeConfig file
					if config, err := clientcmd.BuildConfigFromFlags("", kubeConfig); err != nil {
						return nil, fmt.Errorf("[ERROR] error getting Kubernetes Client rest.Config: %s", err)
					} else {
						return config, nil
					}
				} else {
					return config, nil
				}
			},
			namespace,
			*credentials,
		); err != nil {
			log.Fatalf("[ERROR] unable to sync Secrets: %s", err)
		}
		log.Printf("[INFO] processing took %s", time.Since(start))
	}
}
