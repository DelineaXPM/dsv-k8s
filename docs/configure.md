# Configure

This focuses on the DSV configuration required to use with Kubernetes.
This applies to both local testing Kubernetes and your own seperate cluster.

## JSON Credentials for Helm Install

The configuration requires a JSON formatted list of Client Credential and Tenant mappings.

The name of the credential (such as `app1` or `default`) is used for matching the annontated credential to the right credentials file to use to connect to the connect tenant.

You can place your temporary config in `.cache/credentials.json` as this is ignored by git, so that you can run the helm install command manually if you aren't doing local development.

<img src="assets/info-markup-default-creds.svg">

```json
{
  "app1": {
    "credentials": {
      "clientId": "93d866d4-635f-4d4e-9ce3-0ef7f879f319",
      "clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxxxx"
    },
    "tenant": "mytenant"
  },
  "default": {
    "credentials": {
      "clientId": "64241412-3934-4aed-af26-95b1eaba0e6a",
      "clientSecret": "xxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxx-xxxxx"
    },
    "tenant": "mytenant"
  }
}
```

### Update Manifests

This would be referenced by a Kubernetes secret with annontations like:

```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: user-domain-pass
  annotations:
    dsv.delinea.com/credentials: app1
    dsv.delinea.com/set-secret: 'k8s:sync:test'
```

If using the provided examples, you can edit: `.cache/manifests` and adjust the secrets to map.
You can use all of the provided manifests to test the different behavior, or just deploy one if desired.

## Configuring Credentials in Kubernetes To Talk to DSV

## Configuring DSV

The following is an example of the steps to setup for testing, but can be modified to support your use case.

Create the role that will allow creating a client for programmatic access

```shell
dsv role create --name 'k8s' --desc 'test profile for k8s'
dsv secret create --path 'k8s:sync:test' --data '{"password": "admin","username": "admin"}'
```

Create a policy that allows the local user to read the secret, modify this to the correct user/group mapping:

```shell
dsv policy create -- actions 'read' --path 'secrets:k8s' --desc 'test access to secret' --resources 'secrets:k8s:<.*>' --subjects 'roles:k8s'
dsv client create --role k8s
```
