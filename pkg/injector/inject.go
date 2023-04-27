package injector

import (
	"fmt"

	"github.com/DelineaXPM/dsv-k8s/v2/pkg/config"
	patch "github.com/DelineaXPM/dsv-k8s/v2/pkg/patch"
	"github.com/rs/zerolog"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Inject adds to, updates or replaces the k8s Secret.Data with dsv Secret.Data (see above)
func Inject(secret corev1.Secret, UID types.UID, credentials config.Credentials, log zerolog.Logger) (*v1.AdmissionResponse, error) {
	if jsonPatch, err := patch.GenerateJsonPatch(secret, credentials); err != nil {
		log.Error().Err(err).
			Str("secret_name", secret.Name).
			Str("namespace", secret.Namespace).
			Msg("unable to patch secret")
		return nil, fmt.Errorf("unable to generate JSON patch for Secret '%s': %s", secret.Name, err)
	} else if jsonPatch != nil {
		patchType := v1.PatchTypeJSONPatch

		log.Debug().Str("secret_name", secret.Name).Msg("patching secret")
		return &v1.AdmissionResponse{
			Allowed: true,
			Result: &metav1.Status{
				Status: metav1.StatusSuccess,
			},
			UID:       UID,
			PatchType: &patchType,
			Patch:     jsonPatch,
		}, nil
	} else {
		log.Debug().Str("secret_name", secret.Name).Msg("no patching required")
		return &v1.AdmissionResponse{
			Allowed: true,
			Result: &metav1.Status{
				Status: metav1.StatusSuccess,
			},
			UID: UID,
		}, nil
	}
}
