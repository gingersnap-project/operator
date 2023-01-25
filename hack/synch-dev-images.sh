#!/bin/bash

IMAGE_NAMES="db-syncer cache-manager-mssql cache-manager-mysql cache-manager-postgres"
IMAGE_TAG_BASE=quay.io/gingersnap

for IMG_NAME in ${IMAGE_NAMES}; do
  IMAGE="${IMAGE_TAG_BASE}/${IMG_NAME}"
  docker pull "${IMAGE}"
  kind load docker-image "${IMAGE}"
done
