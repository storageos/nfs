#!/usr/bin/env bash

set -Eeuxo pipefail

cluster=kind-1


run_kind() {
    if [ ! -x /usr/local/bin/kind ]; then
        echo "Download kind binary..."
        wget -O kind 'https://github.com/kubernetes-sigs/kind/releases/download/v0.5.1/kind-linux-amd64' --no-check-certificate && chmod +x kind && sudo mv kind /usr/local/bin/
    fi

    if [ ! -x /usr/local/bin/kubectl ]; then
        echo "Download kubectl..."
        curl -Lo kubectl https://storage.googleapis.com/kubernetes-release/release/"${K8S_VERSION}"/bin/linux/amd64/kubectl && chmod +x kubectl && sudo mv kubectl /usr/local/bin/
    fi

    if ! kind get clusters | grep ${cluster}; then
        echo "Create Kubernetes cluster with kind..."
        # kind create cluster --image=kindest/node:"$K8S_VERSION"
        kind create cluster --image storageos/kind-node:"$K8S_VERSION" --name ${cluster}
    fi

    echo "Export kubeconfig..."
    export KUBECONFIG="$(kind get kubeconfig-path --name="${cluster}")"
    echo "export KUBECONFIG=${KUBECONFIG}"

    echo "Get cluster info..."
    kubectl cluster-info
    echo

    echo "Wait for kubernetes to be ready"
    JSONPATH='{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{end}'; until kubectl get nodes -o jsonpath="$JSONPATH" 2>&1 | grep -q "Ready=True"; do sleep 1; done
    echo
}

run_nfs() {
    # Copy the build container image into KinD.
    x=$(docker ps -f name=${cluster}-control-plane -q)
    docker save storageos/nfs:test > nfs.tar
    docker cp nfs.tar $x:/nfs.tar

    # containerd load image from tar archive (KinD with containerd).
    docker exec $x bash -c "ctr -n k8s.io images import --base-name docker.io/storageos/nfs:test /nfs.tar"

    echo "Creating tmpfs for /export"
    # See:
    # https://github.com/kubernetes-sigs/kind/issues/118#issuecomment-524239315
    docker exec $x bash -c "if [ ! -d /export ]; then mkdir /export; fi"
    docker exec $x bash -c "if df /export | grep overlay; then mount -t tmpfs -o rw,size=100M tmpfs /export; fi"

    # Start the nfs container
    kubectl apply -f test/e2e/nfs.yaml

    echo "Waiting for the nfs pod to start"
    until kubectl get pod nfs-pod --no-headers -o go-template='{{.status.phase}}' | grep -q Running; do sleep 5; done
    echo "nfs-pod is running"
}

# Give nfsvers (40, 41, 42 as param)
test_client() {
    nfsver=$1
    echo "Running NFS v${nfsver} test"
    kubectl apply -f test/e2e/pv-${nfsver}.yaml
    kubectl apply -f test/e2e/pvc.yaml
    kubectl apply -f test/e2e/client.yaml

    # arbitrary wait time, may need longer
    sleep 5
    kubectl describe pod client-pod
    kubectl describe pvc
    kubectl describe pv
    kubectl exec -it client-pod dd if=/dev/urandom of=/mnt/testdata${nfsver} bs=4k count=1024

    # will fail if not written to NFS share
    docker exec -it ${cluster}-control-plane rm /export/testdata${nfsver}

    # check metrics (only available for 4.0 and 4.1, not 4.2)
    case $nfsver in
    4[01])
        if [ $(docker exec -it ${cluster}-control-plane curl http://127.0.0.1/metrics | grep -c ^storageos_clients_nfs_v${nfsver}) != 12 ]; then
            echo "NFS v${nfsver} metrics not found"
            exit 1
        fi
        ;;
    42)
        if [ $(docker exec -it ${cluster}-control-plane curl http://127.0.0.1/metrics | grep -c ^storageos_clients_nfs_v${nfsver}) != 0 ]; then
            echo "NFS v${nfsver} found but not expected"
            exit 1
        fi
        ;;
    esac


    # cleanup
    kubectl delete -f test/e2e/client.yaml --force --grace-period=0
    kubectl delete -f test/e2e/pvc.yaml
    kubectl delete -f test/e2e/pv-${nfsver}.yaml
}

test_health() {
    if [ $(docker exec -it ${cluster}-control-plane curl http://127.0.0.1/healthz) != "ok" ]; then
        echo "Health endpoint did not return ok"
        exit 1
    fi
    if [ $(docker exec -it ${cluster}-control-plane curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1/healthz) != "200" ]; then
        echo "Health endpoint did not return status 200"
        exit 1
    fi
}

main() {
    run_kind

    echo "Ready for testing"

    make unittest
    make image
    run_nfs

    for ver in 40 41 42; do
        test_client $ver
    done

    test_health

    echo "nfs container logs"
    kubectl logs nfs-pod -c storageos-nfs

    # Stop the nfs container
    echo "Removing the nfs pod"
    kubectl delete -f test/e2e/nfs.yaml
}

main
