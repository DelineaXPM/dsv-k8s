# Delinea DevOps Secrets Vault Kubernetes Secret Injector and Syncer

[![Tests](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/tests.yml/badge.svg)](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/tests.yml) [![Docker](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/docker.yml/badge.svg)](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/docker.yml) [![GitHub](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/github.yml/badge.svg)](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/github.yml) [![Red Hat Quay](https://quay.io/repository/delinea/dsv-k8s/status "Red Hat Quay")](https://quay.io/repository/delinea/dsv-k8s)

A [Kubernetes](https://kubernetes.io/)
[Mutating Webhook](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#admission-webhooks)
that injects Secret data from Delinea DevOps Secrets Vault (DSV) into Kubernetes Secrets and a
[CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/)
that subsequently periodically synchronizes them from the source, DSV.
The webhook can be hosted as a pod or as a stand-alone service.
Likewise, the cronjob can run inside or outside the cluster.

The webhook intercepts `CREATE` Secret admissions and then mutates the Secret with data from DSV.
The syncer scans the cluster (or a single namespace) for Secrets that were mutated and,
upon finding a mutated secret,
it compares the version of the DSV Secret with the version it was mutated with and,
if the version in DSV is newer, then the mutation is repeated.

The common configuration consists of one or more Client Credential Tenant mappings.
The credentials are then specified in an [Annotation](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/)
on the Kubernetes Secret to be mutated.
See [below](#use).

The webhook and syncer use the [Golang SDK](https://github.com/DelineaXPM/dsv-sdk-go)
to communicate with the DSV API.

They were tested with [Docker Desktop](https://www.docker.com/products/docker-desktop/)
and [Minikube](https://minikube.sigs.k8s.io/).
They also work on [OpenShift](https://www.redhat.com/en/technologies/cloud-computing/openshift),
[Microk8s](https://microk8s.io/)
and others.

## Injector and Syncer Differences

- Injector: This is a mutating webhook using AdmissionController. This means it operates on the `CREATE` of a Secret, and ensures it modified before finishing the creation of the resource in Kubernetes. This only runs on the creation action triggered by the server.
- Syncer: In contrast, the syncer is a normal cronjob operating on a schedule, checking for any variance in the data between the Secret data between the resource in Kubernetes and the expected value from DSV.

## Which Should I Use?

- Both: If you want a secret to be injected on creation and also synced on your cron schedule then use the Injector and Syncer.
- Injector: If you want the secret to be static despite the change upstream in DSV, and will recreate the secret on any need to upgrade, then the injector. This will reduce the API calls to DSV as well.
- Syncer: If you want the secret value to be updated within the targeted schedule automatically. If this is run by itself without the injector, there can be a lag of up to a minute before the syncer will update the secret. Your application should be able to handle retrying the load of the credential to avoid using the cached credential value that might have been loaded on app start-up in this case.

## Local Development Tooling

- Make: Makefiles provide core automation.
- Mage: Mage is a Go based automation alternative to Make and provides newer functionality for local Kind cluster setup, Go development tooling/linting, and more. Requires Go 17+ and is easily installed via: `go install github.com/magefile/mage@latest`. Run `mage -l` to list all available tasks, and `mage init` to setup developer tooling.
- Pre-Commit: Requires Python3. Included in project, this allows linting and formatting automation before committing, improving the feedback loop.
- Optional:
  - Devcontainer configuration included for VSCode to work with Devcontainers and Codespaces in a pre-built development environment that works on all platforms, and includes nested Docker + ability to run Kind kubernetes clusters without any installing any of those on the Host OS.
  -  Direnv: Default test values are loaded on macOS/Linux based system using [direnv](https://direnv.net/docs/installation.html).
    Run `direnv allow` in the directory to load default env configuration for testing.
  - macOS/Linux: [Trunk.io](https://trunk.io/) to provide linting and formatting on the project. Included in recommended extensions.
    - `trunk install`, `trunk check`, and `trunk fmt` simplifies running checks.

## Configure

The configuration requires a JSON formatted list of Client Credential and Tenant mappings.

```json
{
  "app1": {
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

> *** note ***
> the injector uses the _default_ credentials when mutating a Kubernetes Secret without a _credentialAnnotation_.
> See [below](#use)

## Local

### Run

The injector uses the HTTPS server built-in to the Golang [http](https://pkg.go.dev/net/http)
package to host the Kubernetes Mutating Webhook Webservice.

```bash
$ ./dsv-injector -h
Usage of ./dsv-injector:
  -address string
        the address to listen on, e.g., 'localhost:8080' or ':8443' (default ":18543")
  -cert string
        the path of the public certificate file in PEM format (default "tls/cert.pem")
  -credentials string
        the path of JSON formatted credentials file (default "credentials/config.json")
  -key string
        the path of the private key file in PEM format (default "tls/key.pem")
```

Thus the injector can run "anywhere," but, typically,
the injector runs as a POD in the Kubernetes cluster that uses it.

The syncer is a simple Golang executable.
It typically runs as a Kubernetes CronJob, but it will run outside the cluster.

```bash
$ ./dsv-syncer -h
Usage of ./dsv-syncer:
  -credentials string
        the path of JSON formatted credentials file (default "credentials/config.json")
  -kubeConfig string
        the Kubernetes Client API configuration file; ignored when running in-cluster (default "/home/user/.kube/config")
  -namespace string
        the Kubernetes namespace containing the Secrets to sync; "" (the default) for all namespaces
```

### Build

> *** note ***
> Building the `dsv-injector` image is not required to install it as it is.
> It is available on multiple public registries.

Building the image requires [Docker](https://www.docker.com/) or [Podman](https://podman.io/) and [GNU Make](https://www.gnu.org/software/make/).

To build it, run: `make`.

This will build the injector and syncer as platform binaries and store them in the project root.
It will also build the image (which will build and store its own copy of the binaries) using `$(DOCKER)`.

### Test

The tests expect a few environmental conditions to be met.

> *** note ***
> For more detailed setup see collapsed section below for DSV Test Configuration Setup.

- A valid DSV tenant.
- A secret created with the data format below:
        {
          "data": {
            "password": "admin",
            "username": "admin"
          },
          "version": "0"
        }
- A `configs/credentials.json` to be created manually that contains the client credentials.
- The `configs/credentials.json` credential to be structured like this:

    {
        "app1": {
            "credentials": {
                "clientId": "",
                "clientSecret": ""
            },
            "tenant": "app1"
        }
    }

> *** warning ***
> `app1` is required and using any other will fail test conditions.

<details closed>
<summary>🧪 DSV Test Configuration Setup</summary>

- Using dsv cli (grab from [downloads](https://dsv.secretsvaultcloud.com/downloads) and install with snippet adjusted to version: `$(curl -fSSL https://dsv.secretsvaultcloud.com/downloads/cli/1.35.2/dsv-linux-x64 -o ./dsv && chmod +x dsv-linux-x64 && sudo mv ./dsv-linux-x64 /usr/local/bin/dsv && dsv --version`
- `dsv init` (Use a local user)
- Create the role that will allow creating a client for programmatic access: `dsv role create --name 'k8s' --desc 'test profile for k8s'`
- `dsv secret create --path 'k8s:sync:test' --data '{"password": "admin","username": "admin"}'`
- Create a policy that allows the local user to read the secret, modify this to the correct user/group mapping: `dsv policy create -- actions 'read' --path 'secrets:k8s' --desc 'test access to secret' --resources 'secrets:k8s:<.*>' --subjects 'roles:k8s'`
- Create the client: `dsv client create --role k8s`
- Use those credentials in the structure mentioned above.

</details>

To invoke the tests, run:

```sh
make test
```

Set `$(GO_TEST_FLAGS)` to `-v` to get DEBUG output.

They require a `credentials.json` as either a file or a string.
They also require the path to a secret to test to use.
Use environment variables to specify both:

| Environment Variable       | Default                          | Explanation                                                 |
| -------------------------- | -------------------------------- | ----------------------------------------------------------- |
| `DSV_K8S_TEST_CONFIG`      | _none_                           | Contain a JSON string containing a valid `credentials.json` |
| `DSV_K8S_TEST_CONFIG_FILE` | `../../configs/credentials.json` | The path to a valid `credentials.json`                      |
| `DSV_K8S_TEST_SECRET_PATH` | `/test/secret`                   | The path to the secret to test against in the vault         |

ℹ️ NOTE: `DSV_K8S_TEST_CONFIG` takes precedence over `DSV_K8S_TEST_CONFIG_FILE`

For example:

```sh
DSV_K8S_TEST_CONFIG='{"app1":{"credentials":{"clientId":"93d866d4-635f-4d4e-9ce3-0ef7f879f319","clientSecret":"xxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxxxx"},"tenant":"mytenant"}}' \
DSV_K8S_TEST_SECRET_PATH=my:test:secret \
make test GO_TEST_FLAGS='-v'
```

To remove the binaries and Docker image so that the next build is from scratch, run:

```sh
make clean
```

For Go development, another option is to run gotestsum (installed automatically with `mage init`) with a filewatch option to get regular test output:

```shell
gotestsum --format dots-v2 --watch ./... -- -v
```

## Install

Installation requires [Helm](https://helm.sh).
There are two separate charts for the injector and the syncer.
The `Makefile` demonstrates a typical installation of both.

The dsv-injector chart imports `credentials.json` from the filesystem and stores it in a Kubernetes Secret.
The dsv-syncer chart refers to that Secret instead of creating its own.

The Helm `values.yaml` file `image.repository` is `quay.io/delinea/dsv-k8s`:

```yaml
image:
  repository: quay.io/delinea/dsv-k8s
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: ""
```

That means, by default, `make install` will pull from Red Hat Quay.

```sh
make install
```

However,
the `Makefile` contains an `install-image` target that configures Helm to use the image built with `make image`:

```sh
make install-image
```

`make uninstall` uninstalls the Helm Charts.

### Docker Desktop

To install the locally built image into Docker Desktop, run:

```sh
make install-image
```

ℹ️ NOTE: Kubernetes must be [enabled](https://docs.docker.com/desktop/kubernetes/)
for this to work.

### Remote Cluster

Deploying the locally built image into a remote Cluster will require a container registry.
The container registry must be accessible from the build and the cluster by the same DNS name.
`make` will run the `release` target, which will push the image into the container registry,
`install-cluster` will cause the cluster to pull it from there.

```sh
make install-cluster REGISTRY=internal.example.com:5000
```

### Minikube

#### Docker driver

To deploy to Minikube running on the [Docker driver](https://minikube.sigs.k8s.io/docs/drivers/docker/),
run `eval $(minikube docker-env)` so that the environment shares Minikube's docker context,
then follow the [Docker Desktop](#docker-desktop)
instructions.

#### VM driver

To deploy to Minikube set-up with the VM driver, e.g., Linux [kvm2](https://minikube.sigs.k8s.io/docs/drivers/kvm2/)
or Microsoft [Hyper-V](https://minikube.sigs.k8s.io/docs/drivers/hyperv/),
enable the Minikube built-in registry and use it to make the image available to the Minikube VM:

```sh
minikube addons enable registry
```

❗NOTE: run Minikube [tunnel](https://minikube.sigs.k8s.io/docs/commands/tunnel/)
in a separate terminal to make the registry service available to the host.

```sh
minikube tunnel
```

_It will run continuously, and stopping it will render the registry inaccessible._

Next, get the _host:port_ of the registry:

```sh
kubectl get -n kube-system service registry -o jsonpath="{.spec.clusterIP}{':'}{.spec.ports[0].port}"
```

Finally, follow the [Remote Cluster](#remote-cluster)
instructions using it as `$(REGISTRY)`

### Host (for debugging)

Per above, typically, the injector runs as a POD in the cluster but running it on the host makes debugging easier.

```sh
make install-host EXTERNAL_NAME=laptop.mywifi.net CA_BUNDLE=$(cat /path/to/ca.crt | base64 -w0 -)
```

For it to work:

- The certificate that the injector presents must validate against the `$(CA_BUNDLE)`.
- The certificate must also have a Subject Alternative Name for `$(INJECTOR_NAME).$(NAMESPACE).svc`.
  By default that's `dsv-injector.dsv.svc`.

- The `$(EXTERNAL_NAME)` is a required argument, and the name itself must be resolvable _inside_ the cluster.
__localhost will not work__.

If the `$(CA_BUNDLE)` is argument is omitted, `make` will attempt to extract it from `kubectl config`:

```make
install-host: CA_BUNDLE_KUBE_CONFIG_INDEX = 0
install-host: CA_BUNDLE_JSON_PATH = {.clusters[$(CA_BUNDLE_KUBE_CONFIG_INDEX)].cluster.certificate-authority-data}
install-host: CA_BUNDLE=$(shell $(KUBECTL) config view --raw -o jsonpath='$(CA_BUNDLE_JSON_PATH)' | tr -d '"')

```

which will make:

```sh
kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' | tr -d '"'
```

Optionally set `$(CA_BUNDLE_KUBE_CONFIG_INDEX)` to use `1`, to use the second cluster in your configuration,
`2` for the third and so on.

ℹ️ All this assumes that the injector uses a certificate signed by the cluster CA.
There are several options like [cert-manager](https://cert-manager.io/)
for getting cluster-signed certs, however,
this simple [bash script](https://gist.github.com/amigus/b4e6e642f88e756be1996e44a1c35349)
will request and grant a suitable certificate from the cluster using cURL and OpenSSL.
To use it:

```sh
get_k8s_cert.sh -n dsv-injector -N dsv
```

Now run it:

```sh
./dsv-injector -cert ./dsv-injector.pem -key ./dsv-injector.key -credentials ./configs/credentials.json -address :8543
```

## Use

Once the injector is available in the Kubernetes cluster,
and the webhook is in place,
any correctly annotated Kubernetes Secrets are modified on create and update.

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

The `setAnnotation`, `addAnnotation` and `updateAnnotation`,
must contain the path to the DSV Secret that the injector will use to mutate the Kubernetes Secret.

- `addAnnotation` adds missing fields without overwriting or removing existing fields.
- `updateAnnotation` adds and overwrites existing fields but does not remove fields.
- `setAnnotation` overwrites fields and removes fields that do not exist in the DSV Secret.

NOTE: A Kubernetes Secret should specify only one of the "add," "update,"
or "set" annotations. The order of precedence is `setAnnotation`,
then `addAnnotation`, then `updateAnnotation` when multiple are present.

### Examples

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: example-secret
  annotations:
    dsv.delinea.com/credentials: app1
    dsv.delinea.com/set-secret: /test/secret
type: Opaque
data:
  username: dW5tb2RpZmllZC11c2VybmFtZQ==
  domain: dW5tb2RpZmllZC1kb21haW4=
  password: dW5tb2RpZmllZC1wYXNzd29yZA==
```

The above example specifies credentials,
so a mapping for those credentials must exist in the current webhook configuration.
It uses the `setAnnotation`,
so the data in the injector will overwrite the existing contents of the Kubernetes Secret;
if `/test/secret` contains a `username` and `password` but no `domain`,
then the Kubernetes Secret would get the `username` and `password` from the DSV Secret Data but,
the injector will remove the `domain` field.

There are more examples in the `examples` directory.
They show how the different annotations work.
