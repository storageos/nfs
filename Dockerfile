FROM golang:1.13.5 AS build

WORKDIR /go/src/github.com/storageos/nfs/
COPY . /go/src/github.com/storageos/nfs/
RUN make build

FROM storageos/nfs-base:20190903-0018
LABEL name="StorageOS Shared File" \
    maintainer="support@storageos.com" \
    vendor="StorageOS" \
    version="1.0.1" \
    release="1" \
    distribution-scope="public" \
    architecture="x86_64" \
    url="https://docs.storageos.com" \
    io.k8s.description="The StorageOS Shared File container is used by the Cluster Operator to provide ReadWriteMany (RWX) volumes." \
    io.k8s.display-name="StorageOS Shared File" \
    io.openshift.tags="" \
    summary="Highly-available persistent block and shared file storage for containerized applications." \
    description="StorageOS transforms commodity server or cloud based disk capacity into enterprise-class storage to run persistent workloads such as databases in containers. Provides high availability, low latency persistent block storage. No other hardware or software is required."

COPY --from=build /go/src/github.com/storageos/nfs/LICENSE /licenses/
COPY --from=build /go/src/github.com/storageos/nfs/build/_output/bin/nfs /nfs

# Use for testing only.  Exports /export
COPY --from=build /go/src/github.com/storageos/nfs/export.conf /export.conf

# Disable sssd and systemd name services to allow faster DBus startup.
COPY --from=build /go/src/github.com/storageos/nfs/nsswitch.conf /etc/nsswitch.conf

RUN mkdir /run/dbus && chmod 755 /run/dbus && chown dbus:dbus /run/dbus

# NFS port and http daemon
EXPOSE 2049 80

# Start Ganesha NFS daemon by default
CMD ["/nfs"]
