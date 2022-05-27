package injector

import (
	"fmt"
	"log"

	"github.com/DelineaXPM/dsv-k8s/v2/pkg/config"
	patch "github.com/DelineaXPM/dsv-k8s/v2/pkg/patch"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Inject adds to, updates or replaces the k8s Secret.Data with dsv Secret.Data (see above)
func Inject(secret corev1.Secret, UID types.UID, credentials config.Credentials) (*v1.AdmissionResponse, error) {
	if jsonPatch, err := patch.GenerateJsonPatch(secret, credentials); err != nil {
		return nil, fmt.Errorf("unable to generate JSON patch for Secret '%s': %s", secret.Name, err)
	} else if jsonPatch != nil {
		patchType := v1.PatchTypeJSONPatch

		log.Printf("[INFO] patching k8s Secret '%s'", secret.Name)
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
		log.Printf("[DEBUG] k8s Secret '%s' did not require patching", secret.Name)
		return &v1.AdmissionResponse{
			Allowed: true,
			Result: &metav1.Status{
				Status: metav1.StatusSuccess,
			},
			UID: UID,
		}, nil
	}
}
