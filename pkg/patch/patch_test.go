// Patch test file runs integration tests against a secret, based on reaching out to DSV.
//
// - Create a secret with the following command to test against: `dsv secret create --path 'k8s:sync:test' --data '{"password": "","username": ""}'`
// - Hard delete to reset test: `dsv secret delete --path 'k8s:sync:test' --force`
// - Rollback to a prior version: `dsv secret rollback --path 'k8s:sync:test' --version 0`.
package patch

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"

	"github.com/DelineaXPM/dsv-k8s/v2/internal/test"
	"github.com/mattbaird/jsonpatch"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//nolint:gochecknoglobals // Test file, ok to exclude from this condition, as long as Parallel testing is handled.
var (
	secretPath  = test.SecretPath()
	credentials = test.Credentials()
)

const (

	// The credentials annotation used for authentication.
	credentialsAnnotationValue = "default"
)

// Ensure log output doesn't pollute tests.
func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
}

// getCredentialAnnotationValue is a test helper function to get the credentialsAnnotationValue from either the environment variable or default to the constant value if not set.
func getCredentialAnnotationValue(t *testing.T) string {
	t.Helper()

	if os.Getenv("DSV_CREDENTIALS_ANNOTATION_VALUE") == "default" {
		t.Fatal("Skipping test as credentialsAnnotationValue is 'default', should have 'app1' to run full tests")
	}

	if os.Getenv("DSV_CREDENTIALS_ANNOTATION_VALUE") != "" {
		return os.Getenv("DSV_CREDENTIALS_ANNOTATION_VALUE")
	}
	return credentialsAnnotationValue
}

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

// prettyPrintJSON is a test helper function to pretty print a JSON string.
func prettyPrintJSON(t *testing.T, in string) string {
	t.Helper()
	var buf bytes.Buffer

	if err := json.Indent(&buf, []byte(in), "", "  "); err != nil {
		t.Errorf("non-critical helper for pretty printing json: %v", err)
	}
	return buf.String()
}

// tp tests the patching of a secret
func tp(t *testing.T, secret v1.Secret, count int, op, path, value string) []jsonpatch.JsonPatchOperation {
	t.Helper()
	var ops []jsonpatch.JsonPatchOperation
	ops, err := makePatchOperations(secret, credentials)
	// If the patch operations errors out entirely, then return this as a failure.
	if err != nil {
		t.Logf("makePatchOperations: %+v", ops)
		t.Error(err)
		return nil
	}

	// Output the operations json for review, followed by evaluating operation and value to see if they match.
	var opsDebugOutput string
	for _, item := range ops {
		opsDebugOutput += prettyPrintJSON(t, item.Json())
	}
	t.Logf(opsDebugOutput)

	// When the count of operations expected isn't matching.
	if count != len(ops) {
		t.Errorf("expected %d operations, got %d", count, len(ops))
		return ops
	}

	// The count of operations matches, but the operation didn't seem to proceed as expected.
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

	return ops
}

// This requires a secret created in DSV to be existing as defined in the README.md
// Look at the Test setup requirements section.
func TestCredentialsLogic(t *testing.T) {
	annotations := map[string]string{setAnnotation: secretPath}
	secret := makeSecret(annotations, map[string][]byte{})

	// TestCredentialsLogic should return zero operations when providing an nonexistent credential mapping.
	tp(t, secret, 0, "", "", "")

	secretData := map[string][]byte{
		"username": []byte("root"),
		"password": []byte("password"),
	}

	secret.Data = secretData
	tp(t, secret, 0, "", "", "")
	secret.ObjectMeta.Annotations[credentialsAnnotation] = getCredentialAnnotationValue(t)
	tp(t, secret, 4, "replace", "/data/password", "admin")
}

func TestVersionLogic(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: getCredentialAnnotationValue(t),
		setAnnotation:         secretPath,
		versionAnnotation:     "0",
	}
	secret := makeSecret(annotations, map[string][]byte{})

	tp(t, secret, 0, "", "", "")
	delete(secret.ObjectMeta.Annotations, versionAnnotation)
	tp(t, secret, 4, "add", "/data/password", "admin")
}

func TestAddOperationSetAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: getCredentialAnnotationValue(t),
		setAnnotation:         secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{})

	tp(t, secret, 4, "add", "/data/password", "admin")
}

func TestReplaceOperationSetAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: getCredentialAnnotationValue(t),
		setAnnotation:         secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{
		"username": []byte("root"),
		"password": []byte("password"),
	})

	tp(t, secret, 4, "replace", "/data/password", "admin")
}

func TestRemoveOperationSetAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: getCredentialAnnotationValue(t),
		setAnnotation:         secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{
		"username": []byte("root"),
		"password": []byte("password"),
		"domain":   []byte("anything"),
	})

	tp(t, secret, 5, "remove", "/data/domain", "")
}

func TestRemoveOperationUpdateAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: getCredentialAnnotationValue(t),
		updateAnnotation:      secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{
		"username": []byte("root"),
		"password": []byte("password"),
		"domain":   []byte("anything"),
	})

	ops := tp(t, secret, 4, "replace", "/data/password", "admin")
	for _, op := range ops {
		if op.Operation == "remove" {
			t.Error("expected no remove operations")
		}
	}
}

func TestAddOperationAddAnnotation(t *testing.T) {
	annotations := map[string]string{
		credentialsAnnotation: getCredentialAnnotationValue(t),
		addAnnotation:         secretPath,
	}
	secret := makeSecret(annotations, map[string][]byte{"username": []byte("root")})

	ops := tp(t, secret, 3, "add", "/data/password", "admin")
	for _, op := range ops {
		if op.Path == "/metadata/annotations" {
			continue
		}
		if op.Operation == "replace" || op.Operation == "remove" {
			t.Error("expected no remove or replace operations")
		}
	}
}
