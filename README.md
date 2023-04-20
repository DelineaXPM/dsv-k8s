# Delinea DevOps Secrets Vault Kubernetes Secret Injector and Syncer

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->

[![All Contributors](https://img.shields.io/badge/all_contributors-7-orange.svg?style=flat-square)](#contributors)

<!-- ALL-CONTRIBUTORS-BADGE:END -->

![Docker Pulls](https://img.shields.io/docker/pulls/delineaxpm/dsv-k8s?style=for-the-badge)

![Docker Image Version (latest semver)](https://img.shields.io/docker/v/delineaxpm/dsv-k8s?style=for-the-badge)

[![Tests](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/test.yml/badge.svg)](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/test.yml)

[![Release](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/release.yml/badge.svg)](https://github.com/DelineaXPM/dsv-k8s/actions/workflows/release.yml)

[![Red Hat Quay](https://quay.io/repository/delinea/dsv-k8s/status 'Red Hat Quay')](https://quay.io/repository/delinea/dsv-k8s)

The DSV Kubernetes Injector and Syncer are components for [Kubernetes][kubernetes].
The [Mutating Webhook][mutating-webhook] injects Secret data from the Delinea DevOps Secrets Vault (DSV) into Kubernetes Secrets, and a CronJob periodically synchronizes them.
The webhook can run as a pod or a stand-alone service, and the cronjob can run inside or outside the cluster.

- `dsv-injector`: Intercepts `CREATE` Secret admissions and then mutates the Secret with data from DSV.
- `dsv-syncer`: The syncer scans the cluster (or a single namespace) for Secrets that were mutated, compares the secret version, and updates if the secret has changed versions.

The common configuration consists of one or more Client Credential Tenant mappings.
The credentials are then specified in an [Annotation][annotation] on the Kubernetes Secret to be mutated.

The webhook and syncer use the [Golang SDK][dsv-go-sdk] to communicate with the DSV API.
They were tested with [Docker Desktop][docker-desktop] and [Minikube][minikube].
They also work on [OpenShift][openshift], [Microk8s][microk8s] and others.

## Contents

- [Delinea DevOps Secrets Vault Kubernetes Secret Injector and Syncer](#delinea-devops-secrets-vault-kubernetes-secret-injector-and-syncer)
  - [Contents](#contents)
  - [Supporting Docs](#supporting-docs)
  - [Injector \& Syncer Differences](#injector--syncer-differences)
    - [Which Should I Use?](#which-should-i-use)
  - [Quick Start](#quick-start)
    - [Build](#build)
  - [Test](#test)
  - [Reference Mage Tasks](#reference-mage-tasks)
  - [Contributors](#contributors)

## Supporting Docs

The [docs](docs/) directory has supporting documentation that goes into more detail on the developer workflows, test setup, configuration, helm install commands, and more.

## Injector & Syncer Differences

- Injector: This is a mutating webhook using AdmissionController.
  This means it operates on the `CREATE` of a Secret, and ensures it modified before finishing the creation of the resource in Kubernetes.
  This only runs on the creation action triggered by the server.
- Syncer: In contrast, the syncer is a normal cronjob operating on a schedule, checking for any variance in the data
  between the Secret data between the resource in Kubernetes and the expected value from DSV.

### Which Should I Use?

- Both: If you want a secret to be injected on creation and also synced on your cron schedule then use the Injector and Syncer.
- Injector: If you want the secret to be static despite the change upstream in DSV, and will recreate the secret on any need to upgrade, then the injector.
  This will reduce the API calls to DSV as well.
- Syncer: If you want the secret value to be updated within the targeted schedule automatically.
  If this is run by itself without the injector, there can be a lag of up to a minute before the syncer will update the secret.
  Your application should be able to handle retrying the load of the credential to avoid using the cached credential value that might have been loaded on app start-up in this case.

## Quick Start

Since there's a mix of users for this repo, here's where to go for getting up and running as quickly as possible.

| Who                                                               | Where do I start?                                           |
| ----------------------------------------------------------------- | ----------------------------------------------------------- |
| 👉 I just want to install the helm charts against my own cluster. | Clone, and use `helm install` against the charts directory. |
| 👉 I'm a contributor/developer and want to test/build locally     | Use the [setup-developer](docs/setup-developer.md) guide.   |
| 👉 I'm a contributor and need to create a release.                | Use the [release](docs/release.md) guide.                   |

### Build

<img src="docs/assets/random-dont-need-to-install.svg">

To build run: `mage init build`.
For more detailed directions on local development (such as Mage), see [setup-developer](docs/setup-developer.md)

## Test

See details in [local-testing](docs/local-testing.md)

## Reference Mage Tasks

> Manually updated, for most recent Mage tasks, run `mage -l`.

| Target           | Description                                                                                                                          |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| build            | 🔨 Build builds the project for the current platform.                                                                                |
| buildAll         | 🔨 BuildAll builds all the binaries defined in the project, for all platforms.                                                       |
| clean            | up after yourself.                                                                                                                   |
| go:doctor        | 🏥 Doctor will provide config details.                                                                                               |
| go:fix           | 🔎 Run golangci-lint and apply any auto-fix.                                                                                         |
| go:fmt           | ✨ Fmt runs gofumpt.                                                                                                                 |
| go:init          | ⚙️ Init runs all required steps to use this package.                                                                                 |
| go:lint          | 🔎 Run golangci-lint without fixing.                                                                                                 |
| go:lintConfig    | 🏥 LintConfig will return output of golangci-lint config.                                                                            |
| go:test          | 🧪 Run go test.                                                                                                                      |
| go:testSum       | 🧪 Run gotestsum (Params: Path just like you pass to go test, ie ./..., pkg/, etc ).                                                 |
| go:tidy          | 🧹 Tidy tidies.                                                                                                                      |
| go:wrap          | ✨ Wrap runs golines powered by gofumpt.                                                                                             |
| helm:docs        | generates helm documentation using `helm-doc` tool.                                                                                  |
| helm:init        | ⚙️ Init sets up the required files to allow for local editing/overriding from CacheDirectory.                                        |
| helm:install     | 🚀 Install uses Helm to install the chart.                                                                                           |
| helm:lint        | 🔍 Lint uses Helm to lint the chart for issues.                                                                                      |
| helm:render      | 💾 Render uses Helm to output rendered yaml for testing helm integration.                                                            |
| helm:uninstall   | 🚀 Uninstall uses Helm to uninstall the chart.                                                                                       |
| init             | runs multiple tasks to initialize all the requirements for running a project for a new contributor.                                  |
| installTrunk     | installs trunk.io tooling if it isn't already found.                                                                                 |
| job:init         | runs the setup tasks to initialize the local resources and files, without trying to apply yet.                                       |
| job:redeploy     | removes kubernetes resources and helm charts and then redeploys with log streaming by default.                                       |
| k8s:apply        | applies a kubernetes manifest.                                                                                                       |
| k8s:delete       | Apply applies a kubernetes manifest.                                                                                                 |
| k8s:init         | copies the k8 yaml manifest files from the examples directory to the cache directory for editing and linking in integration testing. |
| k8s:logs         | streams logs until canceled for the dsv syncing jobs, based on the label `dsv.delinea.com: syncer`.                                  |
| kind:destroy     | 🗑️ Destroy tears down the Kind cluster.                                                                                              |
| kind:init        | ➕ Create creates a new Kind cluster and populates a kubeconfig in cachedirectory.                                                   |
| minikube:destroy | 🗑️ Destroy tears down the Kind cluster.                                                                                              |
| minikube:init    | ➕ Create creates a new Minikube cluster and populates a kubeconfig in cachedirectory.                                               |
| release          | 🔨 Release generates a release for the current platform.                                                                             |
| trunkInit        | ensures the required runtimes are installed.                                                                                         |

## Contributors

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center"><a href="https://mig.us/adam"><img src="https://avatars.githubusercontent.com/u/119477?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Adam C.
Migus</b></sub></a><br /><a href="https://github.com/DelineaXPM/dsv-k8s/commits?author=amigus" title="Code">💻</a> <a href="https://github.com/DelineaXPM/dsv-k8s/commits?author=amigus" title="Documentation">📖</a> <a href="https://github.com/DelineaXPM/dsv-k8s/commits?author=amigus" title="Tests">⚠️</a></td>
      <td align="center"><a href="https://www.sheldonhull.com"><img src="https://avatars.githubusercontent.com/u/3526320?v=4?s=100" width="100px;" alt=""/><br /><sub><b>sheldonhull</b></sub></a><br /><a href="https://github.com/DelineaXPM/dsv-k8s/commits?author=sheldonhull" title="Code">💻</a> <a href="https://github.com/DelineaXPM/dsv-k8s/commits?author=sheldonhull" title="Documentation">📖</a> <a href="https://github.com/DelineaXPM/dsv-k8s/commits?author=sheldonhull" title="Tests">⚠️</a></td>
      <td align="center"><a href="https://github.com/hansboder"><img src="https://avatars.githubusercontent.com/u/36736535?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Hans Boder</b></sub></a><br /><a href="https://github.com/DelineaXPM/dsv-k8s/issues?q=author%3Ahansboder" title="Bug reports">🐛</a></td>
      <td align="center"><a href="https://github.com/tylerezimmerman"><img src="https://avatars.githubusercontent.com/u/100804646?v=4?s=100" width="100px;" alt=""/><br /><sub><b>tylerezimmerman</b></sub></a><br /><a href="#maintenance-tylerezimmerman" title="Maintenance">🚧</a></td>
      <td align="center"><a href="https://github.com/delineaKrehl"><img src="https://avatars.githubusercontent.com/u/105234788?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Tim Krehl</b></sub></a><br /><a href="#maintenance-delineaKrehl" title="Maintenance">🚧</a></td>
      <td align="center"><a href="http://endlesstrax.com"><img src="https://avatars.githubusercontent.com/u/17141891?v=4?s=100" width="100px;" alt=""/><br /><sub><b>Ricky White</b></sub></a><br /><a href="#maintenance-EndlessTrax" title="Maintenance">🚧</a></td>
      <td align="center"><a href="https://github.com/forced-request"><img src="https://avatars.githubusercontent.com/u/961246?v=4?s=100" width="100px;" alt=""/><br /><sub><b>John Poulin</b></sub></a><br /><a href="#maintenance-forced-request" title="Maintenance">🚧</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification.
Contributions of any kind welcome!

[kubernetes]: https://kubernetes.io/
[mutating-webhook]: https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#admission-webhooks
[annotation]: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
[dsv-go-sdk]: https://github.com/DelineaXPM/dsv-sdk-go
[docker-desktop]: https://www.docker.com/products/docker-desktop/
[minikube]: https://minikube.sigs.k8s.io/
[openshift]: https://www.redhat.com/en/technologies/cloud-computing/openshift
[microk8s]: https://microk8s.io/
