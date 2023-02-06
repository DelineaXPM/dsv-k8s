# Local Testing

## First Time Setup

- A valid DSV tenant.
- Creation of a dsv secret and the configured client credentials to setup this.
  See [configuring-dsv](configure.md#configuring-dsv)
- The test helm chart values file to be updated with the credentials.
  It's located at: `.cache/dsv-injector/values.yaml`.

<img src="assets/warning-app1-required-for-tests.svg">

## PENDING Contributor Improvements

For dsv-team members, the goal is to load all this directly from a team vault. Right now this project has not been migrated to this, so you have to setup manually the first time.

## Test Environment Configuration

| Environment Variable                   | Default                          | Explanation                                                 |
| -------------------------------------- | -------------------------------- | ----------------------------------------------------------- |
| `DSV_K8S_TEST_CONFIG`                  | _none_                           | Contain a JSON string containing a valid `credentials.json` |
| `DSV_K8S_TEST_SECRET_PATH`             | `/test/secret`                   | The path to the secret to test against in the vault         |
| DEPRECATED: `DSV_K8S_TEST_CONFIG_FILE` | `../../configs/credentials.json` | The path to a valid `credentials.json`                      |

ℹ️ NOTE: `DSV_K8S_TEST_CONFIG` takes precedence over `DSV_K8S_TEST_CONFIG_FILE`

For example:

```shell
DSV_K8S_TEST_CONFIG='{"app1":{"credentials":{"clientId":"","clientSecret":""},"tenant":"mytenant"}}' \
DSV_K8S_TEST_SECRET_PATH=my:test:secret \
mage go:testsum ./...
```
