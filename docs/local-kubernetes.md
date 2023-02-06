# Local Kubernetes

For easier development workflow, the project has prebuilt tasks for both minikube and kind kubernetes tools.

As of 2023-02, the default behavior is to use minikube since it handles updating the kubeconfig locally a little more consistently in the tests run.

## Working With Kubernetes & Stack Locally

> **_NOTE_**
> For any tasks get more help with `-h`, for example, run `mage -h k8s:init`

For local development, Mage tasks have been created to automate most of the setup and usage for local testing.

- Ensure your local `configs/credentials.json` exists.
- run `mage job:init` to setup a local k8s cluster, initial local copies of the helm chart and kubernetes manifest files.
- Modify the `.cache/dsv-injector/values.yaml` with the embedded credentials.json contents matching your `configs/credentials.json`.
- Modify the `.cache/manifests/*.yaml` files to match the credentials you want to test against.
- To deploy (or redeploy after changes) all the helm charts and kuberenetes manifests run `mage job:redeploy`.

## Using Minikube With VM Driver

<details>
<summary>ℹ️ Using Minikube With VM Driver</summary>

To deploy to Minikube set-up with the VM driver, e.g., Linux [kvm2](https://minikube.sigs.k8s.io/docs/drivers/kvm2/)
or Microsoft [Hyper-V](https://minikube.sigs.k8s.io/docs/drivers/hyperv/),
enable the Minikube built-in registry and use it to make the image available to the Minikube VM:

```shell
minikube addons enable registry
```

❗NOTE: run Minikube [tunnel](https://minikube.sigs.k8s.io/docs/commands/tunnel/)
in a separate terminal to make the registry service available to the host.

```shell
minikube tunnel
```

_It will run continuously, and stopping it will render the registry inaccessible._

Next, get the _host:port_ of the registry:

```shell
kubectl get -n kube-system service registry -o jsonpath="{.spec.clusterIP}{':'}{.spec.ports[0].port}"
```

Finally, follow the [Remote Cluster](#remote-cluster)
instructions using it as `$(REGISTRY)`

</details>
