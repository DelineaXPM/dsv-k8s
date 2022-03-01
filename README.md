# Delinea DevOps Secrets Vault Kubernetes Secret Injector

![Docker](https://github.com/thycotic/dsv-k8s/workflows/Docker/badge.svg)
![GitHub Package Registry](https://github.com/thycotic/dsv-k8s/workflows/GitHub%20Package%20Registry/badge.svg)
![Red Hat Quay](https://github.com/thycotic/dsv-k8s/workflows/Red%20Hat%20Quay/badge.svg)

A [Kubernetes](https://kubernetes.io/)
[Mutating Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#admission-webhooks)
that injects Secret data from Delinea DevOps Secrets Vault (DSV) into Kubernetes
Secrets. The webhook can be hosted as a pod or as a stand-alone service.

The webhook works by intercepting `CREATE` and `UPDATE` Secret admissions and
mutating the Secret with data from DSV. The webhook configuration consists of
one or more _role_ to Client Credential Tenant mappings. It updates Kubernetes
Secrets based on annotations on the Secret itself.

The webhook uses the [Golang SDK](https://github.com/thycotic/dsv-sdk-go) to
communicate with the DSV API.

It was tested with [Minikube](https://minikube.sigs.k8s.io/) and
[Minishift](https://docs.okd.io/3.11/minishift/index.html).

## Configure

The webhook requires a JSON formatted list of _role_ to Client Credential and
Tenant mappings. The _role_ is a simple name that does not relate to Kubernetes
Roles. It simply selects which credentials to use to get the Secret from DSV.

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

> NOTE: the injector uses the _default_ role when it mutates a Kubernetes Secret
> that does not have a _roleAnnotation_. [See below](#use)

## Run

The injector is a Golang executable that runs a built-in HTTPS server hosting
the Kubernetes Mutating Webhook Webservice.

```bash
$ /usr/bin/dsv-injector-svc -?
flag provided but not defined: -?
Usage of ./dsv-injector-svc:
  -cert string
        the path of the certificate file in PEM format (default "injector.pem")
  -hostport string
        the host:port e.g. localhost:8080 (default ":18543")
  -key string
        the path of the certificate key file in PEM format (default "injector.key")
  -roles string
        the path of JSON formatted roles file (default "roles.json")
```

Thus the injector can run "anywhere," but, typically, the injector runs as a POD
in the Kubernetes cluster that uses it.

## Build

> NOTE: Building the `dsv-injector` image is not required to install it as it is
> available on multiple public registries.

Building the injector requires [Docker](https://www.docker.com/) or
[Podman](https://podman.io/). To build it, run:

```sh
make image
```

### Minikube and Minishift

Remember to run `eval $(minikube docker-env)` in the shell to push the image to
Minikube's Docker daemon.💡 Likewise for Minishift except its
`eval $(minishift docker-env)`.

### Install

Installation requires [Helm](https://helm.sh).

The `Makefile` demonstrates a typical installation via the
[Helm](https://helm.sh/) chart. It imports `roles.json` as a file that it
templates as a Kubernetes Secret for the injector.

The Helm `values.yaml` file `image.repository` is `thycotic/dsv-injector`:

```yaml
image:
  repository: thycotic/dsv-injector
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: ""
```

That means, by default, `make install` will pull from Docker, GitHub, or Quay.

```sh
make install
```

However, the `Makefile` contains an `install-image` target that configures Helm
to use the image built with `make image`:

```sh
make install-image
```

`make uninstall` uninstalls the Helm Chart.

`make clean` removes the Docker image.

## Use

Once the injector is available in the Kubernetes cluster, and the
[MutatingAdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook)
is in place, any appropriately annotated Kubernetes Secrets are modified on
create and update.

The four annotations that affect the behavior of the webhook are:

```golang
const(
    roleAnnotation   = "dsv.thycotic.com/role"
    setAnnotation    = "dsv.thycotic.com/set-secret"
    addNotation      = "dsv.thycotic.com/add-to-secret"
    updateAnnotation = "dsv.thycotic.com/update-secret"
)
```

`roleAnnotation` selects the credentials that the injector uses to retrieve the
DSV Secret. If the role is present, it must map to Client Credential and Tenant
mapping. If the role is absent, the injector will use the _default_ Credential
and Tenant a mapping.

The `setAnnotation`, `addAnnotation` and `updateAnnotation` contain the path to
the DSV Secret that the injector will use to mutate the Kubernetes Secret.

- `addAnnotation` adds missing fields without overwriting or removing existing
  fields.
- `updateAnnotation` adds and overwrites existing fields but does not remove
  fields.
- `setAnnotation` overwrites fields and removes fields that do not exist in the
  DSV Secret.

NOTE: A Kubernetes Secret should specify only one of the "add," "update," or
"set" annotations. The order of precedence is `setAnnotation`, then
`addAnnotation`, then `updateAnnotation` when multiple are present.

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

The above example specifies a Role, so a mapping for that role must exist in the
current webhook configuration. It uses the `setAnnotation` so the data in the
injector will overwrite the existing contents of the Kubernetes Secret; if
`/test/secret` contains a `username` and `password` but no `domain`, then the
Kubernetes Secret would get the `username` and `password` from the DSV Secret
Data but, the injector will remove the `domain` field.

There are more examples in the `examples` directory. Each one shows how each
annotation works when run against an example with only a username and
private-key in it but no domain.
