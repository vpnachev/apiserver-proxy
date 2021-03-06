# API Server proxy

This component consists of `apiserver-proxy-sidecar` which runs on every `Node` in a `Shoot` cluster.
It does the following:

1. creates a dummy interface `apisrv0` (configurable with the optional `--interface` flag)
1. gets the IP address of the `default/kubernetes` service (`KUBERNETES_SERVICE_HOST` environment variable or directly from the `--ip-address` flag).
1. execute the following `iptables` rules (in this case the IP address of the master service is `10.96.0.2`):

    ```text
    -A PREROUTING -t raw -d 10.96.0.2/32 -p tcp -m tcp --dport 443 -j NOTRACK
    -A OUTPUT -t raw -d 10.96.0.2/32 -p tcp -m tcp --dport 443 -j NOTRACK
    -A OUTPUT -t raw -s 10.96.0.2/32 -p tcp -m tcp --sport 443 -j NOTRACK
    -A INPUT -t filter -d 10.96.0.2/32 -p tcp -m tcp --dport 443 -j ACCEPT
    -A OUTPUT -t filter -s 10.96.0.2/32 -p tcp -m tcp --sport 443 -j ACCEPT
    ```

    Those rules allow traffic to this IP address and disable conntrack as the IP address is local to the machine.

1. Every 1 min repeats the process and start from `1.`

After this, the actual `apiserver-proxy` can listen on this IP address (`10.96.0.2`) and send traffic to the correct kube-apiserver.
The implementation of that proxy is fully transparent and can be replaced at any given moment without any modifications to the `apiserver-proxy-sidecar`.

## Command line options

```shell
bazel run //cmd/apiserver-proxy-sidecar -- --help
Usage of apiserver-proxy-sidecar:
      --add_dir_header                   If true, adds the file directory to the header
      --alsologtostderr                  log to standard error as well as files
      --cleanup                          indicates whether created interface and iptables should be removed on exit
      --interface string                 [optional] name of the interface to be created (default "apisrv0")
      --ip-address string                [optional] ip-address on which the proxy is listening.
                                                        If not set, it uses the "KUBERNETES_SERVICE_HOST" environment variable
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --log_file string                  If non-empty, use this log file
      --log_file_max_size uint           Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                      log to standard error instead of files (default true)
      --setupiptables                    indicates whether iptables rules should be setup (default true)
      --skip_headers                     If true, avoid header prefixes in the log messages
      --skip_log_headers                 If true, avoid headers when opening log files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
      --syncinterval duration            [optional] interval to check for iptables rules (default 1m0s)
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
pflag: help requested
```

## Development

[bazel](https://bazel.build/) is used for building and testing. Optional `Dockerfile` is provided.

### Update dependencies

```shell
go mod tidy

bazel run //:gazelle -- update-repos -from_file=go.mod
bazel run //:gazelle
```

### Running on Docker locally

```shell
bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/apiserver-proxy-sidecar:go_image
```

### Just pushing to local Docker

```shell
bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/apiserver-proxy-sidecar:go_image -- --norun
```

### Testing

> Note: It requires a running Docker.

```shell
bazel test "//..."
```
