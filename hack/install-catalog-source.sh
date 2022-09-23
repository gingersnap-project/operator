#!/usr/bin/env bash
set -e

TESTING_NAMESPACE=${TESTING_NAMESPACE-namespace-for-testing}
IMG_REGISTRY=${IMG_REGISTRY-"localhost:5000"}

export VERSION=0.0.1
export DEFAULT_CHANNEL=alpha

export IMAGE_TAG_BASE=${IMG_REGISTRY}/gingersnap-operator
export BUNDLE_IMG=${IMAGE_TAG_BASE}-bundle:v${VERSION}
export CATALOG_IMG=${IMAGE_TAG_BASE}-catalog
export IMG=${IMAGE_TAG_BASE}

# Create the operator image
#make docker-build docker-push

# Create the operator bundle image
make bundle bundle-build bundle-push

# Create the OLM catalog image
make catalog-build catalog-push

# Create the namespace and CatalogSource
kubectl create namespace "${TESTING_NAMESPACE}" || true
kubectl delete CatalogSource test-catalog -n "${TESTING_NAMESPACE}" || true
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: test-catalog
  namespace: ${TESTING_NAMESPACE}
spec:
  displayName: Test Operators Catalog
  image: ${CATALOG_IMG}
  sourceType: grpc
EOF
