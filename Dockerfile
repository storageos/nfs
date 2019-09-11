FROM golang:1.12.7 AS build
WORKDIR /go/src/github.com/storageos/nfs/
COPY . /go/src/github.com/storageos/nfs/
RUN make build

FROM storageos/nfs-base:20190903-0018

LABEL maintainer="support@storageos.com"

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
