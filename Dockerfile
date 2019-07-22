FROM golang:1.12.7 AS build
WORKDIR /go/src/github.com/storageos/nfs/
COPY . /go/src/github.com/storageos/nfs/
RUN make build

FROM storageos/nfs-base:20190830-0004

LABEL maintainer="support@storageos.com"

COPY --from=build /go/src/github.com/storageos/nfs/build/_output/bin/nfs /nfs

# TESTING ONLY
COPY --from=build /go/src/github.com/storageos/nfs/export.conf /export.conf

# NFS port and http daemon
EXPOSE 2049 80

# Start Ganesha NFS daemon by default
CMD ["/nfs"]
