# Devcontainer Based Setup

- [Devcontainer Based Setup](#devcontainer-based-setup)
  - [Prerequisites](#prerequisites)
  - [I'm starting from scratch](#im-starting-from-scratch)
    - [Windows](#windows)
    - [MacOS](#macos)
    - [Linux](#linux)
    - [After You've Setup VSCode](#after-youve-setup-vscode)
  - [I already use devcontainers](#i-already-use-devcontainers)
  - [Spin It Up](#spin-it-up)
    - [After Devcontainer Loads](#after-devcontainer-loads)
  - [Working With Kind & Stack Locally](#working-with-kind--stack-locally)
    - [Kind](#kind)
  - [Troubleshooting](#troubleshooting)
    - [Mismatch With Checksum for Go Modules](#mismatch-with-checksum-for-go-modules)
    - [Connecting to Services Outside of devcontainer](#connecting-to-services-outside-of-devcontainer)

## Prerequisites

- Docker
- Visual Studio Code

## I'm starting from scratch

> **_NOTE_**
> Docker is left out of these directions, just install that from [Docker Desktop](https://www.docker.com/products/docker-desktop/) site.

### Windows

- [Install chocolatey (package manager for Windows)](https://chocolatey.org/install#individual) (provides single line command to run).
- Run `choco install vscode -y`

### MacOS

- [Homebrew](https://brew.sh/)

- Run `brew install visual-studio-code`

### Linux

- You'll have to install the apps manually.

### After You've Setup VSCode

Run `code --install-extension ms-vscode-remote.remote-containers`

- For supporting Codespaces: `code --install-extension GitHub.codespaces`

## I already use devcontainers

- Ensure you've got Remote Containers or Codespace extension installed as mentioned in directions above and you'll be good to start.

## Spin It Up

> **_NOTE_**
>
> 🐎 PERFORMANCE TIP: Using the directions provided for named container volume will optimize performance over trying to just "open in container" as there is no mounting files to your local filesystem.

Use command pallet with vscode (Control+Shift+P or F1) and type to find the command `Remote Containers: Clone Repository in Named Container`.

- Put the git clone url in, for example: `https://github.com/DelineaXPM/dsv-k8s.git`
- Name the volume and directory both dsv-k8s or whatever you prefer.

> **_NOTE_**
> This is a large development image (10GB). The first time you run this it will take a while. However, after this first run, rebuilding the container to start over should be minimal time, as you'll have the majority of Docker image cached locally.

This includes (for updated info just look at dockerfile):

- Embedded docker
- Embedded Kind/Minikube (kubernetes)
- Go
- Dotnet
- Python
- Node
- Go tools for linting, formatting, and testing.
- Extensions for VSCode defined in `.devcontainers`, such as Go, Kubernetes & Docker, and some others.
- Initial placeholder `.zshrc` file included to help initialize usage of `direnv` for automatically loading default `.envrc` which contains local developement default environment variables.

### After Devcontainer Loads

1. Accept "Install Recommended Extensions" from popup, to automatically get all the preset tools, such as Kubernetes, Go and others setup.
1. Open a new `zsh-login` terminal and allow the automatic setup to finish, as this will ensure all other required tools are setup.
   - Make sure to run `direnv allow` as it prompts you, to ensure all project and your personal environment variables (optional).
1. Make sure Go 1.19 is the correct version running with `go version`.
   1. If it's not, run `sudo .devcontainer/library-scripts/go-debian.sh "1.19"`
1. Run setup task:
   - Using CLI: Run `mage init`

## Working With Kind & Stack Locally

> **_NOTE_**
> For any tasks get more help with `-h`, for example, run `mage -h k8s:init`

### Kind

For local development, Mage tasks have been created to automate most of the setup and usage for local testing.

- Ensure your local `configs/credentials.json` exists.
- run `mage job:init` to setup a local kind cluster, initial local copies of the helm chart and kubernetes manifest files.
- Modify the `.cache/dsv-injector/values.yaml` with the embedded credentials.json contents matching your `configs/credentials.json`.
- Modify the `.cache/manifests/*.yaml` files to match the credentials you want to test against.
- To deploy (or redeploy after changes) all the helm charts and kuberenetes manifests run `mage job:redeploy`.

## Troubleshooting

### Mismatch With Checksum for Go Modules

- Run `go clean -modcache && go mod tidy`.

### Connecting to Services Outside of devcontainer

You are in an isolated, self-contained Docker setup.
The ports internally aren't the same as externally in your host OS.
If the port forward isn't discovered automatically, enable it yourself, by using the port forward tab (next to the terminal tab).

1. You should see a port forward once the services are up (next to the terminal button in the bottom pane).
   1. If the click to open url doesn't work, try accessing the path manually, and ensure it is `https`.
      Example: `https://127.0.0.1:9999`

You can choose the external port to access, or even click on it in the tab and it will open in your host for you.
