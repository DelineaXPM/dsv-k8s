# Using

Once installed, any correctly annotated Kubernetes Secrets are modified on create and update.

The four annotations that affect the behavior of the webhook are:

```golang
const(
    credentialsAnnotation = "dsv.delinea.com/credentials"
    setAnnotation         = "dsv.delinea.com/set-secret"
    addAnnotation         = "dsv.delinea.com/add-to-secret"
    updateAnnotation      = "dsv.delinea.com/update-secret"
)
```

`credentialsAnnotation` selects the credentials that the injector uses to retrieve the DSV Secret.
If the credentials are present, it must map to Client Credential and Tenant mapping.
The injector will use the _default_ Credential and Tenant mapping unless the `credentialsAnnotation` is declared.
The `setAnnotation`, `addAnnotation` and `updateAnnotation`, must contain the path to the DSV Secret that the injector will use to mutate the Kubernetes Secret.

- `addAnnotation` adds missing fields without overwriting or removing existing fields.
- `updateAnnotation` adds and overwrites existing fields but does not remove fields.
- `setAnnotation` overwrites fields and removes fields that do not exist in the DSV Secret.
  NOTE: A Kubernetes Secret should specify only one of the "add," "update,"
  or "set" annotations.
  The order of precedence is `setAnnotation`,
  then `addAnnotation`, then `updateAnnotation` when multiple are present.

## Examples

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: example-secret
  annotations:
    dsv.delinea.com/credentials: app1
    dsv.delinea.com/set-secret: test:secret
type: Opaque
data:
  username:
  domain:
  password:
```

The above example specifies credentials, so a mapping for those credentials must exist in the current webhook configuration.
It uses the `setAnnotation`, so the data in the injector will overwrite the existing contents of the Kubernetes Secret.

If `/test/secret` contains a `username` and `password` but no `domain`, then the k8s secret would get the `username` and `password` from dsv secret data and the injector will remove the `domain` field.

There are more examples in the `examples` directory.
They show how the different annotations work.
