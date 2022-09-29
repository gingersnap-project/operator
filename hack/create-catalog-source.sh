#!/usr/bin/env bash
set -e

DIRNAME=$(dirname "$0")
. "$DIRNAME/common.sh"

# Create the operator image
make docker-build docker-push

# Create the operator bundle image
make bundle bundle-build bundle-push

# Create the OLM catalog image
make catalog-build catalog-push
