# Thycotic DevOps Secrets Vault Kubernetes Secret Injector

A [Kubernetes](https://kubernetes.io/)
[Mutating Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#admission-webhooks)
that injects Secret data from Thycotic DevOps Secrets Vault (DSV) into Kubernetes (k8s) cluster Secrets. The webhook is made available to the Kubernetes cluster as the `injector-svc` which can be hosted in a k8s POD or as a stand-alone service.

The webhook intercepts `CREATE` and `UPDATE` Secret admissions and supplements
or overwrites the Secret data with Secret data from DSV. The webhook
configuration is a set of DSV Role to Client Credential and Tenant mappings.
The webhook updates k8s Secrets based on a set of annotations (see below).

The webhook uses the [DSV Go SDK](https://github.com/thycotic/dsv-sdk-go) to
communicate with DSV.

It was built and tested with [Minikube](https://minikube.sigs.k8s.io/).

## Configure

The webhook requires a JSON formatted list of DSV Role to Client Credential and
Tenant mappings, stored in `configs/roles.json`:

```json
{
    "my-role": {
        "credentials": {
            "clientId": "93d866d4-635f-4d4e-9ce3-0ef7f879f319",
            "clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxxxx"
        },
        "tenant": "mytenant"
    },
    "default": {
        "credentials": {
            "clientId": "64241412-3934-4aed-af26-95b1eaba0e6a",
            "clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxxxx"
        },
        "tenant": "mytenant"
    }
}
```

The _default_ Role is used when the k8s secret being modified does not
specify a particular Role.

## Build

The build requires a functional Minikube installation. It uses a the `Makefile`.

To build the [Docker](https://www.docker.com/) image run:

```sh
make image
```

To configure the Kubernetes cluster to to call the webhook as stand-alone
service running on the host at `$(HOST_IP)`:

```sh
make deploy_host
```

Note that `$(HOST_IP)` defaults to the IP address of the host executing the
build. Also note that the `Makefile` does not include`docker login` or
`docker push` as neither are required to deploy the webhook to Minikube.

The `Makefile` requires access to the Kubernetes CA certificate. It assumes that
it is available in `${HOME}/.minikube/ca.crt` but that can be overridden by
setting `$(CA_CRT)`.

## Use

Once the `injector-svc` is made available to the Kubernetes cluster and the
webhook is configured to call it, any appropriately annotated k8s secrets will
be modified by it whenever they are created or updated.

The four annotations that affect the behavior of the webhook are:

```golang
const(
    roleAnnotation   = "dsv.thycotic.com/role"
    setAnnotation    = "dsv.thycotic.com/set-secret"
    addNotation      = "dsv.thycotic.com/add-to-secret"
    updateAnnotation = "dsv.thycotic.com/update-secret"
)
```

`roleAnnotation` sets the DSV Role that k8s should use to access the DSV Secret
that will be used to modify the k8s Secret. If it is present then the Role
must exist in the above mentioned Role mappings that the webhook is configured
to use. If it is absent then the _default_ mapping is used.

The `setAnnotation`, `addAnnoation` and `updateAnnotation` contain the path to
the DSV Secret that will be used to modified the k8s Secret being admitted.

* `addAnnotation` adds missing fields without overrwriting or removing existing fields.
* `updateAnnotation` adds and overwrites existing fields but does not remove fields.
* `setAnnotation` overwrites fields and removes fields that do not exist in the DSV Secret.

Only one of these should be specified on any given k8s Secret, however, if more
than one are defined then the order of precendence is `setAnnotation` then
`addAnnotation` then `updateAnnotation`.

### Examples

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: example-secret
  annotations:
    dsv.thycotic.com/role: my-role
    dsv.thycotic.com/set-secret: /test/secret
type: Opaque
data:
  username: dW5tb2RpZmllZC11c2VybmFtZQ==
  domain: dW5tb2RpZmllZC1kb21haW4=
  password: dW5tb2RpZmllZC1wYXNzd29yZA==
```

The above example specifies a Role so a mapping for that role must exist in the
current webhook configuration. It uses the `setAnnotation` so the data in the
secret will be overwritten; if `/test/secret` contains a `username` and
`password` but no `domain` then the secret would contain the `username` and
`password` from the DSV Secret Data and the `domain` field will be removed.

There are more examples in the `examples` directory.
