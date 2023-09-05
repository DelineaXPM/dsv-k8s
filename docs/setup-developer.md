# Setup Developer

> Important: All the core local workflow tasks to build and deploy to minikube are wrapped up in mage tasks
>
> Try `mage` by itself to list.
> Use `mage job:*` tasks to help simplify the process.

1. [Setup developer tooling](setup-project.md)
2. [Create DSV Credentials for Testing](configure.md)
3. [Configure The Manifests](configure.md#update-manifests)
4. Once credentials are configured in `.cache/dsv-injector/values.yaml`
   1. 1st time: `mage job:init`.
   2. 1st time/Anytime You updated Go code: `mage job:rebuildimages`.
   3. Any time you want to redeploy the kubernetes & helm charts to minikube: `mage job:redeploy`.

As always, the source of truth is `mage` so if the task names in the doc don't work, check the CLI for the proper commands.

## Optional

If you are using codespaces, most of the tooling should be ready out of the box as long as you open `zsh` terminal.
Run `tilt up` and then you can invoke much of this (including watch the logs stream) from the terminal.

## Reference

- Optional: [devcontainer/codespaces](devcontainer.md)
- Local Kubernetes Overview: [Kubernetes](local-kubernetes.md)
