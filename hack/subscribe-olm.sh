#!/usr/bin/env bash
set -e

TESTING_NAMESPACE=${TESTING_NAMESPACE-namespace-for-testing}

kubectl create namespace "${TESTING_NAMESPACE}" || true
cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: engytita
  namespace: ${TESTING_NAMESPACE}
---
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: engytita
  namespace: ${TESTING_NAMESPACE}
spec:
  channel: alpha
  source: test-catalog
  sourceNamespace: ${TESTING_NAMESPACE}
  name: engytita
EOF
