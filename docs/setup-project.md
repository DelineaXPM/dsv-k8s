# Setup Project

## Compatibility

Linux & MacOS is supported out of the box for local development.
Windows with WSL2 should also work fine.

While the majority of this is cross-platform, the automatically linting and some other commands are only compatible with Linux/MacOS.

## Overview

- Mage: Mage is a Go based automation alternative to Make and provides newer functionality for local Kind cluster setup, Go development tooling/linting, and more.
  Use [aqua](#aqua) to automatically install.
- Run `mage` to list all available tasks, and `mage init` to setup developer tooling.
  Get more detail on a task, if it's available by running `mage -h init`.

- Make: Makefiles provide core automation from the original project.
  - This has slowly been phased out for the more robust Mage tasks.
    Most of your usage won't touch Make.
    The only usage that _might_ use this is the local cert based setup for debugger.
    This hasn't been used in years by current maintainers, as log streaming from mage/tilt UI have been primary method.
  - For anything other than running the debugger, use `mage` commands which have been regularly improved and tested and move any new automation to the magefiles.

## Initial Setup

## Aqua

This tool will ensure all the core development tools, including Go, are installed and setup without needing to run `apt` or other package managers.

Install [aqua](https://aquaproj.github.io/docs/tutorial-basics/quick-start#install-aqua) and have it configured in your path per directions.

Run `aqua install` for tooling such as changie or others for the project.

Ensure your profile has this in it:

```shell
export PATH="${AQUA_ROOT_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/aquaproj-aqua}/bin:$PATH" # for those using aqua this will ensure it's in the path with all tools if loading from home
```

## When Using Without Devcontainer/Codespaces

- Install Aqua
  - Alternative: Manually ensure Go is installed.
- Run `aqua policy allow` to allow the custom `go install` package to run.
- Run `mage init` to install tooling.
  - Done automatically by Mage -> Install [trunk](https://trunk.io/products/check) (quick install script: `curl https://get.trunk.io -fsSL | bash`)
  - This will allow faster installs of project tooling by grabbing binaries for your platform more quickly (most of the time release binaries instead of building from source).

> If you get an error with a go installation, just try once more as aqua installs in parallel and might not have finished installing Go before trying to run `go install` first.

## Direnv

This loads environment variables for the project automatically.

You should hook into your shell, for example with zsh: `eval "$(direnv hook zsh)"`.

Other shells are supported, but this project is only tested with zsh.

> [Hook Into Your Shell](https://direnv.net/docs/hook.html)

Direnv: Default test values are loaded on macOS/Linux based system using [direnv](https://direnv.net/docs/installation.html).

Run `direnv allow` in the directory to load default env configuration for testing.
