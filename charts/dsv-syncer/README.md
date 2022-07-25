# dsv-syncer

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

A Helm chart for the Delinea DevOps Secrets Vault (DSV) Kubernetes Synchronizer Job.

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Adam Migus |  |  |
| Sheldon Hull |  |  |
| Delinea DSV Team |  |  |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| cronJobSchedule | string | `"* * * * *"` | cronJobSchedule controls when the syncer runs; five asterisks means "every minute". See [cronjob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#cron-schedule-syntax) @default - every minute, ie '* * * * *' |
| dsvInjectorCredentialsSecretName | string | `"dsv-injector-credentials"` | dsvInjectorCredentialsSecretName is the name of thecredentialsJson secret from the dsv-injector |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"Always"` |  |
| image.repository | string | `"quay.io/delinea/dsv-k8s"` |  |
| image.tag | string | `""` |  |
| imagePullSecrets | list | `[]` |  |
| nameOverride | string | `""` |  |
| podAnnotations | object | `{}` | default annotations to add @default - Adds `dsv-filter-name` to simplify log selector streaming |
| podSecurityContext | object | `{}` |  |
| replicaCount | int | `1` | replicaCount @default - 1 |
| resources | object | No default values, user must specify to set resource limits. | We usually recommend not to specify default resources and to leave this as a conscious choice for the user. This also increases chances charts run on environments with little resources, such as Minikube. If you do want to specify resources, uncomment the following lines, adjust them as necessary, and remove the curly braces after 'resources:'. |
| securityContext | object | `{}` |  |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account @default - Adds `dsv-filter-name` to simplify log selector streaming |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created @default - true |
| serviceAccount.name | string | `""` | If not set and create is true, a name is generated using the fullname template |
