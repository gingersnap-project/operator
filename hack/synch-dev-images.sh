#!/bin/bash

IMG_TAG_BASE=quay.io/gingersnap

docker pull ${IMG_TAG_BASE}/db-syncer
docker pull ${IMG_TAG_BASE}/cache-manager

kind load docker-image ${IMG_TAG_BASE}/db-syncer
kind load docker-image ${IMG_TAG_BASE}/cache-manager

