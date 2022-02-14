# Delinea DevOps Secrets Vault Kubernetes Secret Injector

![Docker](https://github.com/thycotic/dsv-k8s/workflows/Docker/badge.svg)
![GitHub Package Registry](https://github.com/thycotic/dsv-k8s/workflows/GitHub%20Package%20Registry/badge.svg)
![Red Hat Quay](https://github.com/thycotic/dsv-k8s/workflows/Red%20Hat%20Quay/badge.svg)

A [Kubernetes](https://kubernetes.io/) [Mutating Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#admission-webhooks)
that injects Secret data from Delinea DevOps Secrets Vault (DSV) into
Kubernetes (k8s) cluster Secrets. The webhook is made available to the
Kubernetes cluster as the `dsv-injector` which can be hosted in k8s or as a
stand-alone service.

The webhook intercepts `CREATE` and `UPDATE` Secret admissions and supplements
or overwrites the Secret data with Secret data from DSV. The webhook
configuration is one or more _role_ to Client Credential Tenant mappings.
The webhook updates k8s Secrets based on annotations (see below).

The webhook uses the [DSV Go SDK](https://github.com/thycotic/dsv-sdk-go) to
communicate with DSV.

It was built and tested with [Minikube](https://minikube.sigs.k8s.io/).
It was also tested with [Minishift](https://docs.okd.io/3.11/minishift/index.html).

## Configure

The webhook requires a JSON formatted list of _role_ to Client Credential and Tenant mappings.
The _role_ is a simple name that does not relate to DSV or Kubernetes Roles per se.
Declaring the _role_ annotation selects which credentials to use to get the DSV Secret.
Using the name of the DSV Role used to generate the credentials is good practice.

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

The injector uses the _default_ Role when it mutates a k8s _Secret_ that does not have an explicit Role annotation (see below).

## Run

The `Makefile` demonstrates a typical installation via [Helm](https://helm.sh/).
It provides the CA certificate bundle `$(CA_BUNDLE)` that the cluster uses to authenticate the webhook to the Helm Chart.
It provisions the certificate and associated key using `get_cert.sh` and submits them as a Kubernetes Secret.
It also provides the `roles.json` file via a Kubernetes Secret.

The `thycotic/dsv-injector` image contains the `dsv-injector-svc` executable, however,
the certificate, the key, and the `roles.json` file should be mounted into the container at runtime.

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

### Certificate

To provision a certificate and key from your Kubernetes cluster,
ensure that [openssl](https://www.openssl.org/) is available the use `scripts/get_cert.sh`.

```bash
$ sh scripts/get_cert.sh
Usage: get_cert.sh -n NAME [OPTIONS]...

        -n, -name, --name NAME
                Maps to the host portion of the FQDN that is the subject of the
                certificate; also the basename of the certificate and key files.
        -d, -directory, --directory=DIRECTORY
                The location of the resulting certificate and private-key. The
                default is '.'
        -N, -namespace, --namespace=NAMESPACE
                Represents the Kubernetes cluster Namespace and maps to the
                domain of the FQDN that is the subject of the certificate.
                the default is 'default'
        -b, -bits, --bits=BITS
                the RSA key size in bits; default is 2048
```

### The CA Certificate

The `Makefile` assumes that it exists as `${HOME}/.minikube/ca.crt`.
By default, it is at `${HOME}/.minishift/ca.pem` for Minishift.

To get the certificate from a running cluster, run:

```shell
kubectl config view -o jsonpath='{.clusters[*].cluster.certificate-authority}'
```

If that returns `null`, then the certificate is embedded in the cluster configuration.
In that case, set `$(CA_BUNDLE)` to the output of:

```shell
kubectl config view -o jsonpath='{.clusters[*].cluster.certificate-authority-data}' | tr -d '"'
```

## Build

Building the `dsv-injector` image requires [Docker](https://www.docker.com/) or
[Podman](https://podman.io/).
To build it, run:

```sh
make image
```

### Minikube and Minishift

Remember to run `eval $(minikube docker-env)` in the shell to push the image to Minikube's Docker daemon.
Likewise for Minishift but run `eval $(minishift docker-env)` instead.

## Use

Once the `dsv-injector` is up and available to the Kubernetes cluster, and the
[MutatingAdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook) is configured to call it, any
appropriately annotated k8s Secrets will be modified by it whenever they are
created or updated.

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

The `setAnnotation`, `addAnnotation` and `updateAnnotation` contain the path to
the DSV Secret that will be used to modified the k8s Secret being admitted.

* `addAnnotation` adds missing fields without overwriting or removing existing fields.
* `updateAnnotation` adds and overwrites existing fields but does not remove fields.
* `setAnnotation` overwrites fields and removes fields that do not exist in the DSV Secret.

Only one of these should be specified on any given k8s Secret, however, if more
than one are defined then the order of precedence is `setAnnotation` then
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

There are more examples in the `examples` directory. Each one will show
how each annotation works when run against an example with only a username and
password in it.

```shell
$ thy secret read /test/secret -f .data
{
  "password": "alongpassword",
  "username": "someuser"
}
```
