package patch

import (
	"encoding/base64"
	"testing"

	"github.com/DelineaXPM/dsv-k8s/internal/test"
	"github.com/mattbaird/jsonpatch"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var secretPath = test.SecretPath()
var credentials = test.Credentials()

func decodeBase64(v interface{}) string {
	r, _ := base64.StdEncoding.DecodeString(v.(string))
	return string(r)
}

func makeSecret(annotations map[string]string, data map[string][]byte) v1.Secret {
	return v1.Secret{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Annotations: annotations},
		Data:       data,
		StringData: map[string]string{},
		Type:       "",
	}
}

// tp tests the patching of a secret
func tp(t *testing.T, secret v1.Secret, count int, op, path, value string) []jsonpatch.JsonPatchOperation {
	if ops, err := makePatchOperations(secret, credentials); err != nil {
		t.Error(err)
		return nil
	} else {
		if count != len(ops) {
			t.Errorf("expected %d operations, got %d", count, len(ops))
		} else {
			for _, o := range ops {
				if path == o.Path {
					if op != o.Operation {
						t.Errorf("expected %s operation, got %s", op, o.Operation)
					}
					if o.Value != nil && value != decodeBase64(o.Value) {
						t.Errorf("expected %s, got %s", value, decodeBase64(o.Value))
					}
				}
			}
		}
		return ops
	}
}

/*
	These below assume a secret with path dsvSecretPath containing data and version:
		{
			"data": {
				"password": "admin",
				"username": "admin"
			},
			"version": "0"
		}
*/
func TestCredentialsLogic(t *testing.T) {
	annotations := map[string]string{setAnnotation: secretPath}
	secret := makeSecret(annotations, map[string][]byte{})

	tp(t, secret, 0, "", "", "")

	secretData := map[string][]byte{
		"username": []byte("root"),
		"password": []byte("password"),
	}

	secret.Data = secretData

	tp(t, secret, 0, "", "", "")

	secret.ObjectMeta.Annotations[credentialsAnnotation] = "app1"

	tp(t, secret, 3, "replace", "/data/password", "admin")
}

func TestVersionLogic(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: "app1",
		setAnnotation:         secretPath,
		versionAnnotation:     "0",
	}
	secret := makeSecret(annotations, map[string][]byte{})

	tp(t, secret, 0, "", "", "")
	delete(secret.ObjectMeta.Annotations, versionAnnotation)
	tp(t, secret, 3, "add", "/data/password", "admin")
}

func TestAddOperationSetAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: "app1",
		setAnnotation:         secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{})

	tp(t, secret, 3, "add", "/data/password", "admin")
}

func TestReplaceOperationSetAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: "app1",
		setAnnotation:         secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{
		"username": []byte("root"),
		"password": []byte("password"),
	})

	tp(t, secret, 3, "replace", "/data/password", "admin")
}

func TestRemoveOperationSetAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: "app1",
		setAnnotation:         secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{
		"username": []byte("root"),
		"password": []byte("password"),
		"domain":   []byte("anything"),
	})

	tp(t, secret, 4, "remove", "/data/domain", "")
}

func TestRemoveOperationUpdateAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: "app1",
		updateAnnotation:      secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{
		"username": []byte("root"),
		"password": []byte("password"),
		"domain":   []byte("anything"),
	})

	ops := tp(t, secret, 3, "replace", "/data/password", "admin")
	for _, op := range ops {
		if op.Operation == "remove" {
			t.Error("expected no remove operations")
		}
	}
}

func TestAddOperationAddAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: "app1",
		addAnnotation:         secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{"username": []byte("root")})

	ops := tp(t, secret, 2, "add", "/data/password", "admin")
	for _, op := range ops {
		if op.Path == "/metadata/annotations" {
			continue
		}
		if op.Operation == "replace" || op.Operation == "remove" {
			t.Error("expected no remove or replace operations")
		}
	}
}
