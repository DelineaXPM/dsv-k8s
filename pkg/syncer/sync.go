// package syncer handles the syncing actions for secret reading, patching and injecting into kubernetes secrets.
package syncer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/DelineaXPM/dsv-k8s/v2/internal/k8s"
	"github.com/DelineaXPM/dsv-k8s/v2/pkg/config"
	"github.com/DelineaXPM/dsv-k8s/v2/pkg/patch"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Sync does the same thing as Inject, but by iterating over the existing k8s Secrets
func Sync(config k8s.Config, namespace string, credentials config.Credentials, log zerolog.Logger) error {
	secretsClient, err := k8s.GetSecretsClient(config, namespace)
	if err != nil {
		return fmt.Errorf("[ERROR] error getting a Kubernetes Client API Secrets Client: %w", err)
	}
	log.Debug().Msgf("getting a list of Secrets in namespace %q", namespace)

	if secrets, err := secretsClient.List(context.TODO(), metav1.ListOptions{}); err != nil {
		return fmt.Errorf("[ERROR] unable to get a list of secrets in namespace %q: %w", namespace, err)
	} else {
		wg := sync.WaitGroup{}

		log.Info().Msgf("processing %d Secrets", len(secrets.Items))
		for _, secret := range secrets.Items {
			log.Debug().Msgf("processing k8s Secret %q", secret.Name)
			wg.Add(1)
			go pp(secret, credentials, config, &wg, log)
		} // TODO: put an upper limit on the number of goroutines to spawn in one go
		wg.Wait()
		if secrets.RemainingItemCount != nil && *secrets.RemainingItemCount > 0 {
			log.Warn().Msgf("this server pages; %d Secrets were not processed", secrets.RemainingItemCount)
		}
	}
	return nil
}

// pp possibly patches the Kubernetes Secret
func pp(secret corev1.Secret, credentials config.Credentials, config k8s.Config, wg *sync.WaitGroup, log zerolog.Logger) {
	defer wg.Done()
	start := time.Now()
	defer func() {
		log.Debug().
			Dur("duration", time.Since(start)).Str("secret_name", secret.Name).
			Msg("possible patch complete")
	}()

	if jsonPatch, err := patch.GenerateJsonPatch(secret, credentials); err != nil {
		log.Error().
			Err(err).
			Str("secret_name", secret.Name).
			Str("secret_namespace", secret.Namespace).
			Msg("patch.GenerateJsonPatch")
	} else if jsonPatch == nil {
		log.Debug().Msgf("k8s Secret %q' did not require patching", secret.Name)
	} else {
		if secretsClient, err := k8s.GetSecretsClient(config, secret.Namespace); err != nil {
			log.Error().
				Err(err).
				Str("secret_name", secret.Name).
				Str("secret_namespace", secret.Namespace).
				Msg("k8s.GetSecretsClient")
		} else {
			if result, err := secretsClient.Patch(
				context.TODO(), secret.Name, types.JSONPatchType, jsonPatch, metav1.PatchOptions{},
			); err != nil {
				log.Error().
					Err(err).
					Str("secret_name", secret.Name).
					Str("secret_namespace", secret.Namespace).
					Msg("secretsClient.Patch")
			} else {
				log.Debug().Msgf("patched k8s Secret %q", result.Name)
			}
		}
	}
}
