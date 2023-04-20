# Troubleshooting

## Obtaining Logs

For both customers and development, stern allows easier debugging by providing a stream of the logs for both syncer & injector in one workflow.

Use Stern to easily stream cross namespace logs with the `dsv-filter-selector` by running:

Aqua installs this automatically, but if you want to do this manually grab from github releases like this

If not using the `aqua`, you can modify the following command to match the version and OS you are using and download directly.

```shell
$(curl -fSSl https://github.com/wercker/stern/releases/download/1.11.0/stern_linux_amd64 -o ./stern) && sudo chmod +x ./stern && sudo mv ./stern /usr/local/bin
```

Alternative install options: [Stern Installation](https://github.com/stern/stern#installation)

### Alternative - If Using Mage For Local Development

While `mage k8s:logs` will run this for you, manually you can invoke like this:

```shell
# For all pods in the namespace run
stern --kubeconfig .cache/config --namespace dsv --timestamps .

# For pods with the selector run
stern --kubeconfig .cache/config --namespace dsv --timestamps --selector 'dsv-filter-name in (dsv-syncer, dsv-injector)'
```

### Example for Providing Logs for Support

If debugging, you can stream logs from Kubernetes with this tool and capture to a log file for providing in support cases.

```shell
stern --kubeconfig .cache/config --namespace dsv --timestamps . > activity.log
```
