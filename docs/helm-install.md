## Helm Install

Installation of charts into the a cluster requires [Helm](https://helm.sh).

There are two separate charts for the `dsv-injector` and the `dsv-syncer`.

- The `dsv-injector` chart imports `credentials.json` from the filesystem and stores it in a Kubernetes Secret.
- The `dsv-syncer` chart refers to that Secret _instead of creating its own_.

See [configure](configure.md#json-credentials-for-helm-install)

```shell
NAMESPACE='testing'
CREDENTIALS_JSON_FILE='.cache/credentials.json'
IMAGE_REPOSITORY='docker.io/delineaxpm/dsv-k8s'

helm install
     --namespace $NAMESPACE
     --create-namespace \
     --set-file credentialsJson=${CREDENTIALS_JSON_FILE} \
      --set image.repository=${IMAGE_REPOSITORY} \
     dsv-injector ./charts/dsv-injector

helm install
     --namespace $NAMESPACE
     --create-namespace \
     --set-file credentialsJson=${CREDENTIALS_JSON_FILE} \
      --set image.repository=${IMAGE_REPOSITORY} \
     dsv-syncer ./charts/dsv-syncer
```
