#!/usr/bin/env bash
set -e

DIRNAME=$(dirname "$0")
. "$DIRNAME/common.sh"

cat <<EOF | kubectl apply -f -
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: gingersnap
  namespace: ${SUBSCRIPTION_NAMESPACE}
spec:
  name: gingersnap
  channel: alpha
  source: test-catalog
  sourceNamespace: ${CATALOG_SOURCE_NAMESPACE}
  installPlanApproval: automatic
EOF
