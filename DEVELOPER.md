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
  - [Troubleshooting](#troubleshooting)
    - [Mage or Other CLI Tool Not Found](#mage-or-other-cli-tool-not-found)
    - [Mismatch With Checksum for Go Modules](#mismatch-with-checksum-for-go-modules)
    - [Connecting to Services Outside of devcontainer](#connecting-to-services-outside-of-devcontainer)

## Prerequisites

- Docker
- Visual Studio Code

## I'm starting from scratch

> ***NOTE***
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



> ***NOTE***
>
> 🐎 PERFORMANCE TIP: Using the directions provided for named container volume will optimize performance over trying to just "open in container" as there is no mounting files to your local filesystem.

Use command pallet with vscode (Control+Shift+P or F1) and type to find the command `Remote Containers: Clone Repository in Named Container`.

- Put the git clone url in, for example: `https://github.com/DelineaXPM/dsv-k8s.git`
- Name the volume and directory both dsv-k8s or whatever you prefer.

> ***NOTE***
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

> PROTIP: You can run all the commands in this doc with a click using the recommended Command Runner extension (`zvik.command-runner`). You'll get a button above code snippets to run with it installed. If you run into any endless terminal loading loops disable it. I had some issue in the past that I can't reproduce, so it's just an FYI.

1. Open a new `zsh-login` terminal and allow the automatic setup to finish, as this will ensure all other required tools are setup.
    - Make sure to run `direnv allow` as it prompts you, to ensure all project and your personal environment variables (optional).
2. Run setup task:
    - Using CLI: Run `mage init`


## Troubleshooting

### Mage or Other CLI Tool Not Found

If mage command isn't found, just run `go run mage.go init` and it should setup mage if the other tooling failed to.


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

