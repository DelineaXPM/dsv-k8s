# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html),
and is generated by [Changie](https://github.com/miniscruff/changie).

## v1.2.4 - 2024-10-02

### 🤖 CI & Build

- Improve mage tasks and linting with fixes. Less issues with err output on cleanup now, and also load the docker version of images into minikube setup proactively.

### ⬆️ Dependencies

- Maintenance release due to updated dependencies.

## v1.2.3 - 2024-08-12

### 🤖 CI & Build

- Add a buildName metadata to binary so easy to see if caching issue with container loading. Handle `dev.local/dsv-k8s` as standard image name to better reflect standard approach I've been using. Improve validation checks. Goreleaser upgrade schema and more. Lots of quality of life improvements for dev, and aqua updates.

### 🔨 Refactor

- Improve `values.yml` for the dsv-injector to expose the days till expiration of the self signed cert.
  Include minor doc improvements to this as well to better handle.

## v1.2.2 - 2024-01-15

### ⬆️ Dependencies

- Update dependent libraries and go version. No user facing changes, just continued maintenance for improved security & stability.

## v1.2.1 - 2023-09-05

### 📘 Documentation

- Include detail on providing `tld` in the configuration, allowing `eu` and other TLDs to be used.
- Mention `tilt up` in the initial setup config as viable option.

### 🤖 CI & Build

- Improve mage tasks with secret setup and tear down for better development support and troubleshooting.
- Bump go version in release pipeline to use `1.21` as can include standard library security improvements.
- Remove failing error condition on `mage job:rebuild` to better allow default setup without running local builds, such as just using the published docker image.
  This supports easier demo/test usage by support.

### 🔨 Refactor

- Improve logging with error wrapping and remove deprecated Go `ioutil` usage.

### ⬆️ Dependencies

- Bump tooling such as changie, release, trunk, more security scanners.
- Other dependency bumps such as `golang.org/x/net`.

## v1.2.0 - 2023-04-27

### 🤖 CI & Build

- Improve mage tasks for minikube, including list images, remove/load images.
- Update aqua to the latest version.
- Include aqua configuration tags for Github actions.
- Pin the version examples in the chart install and yaml values file as a better practice for production chart usage.
- Chart versions no longer are independently versioned.
  Instead the docker image, syncing chart, and injector chart are all aligned with the same version number.
  This is automatically kept to the correct version when running `changie merge`.

### 🔨 Refactor

- Improve output to structured json logs.
- Change invocation to use environment variable configuration.
- Allow users to provide config maps dynamically in Values.yaml.
- Note the current version, commit sha, and date of the currently used binary at startup.

### 🔒 Security

- Improve the base image to use Chainguard's static nonroot image.
- Ensure default timeout is set for `ReadHeaderTimeout`.

### 🤖 Development

- Improve the local development experience with Tilt.
- Improve documentation and include stern/kubectl log to file snippets.
- Add local development checks for values files.
  This will catch incorrectly configured values that are difficult to troubleshoot in local development workflows.
- Add missing `--overwrite` tag to `mage minikube:loadimages` task.

## v1.1.6 - 2023-02-06

### 🤖 CI

- Adjust Quay registry to be a target publishing platform instead of building on every single commit to main.
  This aligns Quay registry with the more controlled semver versioning in docker hub instead of publishing under latest only.

### Related

- fixed AB#485565

### Contributors

- [sheldonhull](https://github.com/sheldonhull)

## v1.1.5 - 2023-01-24

### 🔨 Refactor

- Point the helm charts towards docker hub based images, instead of quay, as these are now iterated on with changelog driven release instead of each commit.
  This should reduce frequency of needless version updates.

### 🐛 Bug Fix

- Docker Hub published images did not have the correct path to the injector and syncer, resulting in an invalid entrypoint.
  This is fixed and should now correctly resolve when using the updated helm charts that provide a qualified path.
  For example: `/app/dsv-injector` instead of just saying `dsv-injector` now.
  This is due to using a minimal distroless image and not copying binaries into a path that is assumed to be resolved automatically by `PATH`, such as `/usr/local/bin`.
  Now the path to the binary is explicitly set and should resolve any path resolution issues.

### 🤖 Development

- Bump aqua tooling and include dsv-cli in the project setup.
- Included `CGO_ENABLED=0` to avoid issues with running commands in devcontainers & codespaces.
- Improved mage tasks to support minikube as default to see if this helps with Codespace timeouts being experienced.
- Bumped the docker feature kit to 2.0 as well to attempt to resolve timeouts in devcontainer/codespaces.

### Related

- fixes AB#483421
- [Issue 83 Fixed](https://github.com/DelineaXPM/dsv-k8s/issues/83).
  Thank you @JulianPedro for helping identify this and opening the descriptive issue. 👍

### Contributors

- [sheldonhull](https://github.com/sheldonhull)

## v1.1.4 - 2022-10-11

### Security

Update kubernetes package dependencies.

## v1.1.3 - 2022-10-10

### Added

- Changelog generation triggers docker release instead of every commit

### Fixed

- Resolve dockerhub publishing.

## [v1.1.2] (2022-10-10)

### Added

- feat: Configure Mend for GitHub.com (#26) (2022-08-26)

### Fixed

- fix(mend): 🐛 remove trailing space in json (2022-08-26)
- fix(tests): 🧪 resolve failing test cases and improve local testing automation (#7) (2022-07-27)

### Others

- chore(deps): update trunk-io/trunk-action digest to 22e948f (#59) (2022-10-07)
- chore(deps): update actions/stale digest to 5ebf00e (#58) (2022-10-07)
- chore: add trunk config checks as linting tool (#40) (2022-10-07)
- chore(deps): update amannn/action-semantic-pull-request digest to 505e44b (#46) (2022-09-28)
- docs: add forced-request as a contributor for maintenance (#53) (2022-09-28)
- docs: add EndlessTrax as a contributor for maintenance (#52) (2022-09-28)
- docs: add delineaKrehl as a contributor for maintenance (#51) (2022-09-27)
- docs: add tylerezimmerman as a contributor for maintenance (#50) (2022-09-27)
- docs: add sheldonhull as a contributor for code, doc, and test (#48) (2022-09-27)
- docs: add hansboder as a contributor for bug (#54) (2022-09-27)
- docs: add amigus as a contributor for code, doc, and test (#49) (2022-09-27)
- chore(deps): update actions/stale action to v6 (#43) (2022-09-27)
- chore(deps): update actions/stale digest to 9c1b1c6 (#38) (2022-09-01)
- chore(vscode): add missing linting settings (2022-09-01)
- chore: align github workflows and tooling (#37) (2022-09-01)
- ci(action/docker): 🔨 eliminate github event name push from skipping workflow dispatch (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.3 (#35) (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.2 (#34) (2022-08-29)
- chore(deps): update kubernetes packages to v0.25.0 (#25) (2022-08-26)
- ci(mend): disable issues for mend but still allow pull and renovate execution (2022-08-26)
- chore(deps): update github.com/mattbaird/jsonpatch digest to 098863c (#13) (2022-08-23)
- chore: update maintainers (#24) (2022-08-23)
- chore(devcontainer): 🧰 improvements to experience with setup directions, prebuilt tooling, and reliability (#23) (2022-08-20)
- chore(deps): update kubernetes packages to v0.24.4 (#22) (2022-08-18)
- ci(tests): fix name of tasks (2022-08-05)
- chore(deps): update ⬆️ golang module go to 1.19 (#17) (2022-08-05)
- ci(tests): remove push filter condition (2022-08-05)
- ci(tests): use mage to run tests (2022-08-05)
- ci(tests): 🔨 adjust test to use gotestsum (2022-08-05)
- ci(tests): ➕ require tests to be run as part of pull request (2022-08-05)
- ci: 🤖 add workflow dispatch to docker building actions (#21) (2022-08-03)
- ci(docker): set registry as org secret (#20) (2022-08-01)
- chore: use new env vars for docker hub publishing (#19) (2022-08-01)
- chore(deps): update ⬆️ golang module github.com/pterm/pterm to v0.12.45 (#16) [skip ci] (2022-07-29)
- chore(deps): update ⬆️ golang module github.com/mittwald/go-helm-client to v0.11.3 (#14) [skip ci] (2022-07-28)
- Configure Renovate AB#449139 (#11) (2022-07-28)
- chore(codeowner): assign default codeowners (#12) (2022-07-27)
- chore(helm): ⬆️ bump chart patch version due to minor changes (#8) (2022-07-27)
- Merge pull request #6 from The-Migus-Group/fix-meta (2022-07-12)
- Merge pull request #5 from DelineaXPM/fix-docker (2022-06-28)
- Merge pull request #4 from DelineaXPM/fix-3 (2022-06-28)
- Merge pull request #1 from DelineaXPM/delineaKrehl-DeepRebrand (2022-06-02)
- Fix for #13 that improves injector error handling. (#14) (2022-05-20)
  ]

## [v1.1.1] (2022-10-10)

### Added

- feat: Configure Mend for GitHub.com (#26) (2022-08-26)

### Fixed

- fix(mend): 🐛 remove trailing space in json (2022-08-26)
- fix(tests): 🧪 resolve failing test cases and improve local testing automation (#7) (2022-07-27)

### Others

- chore(deps): update trunk-io/trunk-action digest to 22e948f (#59) (2022-10-07)
- chore(deps): update actions/stale digest to 5ebf00e (#58) (2022-10-07)
- chore: add trunk config checks as linting tool (#40) (2022-10-07)
- chore(deps): update amannn/action-semantic-pull-request digest to 505e44b (#46) (2022-09-28)
- docs: add forced-request as a contributor for maintenance (#53) (2022-09-28)
- docs: add EndlessTrax as a contributor for maintenance (#52) (2022-09-28)
- docs: add delineaKrehl as a contributor for maintenance (#51) (2022-09-27)
- docs: add tylerezimmerman as a contributor for maintenance (#50) (2022-09-27)
- docs: add sheldonhull as a contributor for code, doc, and test (#48) (2022-09-27)
- docs: add hansboder as a contributor for bug (#54) (2022-09-27)
- docs: add amigus as a contributor for code, doc, and test (#49) (2022-09-27)
- chore(deps): update actions/stale action to v6 (#43) (2022-09-27)
- chore(deps): update actions/stale digest to 9c1b1c6 (#38) (2022-09-01)
- chore(vscode): add missing linting settings (2022-09-01)
- chore: align github workflows and tooling (#37) (2022-09-01)
- ci(action/docker): 🔨 eliminate github event name push from skipping workflow dispatch (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.3 (#35) (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.2 (#34) (2022-08-29)
- chore(deps): update kubernetes packages to v0.25.0 (#25) (2022-08-26)
- ci(mend): disable issues for mend but still allow pull and renovate execution (2022-08-26)
- chore(deps): update github.com/mattbaird/jsonpatch digest to 098863c (#13) (2022-08-23)
- chore: update maintainers (#24) (2022-08-23)
- chore(devcontainer): 🧰 improvements to experience with setup directions, prebuilt tooling, and reliability (#23) (2022-08-20)
- chore(deps): update kubernetes packages to v0.24.4 (#22) (2022-08-18)
- ci(tests): fix name of tasks (2022-08-05)
- chore(deps): update ⬆️ golang module go to 1.19 (#17) (2022-08-05)
- ci(tests): remove push filter condition (2022-08-05)
- ci(tests): use mage to run tests (2022-08-05)
- ci(tests): 🔨 adjust test to use gotestsum (2022-08-05)
- ci(tests): ➕ require tests to be run as part of pull request (2022-08-05)
- ci: 🤖 add workflow dispatch to docker building actions (#21) (2022-08-03)
- ci(docker): set registry as org secret (#20) (2022-08-01)
- chore: use new env vars for docker hub publishing (#19) (2022-08-01)
- chore(deps): update ⬆️ golang module github.com/pterm/pterm to v0.12.45 (#16) [skip ci] (2022-07-29)
- chore(deps): update ⬆️ golang module github.com/mittwald/go-helm-client to v0.11.3 (#14) [skip ci] (2022-07-28)
- Configure Renovate AB#449139 (#11) (2022-07-28)
- chore(codeowner): assign default codeowners (#12) (2022-07-27)
- chore(helm): ⬆️ bump chart patch version due to minor changes (#8) (2022-07-27)
- Merge pull request #6 from The-Migus-Group/fix-meta (2022-07-12)
- Merge pull request #5 from DelineaXPM/fix-docker (2022-06-28)
- Merge pull request #4 from DelineaXPM/fix-3 (2022-06-28)
- Merge pull request #1 from DelineaXPM/delineaKrehl-DeepRebrand (2022-06-02)
- Fix for #13 that improves injector error handling. (#14) (2022-05-20)
- Add Testing 📛 (2022-05-06)
- Run go test ./... (2022-05-06)
- Shorten the name of the github action (2022-05-06)
- Fix build badges. 📛 (2022-05-06)
- Secret Synchronization Mechanism #11 (#12) (2022-05-06)
- Merge pull request #10 from thycotic/chart-upgrade-fixes (2022-03-22)
- Eval role logic after patchMode test; Fixes #8 (#9) (2022-03-10)

## [v1.1.0] (2022-10-10)

### Added

- feat: Configure Mend for GitHub.com (#26) (2022-08-26)

### Fixed

- fix(mend): 🐛 remove trailing space in json (2022-08-26)
- fix(tests): 🧪 resolve failing test cases and improve local testing automation (#7) (2022-07-27)

### Others

- chore(deps): update trunk-io/trunk-action digest to 22e948f (#59) (2022-10-07)
- chore(deps): update actions/stale digest to 5ebf00e (#58) (2022-10-07)
- chore: add trunk config checks as linting tool (#40) (2022-10-07)
- chore(deps): update amannn/action-semantic-pull-request digest to 505e44b (#46) (2022-09-28)
- docs: add forced-request as a contributor for maintenance (#53) (2022-09-28)
- docs: add EndlessTrax as a contributor for maintenance (#52) (2022-09-28)
- docs: add delineaKrehl as a contributor for maintenance (#51) (2022-09-27)
- docs: add tylerezimmerman as a contributor for maintenance (#50) (2022-09-27)
- docs: add sheldonhull as a contributor for code, doc, and test (#48) (2022-09-27)
- docs: add hansboder as a contributor for bug (#54) (2022-09-27)
- docs: add amigus as a contributor for code, doc, and test (#49) (2022-09-27)
- chore(deps): update actions/stale action to v6 (#43) (2022-09-27)
- chore(deps): update actions/stale digest to 9c1b1c6 (#38) (2022-09-01)
- chore(vscode): add missing linting settings (2022-09-01)
- chore: align github workflows and tooling (#37) (2022-09-01)
- ci(action/docker): 🔨 eliminate github event name push from skipping workflow dispatch (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.3 (#35) (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.2 (#34) (2022-08-29)
- chore(deps): update kubernetes packages to v0.25.0 (#25) (2022-08-26)
- ci(mend): disable issues for mend but still allow pull and renovate execution (2022-08-26)
- chore(deps): update github.com/mattbaird/jsonpatch digest to 098863c (#13) (2022-08-23)
- chore: update maintainers (#24) (2022-08-23)
- chore(devcontainer): 🧰 improvements to experience with setup directions, prebuilt tooling, and reliability (#23) (2022-08-20)
- chore(deps): update kubernetes packages to v0.24.4 (#22) (2022-08-18)
- ci(tests): fix name of tasks (2022-08-05)
- chore(deps): update ⬆️ golang module go to 1.19 (#17) (2022-08-05)
- ci(tests): remove push filter condition (2022-08-05)
- ci(tests): use mage to run tests (2022-08-05)
- ci(tests): 🔨 adjust test to use gotestsum (2022-08-05)
- ci(tests): ➕ require tests to be run as part of pull request (2022-08-05)
- ci: 🤖 add workflow dispatch to docker building actions (#21) (2022-08-03)
- ci(docker): set registry as org secret (#20) (2022-08-01)
- chore: use new env vars for docker hub publishing (#19) (2022-08-01)
- chore(deps): update ⬆️ golang module github.com/pterm/pterm to v0.12.45 (#16) [skip ci] (2022-07-29)
- chore(deps): update ⬆️ golang module github.com/mittwald/go-helm-client to v0.11.3 (#14) [skip ci] (2022-07-28)
- Configure Renovate AB#449139 (#11) (2022-07-28)
- chore(codeowner): assign default codeowners (#12) (2022-07-27)
- chore(helm): ⬆️ bump chart patch version due to minor changes (#8) (2022-07-27)
- Merge pull request #6 from The-Migus-Group/fix-meta (2022-07-12)
- Merge pull request #5 from DelineaXPM/fix-docker (2022-06-28)
- Merge pull request #4 from DelineaXPM/fix-3 (2022-06-28)
- Merge pull request #1 from DelineaXPM/delineaKrehl-DeepRebrand (2022-06-02)
- Fix for #13 that improves injector error handling. (#14) (2022-05-20)
- Add Testing 📛 (2022-05-06)
- Run go test ./... (2022-05-06)
- Shorten the name of the github action (2022-05-06)
- Fix build badges. 📛 (2022-05-06)
- Secret Synchronization Mechanism #11 (#12) (2022-05-06)
- Merge pull request #10 from thycotic/chart-upgrade-fixes (2022-03-22)
- Eval role logic after patchMode test; Fixes #8 (#9) (2022-03-10)
- Comments + all = install-image (2022-03-01)
- Publish to Minikube by default. 🪄 (2022-03-01)

## [v1.0.1] (2022-10-10)

### Added

- feat: Configure Mend for GitHub.com (#26) (2022-08-26)

### Fixed

- fix(mend): 🐛 remove trailing space in json (2022-08-26)
- fix(tests): 🧪 resolve failing test cases and improve local testing automation (#7) (2022-07-27)

### Others

- chore(deps): update trunk-io/trunk-action digest to 22e948f (#59) (2022-10-07)
- chore(deps): update actions/stale digest to 5ebf00e (#58) (2022-10-07)
- chore: add trunk config checks as linting tool (#40) (2022-10-07)
- chore(deps): update amannn/action-semantic-pull-request digest to 505e44b (#46) (2022-09-28)
- docs: add forced-request as a contributor for maintenance (#53) (2022-09-28)
- docs: add EndlessTrax as a contributor for maintenance (#52) (2022-09-28)
- docs: add delineaKrehl as a contributor for maintenance (#51) (2022-09-27)
- docs: add tylerezimmerman as a contributor for maintenance (#50) (2022-09-27)
- docs: add sheldonhull as a contributor for code, doc, and test (#48) (2022-09-27)
- docs: add hansboder as a contributor for bug (#54) (2022-09-27)
- docs: add amigus as a contributor for code, doc, and test (#49) (2022-09-27)
- chore(deps): update actions/stale action to v6 (#43) (2022-09-27)
- chore(deps): update actions/stale digest to 9c1b1c6 (#38) (2022-09-01)
- chore(vscode): add missing linting settings (2022-09-01)
- chore: align github workflows and tooling (#37) (2022-09-01)
- ci(action/docker): 🔨 eliminate github event name push from skipping workflow dispatch (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.3 (#35) (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.2 (#34) (2022-08-29)
- chore(deps): update kubernetes packages to v0.25.0 (#25) (2022-08-26)
- ci(mend): disable issues for mend but still allow pull and renovate execution (2022-08-26)
- chore(deps): update github.com/mattbaird/jsonpatch digest to 098863c (#13) (2022-08-23)
- chore: update maintainers (#24) (2022-08-23)
- chore(devcontainer): 🧰 improvements to experience with setup directions, prebuilt tooling, and reliability (#23) (2022-08-20)
- chore(deps): update kubernetes packages to v0.24.4 (#22) (2022-08-18)
- ci(tests): fix name of tasks (2022-08-05)
- chore(deps): update ⬆️ golang module go to 1.19 (#17) (2022-08-05)
- ci(tests): remove push filter condition (2022-08-05)
- ci(tests): use mage to run tests (2022-08-05)
- ci(tests): 🔨 adjust test to use gotestsum (2022-08-05)
- ci(tests): ➕ require tests to be run as part of pull request (2022-08-05)
- ci: 🤖 add workflow dispatch to docker building actions (#21) (2022-08-03)
- ci(docker): set registry as org secret (#20) (2022-08-01)
- chore: use new env vars for docker hub publishing (#19) (2022-08-01)
- chore(deps): update ⬆️ golang module github.com/pterm/pterm to v0.12.45 (#16) [skip ci] (2022-07-29)
- chore(deps): update ⬆️ golang module github.com/mittwald/go-helm-client to v0.11.3 (#14) [skip ci] (2022-07-28)
- Configure Renovate AB#449139 (#11) (2022-07-28)
- chore(codeowner): assign default codeowners (#12) (2022-07-27)
- chore(helm): ⬆️ bump chart patch version due to minor changes (#8) (2022-07-27)
- Merge pull request #6 from The-Migus-Group/fix-meta (2022-07-12)
- Merge pull request #5 from DelineaXPM/fix-docker (2022-06-28)
- Merge pull request #4 from DelineaXPM/fix-3 (2022-06-28)
- Merge pull request #1 from DelineaXPM/delineaKrehl-DeepRebrand (2022-06-02)
- Fix for #13 that improves injector error handling. (#14) (2022-05-20)
- Add Testing 📛 (2022-05-06)
- Run go test ./... (2022-05-06)
- Shorten the name of the github action (2022-05-06)
- Fix build badges. 📛 (2022-05-06)
- Secret Synchronization Mechanism #11 (#12) (2022-05-06)
- Merge pull request #10 from thycotic/chart-upgrade-fixes (2022-03-22)
- Eval role logic after patchMode test; Fixes #8 (#9) (2022-03-10)
- Comments + all = install-image (2022-03-01)
- Publish to Minikube by default. 🪄 (2022-03-01)
- Reenable image publishing. (2022-03-01)
- Parity updates with TSS-K8s (#7) (2022-03-01)
- Add webhookScope (2022-02-17)
- Cleanup ♻️ (2022-02-17)
- Avoid hard-coding. (2022-02-16)
- Editing. ✂️ (2022-02-16)
- Fix for non-default namespaces. (2022-02-16)
- Bugfix: hardcoded namespace. (2022-02-16)
- Clean up. (2022-02-16)
- Further simplify. (2022-02-16)
- Simplify. (2022-02-16)
- Use a self-signed cert. (2022-02-16)
- Update README.md (2022-02-14)
- Wordsmithing. (2022-02-14)
- Initial revision of NOTES.txt for completeness. (2022-02-14)
- Make the Helm Chart use the image we built. 🔨 (2022-02-14)
- Typo and tweaks. ✨ (2022-02-14)
- Grammar. (2022-02-14)
- Minor improvements; fixes #3 (2022-02-14)
- Bugfix: whitespace error. (2022-02-14)
- Fix 'clean' (2022-02-13)
- Cleanup. (2022-02-13)
- Edits for Helm and OpenShift. (2022-02-13)
- v1beta1 compatibility (for OpenShift) (2022-02-12)
- Increase flexiblity and (hopefully) readablity. (2022-02-11)
- Update for Helm (#5) (2022-02-11)
- Avoid infinite loop (#6) (2022-02-11)
- Update to apiv1 (#4) (2022-02-11)
- Use Helm. (#5) (2022-02-11)
- Fix GPR push authentication failure. (2022-02-09)
- Upgrade to dsv-sdk-go v1.0.1. (2020-05-13)

## v1.0.0 (2022-10-10)

### Added

- feat: Configure Mend for GitHub.com (#26) (2022-08-26)

### Fixed

- fix(mend): 🐛 remove trailing space in json (2022-08-26)
- fix(tests): 🧪 resolve failing test cases and improve local testing automation (#7) (2022-07-27)

### Others

- chore(deps): update trunk-io/trunk-action digest to 22e948f (#59) (2022-10-07)
- chore(deps): update actions/stale digest to 5ebf00e (#58) (2022-10-07)
- chore: add trunk config checks as linting tool (#40) (2022-10-07)
- chore(deps): update amannn/action-semantic-pull-request digest to 505e44b (#46) (2022-09-28)
- docs: add forced-request as a contributor for maintenance (#53) (2022-09-28)
- docs: add EndlessTrax as a contributor for maintenance (#52) (2022-09-28)
- docs: add delineaKrehl as a contributor for maintenance (#51) (2022-09-27)
- docs: add tylerezimmerman as a contributor for maintenance (#50) (2022-09-27)
- docs: add sheldonhull as a contributor for code, doc, and test (#48) (2022-09-27)
- docs: add hansboder as a contributor for bug (#54) (2022-09-27)
- docs: add amigus as a contributor for code, doc, and test (#49) (2022-09-27)
- chore(deps): update actions/stale action to v6 (#43) (2022-09-27)
- chore(deps): update actions/stale digest to 9c1b1c6 (#38) (2022-09-01)
- chore(vscode): add missing linting settings (2022-09-01)
- chore: align github workflows and tooling (#37) (2022-09-01)
- ci(action/docker): 🔨 eliminate github event name push from skipping workflow dispatch (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.3 (#35) (2022-08-31)
- chore(deps): update kentaro-m/auto-assign-action action to v1.2.2 (#34) (2022-08-29)
- chore(deps): update kubernetes packages to v0.25.0 (#25) (2022-08-26)
- ci(mend): disable issues for mend but still allow pull and renovate execution (2022-08-26)
- chore(deps): update github.com/mattbaird/jsonpatch digest to 098863c (#13) (2022-08-23)
- chore: update maintainers (#24) (2022-08-23)
- chore(devcontainer): 🧰 improvements to experience with setup directions, prebuilt tooling, and reliability (#23) (2022-08-20)
- chore(deps): update kubernetes packages to v0.24.4 (#22) (2022-08-18)
- ci(tests): fix name of tasks (2022-08-05)
- chore(deps): update ⬆️ golang module go to 1.19 (#17) (2022-08-05)
- ci(tests): remove push filter condition (2022-08-05)
- ci(tests): use mage to run tests (2022-08-05)
- ci(tests): 🔨 adjust test to use gotestsum (2022-08-05)
- ci(tests): ➕ require tests to be run as part of pull request (2022-08-05)
- ci: 🤖 add workflow dispatch to docker building actions (#21) (2022-08-03)
- ci(docker): set registry as org secret (#20) (2022-08-01)
- chore: use new env vars for docker hub publishing (#19) (2022-08-01)
- chore(deps): update ⬆️ golang module github.com/pterm/pterm to v0.12.45 (#16) [skip ci] (2022-07-29)
- chore(deps): update ⬆️ golang module github.com/mittwald/go-helm-client to v0.11.3 (#14) [skip ci] (2022-07-28)
- Configure Renovate AB#449139 (#11) (2022-07-28)
- chore(codeowner): assign default codeowners (#12) (2022-07-27)
- chore(helm): ⬆️ bump chart patch version due to minor changes (#8) (2022-07-27)
- Merge pull request #6 from The-Migus-Group/fix-meta (2022-07-12)
- Merge pull request #5 from DelineaXPM/fix-docker (2022-06-28)
- Merge pull request #4 from DelineaXPM/fix-3 (2022-06-28)
- Merge pull request #1 from DelineaXPM/delineaKrehl-DeepRebrand (2022-06-02)
- Fix for #13 that improves injector error handling. (#14) (2022-05-20)
- Add Testing 📛 (2022-05-06)
- Run go test ./... (2022-05-06)
- Shorten the name of the github action (2022-05-06)
- Fix build badges. 📛 (2022-05-06)
- Secret Synchronization Mechanism #11 (#12) (2022-05-06)
- Merge pull request #10 from thycotic/chart-upgrade-fixes (2022-03-22)
- Eval role logic after patchMode test; Fixes #8 (#9) (2022-03-10)
- Comments + all = install-image (2022-03-01)
- Publish to Minikube by default. 🪄 (2022-03-01)
- Reenable image publishing. (2022-03-01)
- Parity updates with TSS-K8s (#7) (2022-03-01)
- Add webhookScope (2022-02-17)
- Cleanup ♻️ (2022-02-17)
- Avoid hard-coding. (2022-02-16)
- Editing. ✂️ (2022-02-16)
- Fix for non-default namespaces. (2022-02-16)
- Bugfix: hardcoded namespace. (2022-02-16)
- Clean up. (2022-02-16)
- Further simplify. (2022-02-16)
- Simplify. (2022-02-16)
- Use a self-signed cert. (2022-02-16)
- Update README.md (2022-02-14)
- Wordsmithing. (2022-02-14)
- Initial revision of NOTES.txt for completeness. (2022-02-14)
- Make the Helm Chart use the image we built. 🔨 (2022-02-14)
- Typo and tweaks. ✨ (2022-02-14)
- Grammar. (2022-02-14)
- Minor improvements; fixes #3 (2022-02-14)
- Bugfix: whitespace error. (2022-02-14)
- Fix 'clean' (2022-02-13)
- Cleanup. (2022-02-13)
- Edits for Helm and OpenShift. (2022-02-13)
- v1beta1 compatibility (for OpenShift) (2022-02-12)
- Increase flexiblity and (hopefully) readablity. (2022-02-11)
- Update for Helm (#5) (2022-02-11)
- Avoid infinite loop (#6) (2022-02-11)
- Update to apiv1 (#4) (2022-02-11)
- Use Helm. (#5) (2022-02-11)
- Fix GPR push authentication failure. (2022-02-09)
- Upgrade to dsv-sdk-go v1.0.1. (2020-05-13)
