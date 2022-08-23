#!/usr/bin/env bash
set -e

CATALOG_DIR=olm-catalog
DOCKERFILE=${CATALOG_DIR}.Dockerfile
CATALOG=${CATALOG_DIR}/catalog.yaml

BUNDLE_IMGS="${BUNDLE_IMG}"
# Define existing bundle images required in the catalog
#for version in v2.2.1 v2.2.2 v2.2.3 v2.2.4; do
#  BUNDLE_IMGS="${BUNDLE_IMGS} quay.io/operatorhubio/engytita:$version"
#done

rm -rf ${CATALOG_DIR}
mkdir ${CATALOG_DIR}

# Define OLM update graph
cat <<EOF >> ${CATALOG}
---
schema: olm.package
name: engytita
defaultChannel: alpha
---
schema: olm.channel
name: alpha
package: engytita
entries:
- name: engytita.v0.0.1
EOF

set -x

${OPM} render --use-http -o yaml ${BUNDLE_IMGS} >> ${CATALOG}

${OPM} validate ${CATALOG_DIR}
${OPM} generate dockerfile ${CATALOG_DIR}
docker build -f ${DOCKERFILE} -t ${CATALOG_IMG} .

rm -rf ${DOCKERFILE}
