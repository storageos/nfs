IMAGE ?= storageos/nfs:test
GO_BUILD_CMD = go build -v
GO_ENV = GOOS=linux CGO_ENABLED=0

all: unittest build

.PHONY: build

build:
	@echo "Building nfs"
	$(GO_ENV) $(GO_BUILD_CMD) -o ./build/_output/bin/nfs .

image:
	docker build --no-cache . -f Dockerfile -t $(IMAGE)

update:
	docker build --no-cache . -f Dockerfile.update -t $(IMAGE)

unittest:
	go test -v -race `go list -v ./...`

clean:
	rm -rf build/_output

# Run the nfs server on the host.
run:
	docker run --rm \
		-p 2049 \
		-p 80 \
		-v /export:/export \
		--cap-add=SYS_ADMIN \
		--privileged \
		-e GANESHA_CONFIGFILE=/export.conf \
		-e NAME=test \
		-e NAMESPACE=default \
		-e DISABLE_METRICS=false \
		$(IMAGE)