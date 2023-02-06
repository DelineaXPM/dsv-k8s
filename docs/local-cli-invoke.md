# Local CLI Invoke

The cli can be invoked locally.

The injector uses the HTTPS server built-in to the Golang [http](https://pkg.go.dev/net/http)
package to host the Kubernetes Mutating Webhook Webservice.

```bash
$ ./dsv-injector -h
Usage of ./dsv-injector:
  -address string
        the address to listen on, e.g., 'localhost:8080' or ':8443' (default ":18543")
  -cert string
        the path of the public certificate file in PEM format (default "tls/cert.pem")
  -credentials string
        the path of JSON formatted credentials file (default "credentials/config.json")
  -key string
        the path of the private key file in PEM format (default "tls/key.pem")
```

Thus the injector can run "anywhere," but, typically,
the injector runs as a POD in the Kubernetes cluster that uses it.
The syncer is a simple Golang executable.
It typically runs as a Kubernetes CronJob, but it will run outside the cluster.

```bash
$ ./dsv-syncer -h
Usage of ./dsv-syncer:
  -credentials string
        the path of JSON formatted credentials file (default "credentials/config.json")
  -kubeConfig string
        the Kubernetes Client API configuration file; ignored when running in-cluster (default "/home/user/.kube/config")
  -namespace string
        the Kubernetes namespace containing the Secrets to sync; "" (the default) for all namespaces
```
