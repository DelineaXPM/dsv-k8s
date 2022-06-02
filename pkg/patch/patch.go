package patch

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DelineaXPM/dsv-k8s/v2/pkg/config"
	"github.com/DelineaXPM/dsv-sdk-go/v2/vault"
	"github.com/mattbaird/jsonpatch"
	v1 "k8s.io/api/core/v1"
)

/*
credentialsAnnotation contains a name that maps to a set of credentials in Credentials (see below)
setAnnotation, addAnnoation and updateAnnotation contain the path to the DSV Secret
that will be used to modified this Secret.
addAnnotation adds missing fields without overwriting or removing existing fields
updateAnnotation adds and overwrites existing fields but does not remove fields
setAnnotation overwrites fields and removes fields that do not exist in the DSV Secret
*/
const (
	credentialsAnnotation = "dsv.delinea.com/credentials"
	setAnnotation         = "dsv.delinea.com/set-secret"
	addAnnotation         = "dsv.delinea.com/add-to-secret"
	updateAnnotation      = "dsv.delinea.com/update-secret"
	tsAnnotation          = "dsv.delinea.com/modified"
	versionAnnotation     = "dsv.delinea.com/version"
)

/*
noPatch causes Inject to approve without patching
patchAdd causes Inject to add entries without overwriting or removing existing ones
patchUpdate causes Inject to add and update entries without removing existing ones
patchOverwrite causes Inject to completely overwrite all entries including removal
of existing entries that are not present in the DSV Secret
*/
const (
	noPatch = iota
	patchAdd
	patchUpdate
	patchOverwrite
)

// GeneratePatch generates a JSON Patch that applies changes to the k8s Secret based on the DSV Secret that it refers to
func GenerateJsonPatch(secret v1.Secret, credentials config.Credentials) ([]byte, error) {
	if patchOperations, err := makePatchOperations(secret, credentials); err != nil {
		return nil, err
	} else if patchOperations != nil {
		if jsonPatch, err := json.Marshal(patchOperations); err != nil {
			return nil, fmt.Errorf("unable to marshal JsonPatch: %s", err)
		} else {
			return jsonPatch, nil
		}
	} else {
		return nil, nil
	}
}

func makePatchOperations(secret v1.Secret, credentials config.Credentials) ([]jsonpatch.JsonPatchOperation, error) {
	annotations := secret.ObjectMeta.Annotations
	patchMode := noPatch
	var secretPath string
	var ok bool

	if secretPath, ok = annotations[setAnnotation]; ok {
		patchMode = patchOverwrite
	} else if secretPath, ok = annotations[addAnnotation]; ok {
		patchMode = patchAdd
	} else if secretPath, ok = annotations[updateAnnotation]; ok {
		patchMode = patchUpdate
	}

	if patchMode == noPatch {
		return nil, nil
	}

	var config vault.Configuration
	/*
		If there is a credentials annotation, use the credentials by that name
		and return an error if there are no credentials for that name.
		Otherwise, use the default credentials or, finally,
		do nothing if there aren't any.
	*/
	if name, ok := annotations[credentialsAnnotation]; ok {
		if credentials, ok := credentials[name]; ok {
			config = credentials.Configuration
		} else {
			return nil, fmt.Errorf("no credentials for: %s", name)
		}
	} else if credentials, ok := credentials["default"]; ok {
		config = credentials.Configuration
	} else {
		return nil, nil
	}

	var vaultSecret *vault.Secret

	/*
		If there's a patch annotation, and a credentials annotation that matches
		a set of credentials, use them to get a Vault Secret.
	*/
	if vault, err := vault.New(config); err != nil {
		return nil, fmt.Errorf("configuration error: %s", err)
	} else if vaultSecret, err = vault.Secret(secretPath); err != nil {
		return nil, fmt.Errorf("unable to get the secret: %s", err)
	}
	/*
		Use the version annotation to determine if the secret has been modified
	*/
	if vaultSecret.Version == annotations[versionAnnotation] {
		return nil, nil
	}
	/*
		CreatePatch returns the difference between the JSON representation of
		the DSV Secret Data and the k8s Secret Data, as an RFC 6902 JSON Patch.
	*/
	vsdj, _ := json.Marshal(vaultSecret.Data)
	sdj, _ := json.Marshal(secret.Data)
	diff, _ := jsonpatch.CreatePatch(sdj, vsdj)
	ops := []jsonpatch.JsonPatchOperation{}
	/*
		Each patch operation has to be updated so that k8s can apply it to
		the entire Secret:
		1) the Path must be relative to /Secret rather than /Secret/Data
		2) the Values must be Base64 encoded
		3) the Operations that conflict with patchMode must be removed
	*/
	for _, op := range diff {
		op.Path = "/data" + op.Path

		switch v := op.Value.(type) {
		case []byte:
			op.Value = base64.StdEncoding.EncodeToString(v)
		case string:
			op.Value = base64.StdEncoding.EncodeToString([]byte(v))
		case map[string]interface{}:
			if json, err := json.Marshal(v); err != nil {
				return nil, fmt.Errorf("unable to marshal value for %s operation on %s: %s",
					op.Operation, op.Path, err)
			} else {
				op.Value = base64.StdEncoding.EncodeToString([]byte(json))
			}
		}
		// remove operations that conflict with the patchMode
		switch op.Operation {
		case "replace":
			if patchMode == patchAdd {
				continue
			}
		case "remove":
			if patchMode == patchAdd || patchMode == patchUpdate {
				continue
			}
		}
		ops = append(ops, op)
	}
	/*
		If there is at least one patch operation add an operation to replace update the annotations
	*/
	if len(ops) > 0 {
		ops = append(ops, jsonpatch.JsonPatchOperation{
			Operation: "replace",
			Path:      "/metadata/annotations",
			Value: map[string]string{
				tsAnnotation:      time.Now().UTC().Format(time.UnixDate),
				versionAnnotation: vaultSecret.Version,
			},
		})
		return ops, nil
	}
	return nil, nil
}
