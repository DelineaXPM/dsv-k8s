package syncer

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/DelineaXPM/dsv-k8s/v2/internal/k8s"
	"github.com/DelineaXPM/dsv-k8s/v2/pkg/config"
	"github.com/DelineaXPM/dsv-k8s/v2/pkg/patch"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// pp possibly patches the Kubernetes Secret
func pp(secret corev1.Secret, credentials config.Credentials, config k8s.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	logFailure := func(secret corev1.Secret, err error) {
		log.Printf("[ERROR] unable to patch Secret '%s' in namespace '%s': %s", secret.Name, secret.Namespace, err)
	}

	if jsonPatch, err := patch.GenerateJsonPatch(secret, credentials); err != nil {
		logFailure(secret, err)
	} else if jsonPatch == nil {
		log.Printf("[DEBUG] k8s Secret '%s' did not require patching", secret.Name)
	} else {
		if secretsClient, err := k8s.GetSecretsClient(config, secret.Namespace); err != nil {
			logFailure(secret, err)
		} else {
			if result, err := secretsClient.Patch(
				context.TODO(), secret.Name, types.JSONPatchType, jsonPatch, metav1.PatchOptions{},
			); err != nil {
				logFailure(secret, err)
			} else {
				log.Printf("[DEBUG] patched k8s Secret '%s'", result.Name)
			}
		}
	}
}

// Sync does the same thing as Inject, but by iterating over the existing k8s Secrets
func Sync(config k8s.Config, namespace string, credentials config.Credentials) error {
	secretsClient, err := k8s.GetSecretsClient(config, namespace)

	if err != nil {
		return fmt.Errorf("[ERROR] error getting a Kubernetes Client API Secrets Client: %s", err)
	}
	log.Printf("[DEBUG] getting a list of Secrets in namespace '%s'", namespace)

	if secrets, err := secretsClient.List(context.TODO(), metav1.ListOptions{}); err != nil {
		return fmt.Errorf("[ERROR] unable to get a list of secrets in namespace '%s': %s", namespace, err)
	} else {
		wg := sync.WaitGroup{}

		log.Printf("[INFO] processing %d Secrets", len(secrets.Items))
		for _, secret := range secrets.Items {
			log.Printf("[DEBUG] processing k8s Secret '%s'", secret.Name)
			wg.Add(1)
			go pp(secret, credentials, config, &wg)
		} // TODO: put an upper limit on the number of goroutines to spawn in one go
		wg.Wait()
		if secrets.RemainingItemCount != nil && *secrets.RemainingItemCount > 0 {
			log.Printf("[WARN] this server pages; %d Secrets were not processed", secrets.RemainingItemCount)
		}
	}
	return nil
}
