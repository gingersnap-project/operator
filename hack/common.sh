#!/usr/bin/env bash
set -e

_script="$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")"

export PATH="$_script/../bin":$PATH
export TESTING_NAMESPACE=${TESTING_NAMESPACE-namespace-for-testing}
export IMG_REGISTRY=${IMG_REGISTRY-"localhost:5000"}
export IMG_VERSION=${IMG_VERSION-"latest"}

export VERSION=${VERSION-"0.0.1"}
export DEFAULT_CHANNEL=alpha
export CATALOG_SOURCE_NAME=${CATALOG_SOURCE_NAME-test-catalog}
export CATALOG_SOURCE_NAMESPACE=${CATALOG_SOURCE_NAMESPACE-olm}

export IMAGE_TAG_BASE=${IMG_REGISTRY}/operator
export BUNDLE_IMG_TAG=${BUNDLE_IMG_TAG-"v${VERSION}"}
export BUNDLE_IMG=${IMAGE_TAG_BASE}-bundle:${BUNDLE_IMG_TAG}
export CATALOG_IMG=${IMAGE_TAG_BASE}-catalog
export IMG=${IMAGE_TAG_BASE}:${IMG_VERSION}
