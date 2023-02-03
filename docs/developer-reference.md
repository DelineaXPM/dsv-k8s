# Devcontainer Based Setup

- [Devcontainer Based Setup](#devcontainer-based-setup)
  - [Troubleshooting](#troubleshooting)
  - [Error With Permissions On Go Directories](#error-with-permissions-on-go-directories)
    - [Mismatch With Checksum for Go Modules](#mismatch-with-checksum-for-go-modules)
    - [Connecting to Services Outside of devcontainer](#connecting-to-services-outside-of-devcontainer)

## Troubleshooting

## Error With Permissions On Go Directories

Clear the directories with `rm -rf /home/vscode/go` and then try `mage init` to redownload packages.

> Known issue: Haven't figured out why this is being set incorrectly yet

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
