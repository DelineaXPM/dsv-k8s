# Release

## First Time Setup

- Run `mage init` to install tooling.
- Install [trunk](https://trunk.io/products/check) (quick install script: `curl https://get.trunk.io -fsSL | bash`)
- Install [aqua](https://aquaproj.github.io/docs/tutorial-basics/quick-start#install-aqua) and have it configured in your path per directions.
  - This will allow faster installs of project tooling by grabbing binaries for your platform more quickly (most of the time release binaries instead of building from source).
- Run `aqua install` for tooling such as changie or others for the project.
  - At this time, it expects you have to Go pre-installed.

## Release Notes

This project uses an different approach to release, driving it from changelog and versioned changelog notes instead of tagging.

> Use [changie](https://changie.dev/guide/quick-start/) quick start for basic review.

### Creating New Notes

- During development, new changes of note get tracked via `changie new`. This can span many pull requests, whatever makes sense as version to ship as changes to users.
- To release the changes into a version, `changie batch <major|minor|patch>` (unless breaking changes occur, you'll want to stick with minor for feature additions, and patch for fixes or non app work.

Keep your summary of changes that users would care about in the `.changes/` files it will create.

### Release

Update [CHANGELOG.md](CHANGELOG.md) by running `changie merge` which will rebuild the changelog file with all the documented notes.

### Format & Lint

- Run `trunk fmt --all; trunk check --all` to finalize run through.
- Push changelog via PR or direct if you have permissions and this will trigger the [release-composite](.github/workflows/release-composite.yml). If any issues, retrigger manually via `gh workflow run release-composite`.
- Release should be published in the [releases](https://github.com/DelineaXPM/dsv-repo-template/releases)
- Edit the release and click "update release" to ensure it publishes to the marketplace. Unfortunately, creating a release doesn't trigger the marketplace release without doing this step. While this can be automated through other actions, I've opted due to time constraints to leave that last step as a manual one.

## FAQ

### What drives the version number for the release?

Changie notes are named like `v1.0.4.md`.
This version number will be used to set the version of the release, so the docs in essence will be the version source of truth.

### Conventional Commit

We use [conventional commit](https://www.conventionalcommits.org/en).
Pull requests must adhere to this to be merged.

Description should be bullet point list or longer-form content to describe anything the title doesn't make clear.
