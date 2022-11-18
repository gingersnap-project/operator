#!/bin/bash
set -e

PROJECT_ROOT=$1
APPLYCONFIGURATION_GEN=$2
OUTPUT_PACKAGE=$3

PKG_ROOT="github.com/gingersnap-project/operator"
APIS_PKG="${PKG_ROOT}/pkg/apis"
APIS_DIR="${PROJECT_ROOT}/pkg/apis"

rm -rf "${OUTPUT_PACKAGE}"

# client-gen only seems to work with source files in /pkg/apis/<kind>/<version> format, so temporarily create structure
mkdir -p "${APIS_DIR}"/cache/v1alpha1
cp "$PROJECT_ROOT"/api/v1alpha1/cache_types.go "${APIS_DIR}"/cache/v1alpha1/
cp "$PROJECT_ROOT"/api/v1alpha1/lazycacherule_types.go "${APIS_DIR}"/cache/v1alpha1/
cp "$PROJECT_ROOT"/api/v1alpha1/eagercacherule_types.go "${APIS_DIR}"/cache/v1alpha1/
cp "$PROJECT_ROOT"/api/v1alpha1/cacheservice.go "${APIS_DIR}"/cache/v1alpha1/
cp "$PROJECT_ROOT"/api/v1alpha1/zz_*.pb.go "${APIS_DIR}"/cache/v1alpha1/

"${APPLYCONFIGURATION_GEN}" --go-header-file hack/boilerplate.go.txt \
  --input-dirs "${APIS_PKG}"/cache/v1alpha1,github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1 \
  --trim-path-prefix=${PKG_ROOT} \
  --output-package ${PKG_ROOT}/"${OUTPUT_PACKAGE}" \
  --output-base ./

# Update the clientset imports to use the actual types
find "${OUTPUT_PACKAGE}" -type f | xargs sed -i "s#${APIS_PKG}/cache/v1alpha1#${PKG_ROOT}/api/v1alpha1#g"

# Fix compilation issue WithOwnerReferences
find "${OUTPUT_PACKAGE}" -type f | xargs sed -i -e "s#...metav1.OwnerReference#...*v1.OwnerReferenceApplyConfiguration#g" \
  -e "s#(b.OwnerReferences, values\[i\])#(b.OwnerReferences, *values\[i\])#g"

# Clean up tmp apis dir
rm -rf "${APIS_DIR}"
