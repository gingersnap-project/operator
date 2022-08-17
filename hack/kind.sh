#!/usr/bin/env bash
# Modified version of the script found at https://kind.sigs.k8s.io/docs/user/local-registry/#create-a-cluster-and-registry
set -o errexit

CERT_MANAGER_VERSION="v1.8.0"
KINDEST_NODE_VERSION=${KINDEST_NODE_VERSION:-'v1.23.4'}
KIND_SUBNET=${KIND_SUBNET-172.172.0.0}
OLM_VERSION="v0.21.2"

docker network create kind --subnet "${KIND_SUBNET}/16" || true

# create registry container unless it already exists
reg_name='kind-registry'
reg_port=${KIND_PORT-'5000'}
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -d --restart=always -p "127.0.0.1:${reg_port}:5000" --name "${reg_name}" \
    quay.io/infinispan-test/registry:2
fi

# create a cluster with the local registry enabled in containerd
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
containerdConfigPatches:
- |-
  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."localhost:${reg_port}"]
    endpoint = ["http://${reg_name}:${reg_port}"]
nodes:
  - role: control-plane
    image: quay.io/infinispan-test/kindest-node:${KINDEST_NODE_VERSION}
EOF

# connect the registry to the cluster network
# (the network may already be connected)
docker network connect "kind" "${reg_name}" || true

# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: local-registry-hosting
  namespace: kube-public
data:
  localRegistryHosting.v1: |
    host: "localhost:${reg_port}"
    help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

# Attempt to load cert-manager images from local docker registry before installing
kind load docker-image quay.io/jetstack/cert-manager-cainjector:${CERT_MANAGER_VERSION} || true
kind load docker-image quay.io/jetstack/cert-manager-controller:${CERT_MANAGER_VERSION}  || true
kind load docker-image quay.io/jetstack/cert-manager-webhook:${CERT_MANAGER_VERSION} || true

# Install cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/${CERT_MANAGER_VERSION}/cert-manager.yaml

# Install OLM
operator-sdk olm install --version ${OLM_VERSION}

# Sometimes olm install does not wait long enough for deployments to be rolled out
kubectl wait --for=condition=available --timeout=60s deployment/catalog-operator -n olm
kubectl wait --for=condition=available --timeout=60s deployment/olm-operator -n olm
kubectl wait --for=condition=available --timeout=60s deployment/packageserver -n olm

# Wait for cert-manager
kubectl wait --for=condition=available --timeout=60s deployment/cert-manager -n cert-manager
kubectl wait --for=condition=available --timeout=60s deployment/cert-manager-cainjector -n cert-manager
kubectl wait --for=condition=available --timeout=60s deployment/cert-manager-webhook -n cert-manager
