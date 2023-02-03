# Setup Developer

## Dive In

- [devcontainer/codespaces](devcontainer.md)
- [Kubernetes](local-kubernetes.md)
- [Project Setup](setup-project.md)

## Other Dev Tools

Use Stern to easily stream cross namespace logs with the `dsv-filter-selector` by running:

Aqua installs this automatically, but if you want to do this manually grab from github releases like this

```shell
$(curl -fSSl https://github.com/wercker/stern/releases/download/1.11.0/stern_linux_amd64 -o ./stern) && sudo chmod +x ./stern && sudo mv ./stern /usr/local/bin
```

> Or use `brew install stern` or aqua.

While `mage k8s:logs` will run this for you, manually you can invoke like this:

```shell
# For all pods in the namespace run
stern --kubeconfig .cache/config --namespace dsv --timestamps .

# For pods with the selector run
stern --kubeconfig .cache/config --namespace dsv --timestamps --selector 'dsv-filter-name in (dsv-syncer, dsv-injector)'
```
