# NFS Container for StorageOS

[![Build Status](https://travis-ci.org/storageos/nfs.svg?branch=master)](https://travis-ci.org/storageos/nfs)

NFS container used to provide RWX StorageOS volumes.  Intended to be used with
the [StorageOS Cluster Operator](https://github.com/storageos/cluster-operator).

The NFS container makes use of
[nfs-ganesha](http://github.com/nfs-ganesha/nfs-ganesha/), ensuring that all
system dependencies are met and running if required.

Configuration is optimised for sharing StorageOS block volumes as Kubernetes RWX
volumes, with a single NFS server per underlying block volume.  Only NFSv4 is
supported.

For more information, see https://docs.storageos.com.

## Build

```shell
make image IMAGE=storageos/nfs:test
```

Once the image has been built, the `nfs` binary can be recompiled and the image
updated quickly by running:

```shell
make update IMAGE=storageos/nfs:test
```

## Run it on host

Build the init container with `make image` or `make update` and run it on the
host with `make run`.

## Configuration

The NFS container is configured using environment variables.  Only
`GANESHA_CONFIGFILE` is required and must be set to the path of a valid
[nfs-ganesha](http://github.com/nfs-ganesha/nfs-ganesha/) configuration
file.

### Environment variables

| Variable Name             | Valid in versions | Description |
| :------------------       | :---------------- | :---------- |
| `GANESHA_CONFIGFILE`      | 1.0+              | (REQUIRED) Path to a valid [nfs-ganesha](http://github.com/nfs-ganesha/nfs-ganesha/) configuration file |
| `LISTEN_ADDR`             | 1.0+              | HTTP server listen address. Default `:80` |
| `DISABLE_METRICS`         | 1.0+              | Disables the /metrics endpoint if set to `true`. Default `false` |
| `NAME`                    | 1.0+              | Name of the NFS server.  Corresponds to the RWX volume name.  Used to label Prometheus metrics. |
| `NAMESPACE`               | 1.0+              | Namespace of the NFS server. Used to label Prometheus metrics. |

## Health

NFS server health is reported by querying `/healthz` on the HTTP server
`LISTEN_ADDR`.

`HTTP 200/OK` will be returned when the server is operational and sending
heartbeat messages.

`HTTP 503/Service Unavailable` will be returned if the server hasn't sent a
heartbeat message within 10 seconds.

## Prometheus metrics

Prometheus metrics are available by querying `/metrics` on the HTTP server
`LISTEN_ADDR`.

Metrics are returned for:

- Process statistics, including memory usage, threads, cpu time and file
  descriptors.
- NFS server exports, reported per export for each protocol in use.  Read and
  write operations are broken down into:

    - Requested bytes total
    - Transfered bytes total
    - Operations total
    - Cumulative operations latency in seconds
    - Cumulative wait queue in seconds

- NFS clients, reported per client connection for each protocol in use.  Only
  NFSv40 and NFSv41 are supported.  Ganesha has not yet implemented client
  metrics for NFSv42. Read and write operations are broken down into:

    - Requested bytes total
    - Transfered bytes total
    - Operations total
    - Cumulative operations latency in seconds
    - Cumulative wait queue in seconds

If `NAME` and/or `NAMESPACE` environment values are set, metrics are labeled
with `name=NAME` and `namespace=NAMESPACE`.
