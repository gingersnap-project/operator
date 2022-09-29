#!/usr/bin/env bash
set -e

DIRNAME=$(dirname "$0")
. "$DIRNAME/create-catalog-source.sh"

# Create the namespace and CatalogSource
kubectl create namespace "${CATALOG_SOURCE_NAMESPACE}" || true
kubectl delete CatalogSource "${CATALOG_SOURCE_NAME}" -n "${CATALOG_SOURCE_NAMESPACE}" || true
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: ${CATALOG_SOURCE_NAME}
  namespace: ${CATALOG_SOURCE_NAMESPACE}
spec:
  displayName: Test Operators Catalog
  image: ${CATALOG_IMG}
  sourceType: grpc
EOF
