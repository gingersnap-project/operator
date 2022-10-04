#!/usr/bin/env bash
set -e

_script="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"

export PATH="$_script/../bin":$PATH

export CATALOG_SOURCE_NAME=${CATALOG_SOURCE_NAME-test-catalog}
export CATALOG_SOURCE_NAMESPACE=${CATALOG_SOURCE_NAMESPACE-olm}
export SUBSCRIPTION_NAMESPACE=${SUBSCRIPTION_NAMESPACE-operators}
