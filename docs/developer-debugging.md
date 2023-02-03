# Developer Debugging

This documentation comes from the original `make` driven build process.
It hasn't been migrated to `mage` and will be noted here unless deprecated.

## Host (for debugging)

Per above, typically, the injector runs as a POD in the cluster but running it on the host makes debugging easier.

```sh
make install-host EXTERNAL_NAME=laptop.mywifi.net CA_BUNDLE=$(cat /path/to/ca.crt | base64 -w0 -)
```

For it to work:

- The certificate that the injector presents must validate against the `$(CA_BUNDLE)`.
- The certificate must also have a Subject Alternative Name for `$(INJECTOR_NAME).$(NAMESPACE).svc`.
  By default that's `dsv-injector.dsv.svc`.

- The `$(EXTERNAL_NAME)` is a required argument, and the name itself must be resolvable _inside_ the cluster.
  **localhost will not work**.

If the `$(CA_BUNDLE)` is argument is omitted, `make` will attempt to extract it from `kubectl config`:

```make
install-host: CA_BUNDLE_KUBE_CONFIG_INDEX = 0
install-host: CA_BUNDLE_JSON_PATH = {.clusters[$(CA_BUNDLE_KUBE_CONFIG_INDEX)].cluster.certificate-authority-data}
install-host: CA_BUNDLE=$(shell $(KUBECTL) config view --raw -o jsonpath='$(CA_BUNDLE_JSON_PATH)' | tr -d '"')

```

which will make:

```sh
kubectl config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}' | tr -d '"'
```

Optionally set `$(CA_BUNDLE_KUBE_CONFIG_INDEX)` to use `1`, to use the second cluster in your configuration,
`2` for the third and so on.
ℹ️ All this assumes that the injector uses a certificate signed by the cluster CA.
There are several options like [cert-manager](https://cert-manager.io/)
for getting cluster-signed certs, however,
this simple [bash script](https://gist.github.com/amigus/b4e6e642f88e756be1996e44a1c35349)
will request and grant a suitable certificate from the cluster using cURL and OpenSSL.
To use it:

```sh
get_k8s_cert.sh -n dsv-injector -N dsv
```

Now run it:

```sh
./dsv-injector -cert ./dsv-injector.pem -key ./dsv-injector.key -credentials ./configs/credentials.json -address :8543
```
