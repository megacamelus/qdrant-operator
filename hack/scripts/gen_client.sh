#!/bin/sh

if [ $# -ne 1 ]; then
    echo "project root is expected"
fi

PROJECT_ROOT="$1"
TMP_DIR=$( mktemp -d -t qdrant-client-gen-XXXXXXXX )

mkdir -p "${TMP_DIR}/client"
mkdir -p "${PROJECT_ROOT}/pkg/client/qdrant"

"${PROJECT_ROOT}"/bin/applyconfiguration-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-base="${TMP_DIR}/client" \
  --input-dirs=github.com/lburgazzoli/qdrant-operator/api/qdrant/v1alpha1 \
  --output-package=github.com/lburgazzoli/qdrant-operator/pkg/client/qdrant/applyconfiguration

"${PROJECT_ROOT}"/bin/client-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-base="${TMP_DIR}/client" \
  --plural-exceptions="qdrant:qdrants" \
  --input=designer/v1alpha1 \
  --clientset-name "versioned" \
  --input-base=github.com/lburgazzoli/qdrant-operator/api \
  --apply-configuration-package=github.com/lburgazzoli/qdrant-operator/pkg/client/qdrant/applyconfiguration \
  --output-package=github.com/lburgazzoli/qdrant-operator/pkg/client/qdrant/clientset

"${PROJECT_ROOT}"/bin/lister-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-base="${TMP_DIR}/client" \
  --plural-exceptions="qdrant:qdrants" \
  --input-dirs=github.com/lburgazzoli/qdrant-operator/api/qdrant/v1alpha1 \
  --output-package=github.com/lburgazzoli/qdrant-operator/pkg/client/qdrant/listers

"${PROJECT_ROOT}"/bin/informer-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-base="${TMP_DIR}/client" \
  --plural-exceptions="qdrant:qdrants" \
  --input-dirs=github.com/lburgazzoli/qdrant-operator/api/qdrant/v1alpha1 \
  --versioned-clientset-package=github.com/lburgazzoli/qdrant-operator/pkg/client/qdrant/clientset/versioned \
  --listers-package=github.com/lburgazzoli/qdrant-operator/pkg/client/qdrant/listers \
  --output-package=github.com/lburgazzoli/qdrant-operator/pkg/client/qdrant/informers

cp -R "${TMP_DIR}"/client/github.com/lburgazzoli/qdrant-operator/pkg/client/qdrant/* "${PROJECT_ROOT}"/pkg/client/qdrant