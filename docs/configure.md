# Configure

This focuses on the DSV configuration required to use with Kubernetes.
This applies to both local testing Kubernetes and your own seperate cluster.

## Help Getting Started

Run `mage dsv:setupdsv` to create the required DSV configuration for testing.
This requires you to have already run `dsv init` in the project and runs against the profile you specified in `.env`.
You should ensure `direnv allow` has been run and the `.env` file is loaded.
Your `zsh` terminal should warn you if you didn't create the `.env` file.

The order:

- `mage dsv:setupdsv`
- `mage dsv:createsecret`
- `mage dsv:convertClientToCredentials`

To tear down and recreate with new secret, just run `mage dsv:destroy`

## Manually Creating (Prior Method Before Automation)

### JSON Credentials for Helm Install

The configuration requires a JSON formatted list of Client Credential and Tenant mappings.

The name of the credential (such as `app1` or `default`) is used for matching the annontated credential to the right credentials file to use to connect to the connect tenant.

By default, the `tld` is set to `.com` but you can change it by changing `tld` field to `eu` (for example: `"tld": "eu"`).

You can place your temporary config in `.cache/credentials.json` as this is ignored by git, so that you can run the helm install command manually if you aren't doing local development.

<img src="assets/info-markup-default-creds.svg">

```json
{
  "app1": {
    "credentials": {
      "clientId": "",
      "clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxxxx"
    },
    "tenant": "mytenant",
    "tld": "eu"
  },
  "default": {
    "credentials": {
      "clientId": "",
      "clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxxxx"
    },
    "tenant": "mytenant"
  }
}
```

### Update Manifests

This would be referenced by a Kubernetes secret with annotations like:

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: user-domain-pass
  annotations:
    dsv.delinea.com/credentials: app1
    dsv.delinea.com/set-secret: 'tests:dsv-k8s'
data:
  dummy_value: ""
```

If using the provided examples, you can edit: `.cache/manifests` and adjust the secrets to map.
You can use all of the provided manifests to test the different behavior, or just deploy one if desired.

## Configuring Credentials in Kubernetes To Talk to DSV

## Configuring DSV

The following is an example of the steps to setup for testing, but can be modified to support your use case.

Create the role that will allow creating a client for programmatic access

```shell
dsv role create --name 'k8s' --desc 'test profile for k8s'
dsv secret create --path 'tests:dsv-k8s' --data '{"password": "admin","username": "admin"}'
```

Create a policy that allows the local user to read the secret, modify this to the correct user/group mapping:

```shell
dsv policy create -- actions 'read' --path 'secrets:k8s' --desc 'test access to secret' --resources 'secrets:k8s:<.*>' --subjects 'roles:k8s'
dsv client create --role k8s
```
