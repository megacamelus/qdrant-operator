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
  --input-dirs=github.com/megacamelus/qdrant-operator/api/qdrant/v1alpha1 \
  --output-package=github.com/megacamelus/qdrant-operator/pkg/client/qdrant/applyconfiguration

"${PROJECT_ROOT}"/bin/client-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-base="${TMP_DIR}/client" \
  --input=qdrant/v1alpha1 \
  --clientset-name "versioned" \
  --input-base=github.com/megacamelus/qdrant-operator/api \
  --apply-configuration-package=github.com/megacamelus/qdrant-operator/pkg/client/qdrant/applyconfiguration \
  --output-package=github.com/megacamelus/qdrant-operator/pkg/client/qdrant/clientset

"${PROJECT_ROOT}"/bin/lister-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-base="${TMP_DIR}/client" \
  --input-dirs=github.com/megacamelus/qdrant-operator/api/qdrant/v1alpha1 \
  --output-package=github.com/megacamelus/qdrant-operator/pkg/client/qdrant/listers

"${PROJECT_ROOT}"/bin/informer-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-base="${TMP_DIR}/client" \
  --input-dirs=github.com/megacamelus/qdrant-operator/api/qdrant/v1alpha1 \
  --versioned-clientset-package=github.com/megacamelus/qdrant-operator/pkg/client/qdrant/clientset/versioned \
  --listers-package=github.com/megacamelus/qdrant-operator/pkg/client/qdrant/listers \
  --output-package=github.com/megacamelus/qdrant-operator/pkg/client/qdrant/informers


# This should not be needed but for some reasons, the applyconfiguration-gen tool
# sets a wrong APIVersion for the Dapr type (operator/v1alpha1 instead of the one with
# the domain qdrant.megacamelus.github.io/v1alpha1).
#
# See: https://github.com/kubernetes/code-generator/issues/150
sed -i \
  's/WithAPIVersion(\"qdrant\/v1alpha1\")/WithAPIVersion(\"qdrant.megacamelus.github.io\/v1alpha1\")/g' \
  "${TMP_DIR}"/client/github.com/megacamelus/qdrant-operator/pkg/client/qdrant/applyconfiguration/qdrant/v1alpha1/cluster.go
sed -i \
  's/WithAPIVersion(\"qdrant\/v1alpha1\")/WithAPIVersion(\"qdrant.megacamelus.github.io\/v1alpha1\")/g' \
  "${TMP_DIR}"/client/github.com/megacamelus/qdrant-operator/pkg/client/qdrant/applyconfiguration/qdrant/v1alpha1/collection.go


cp -R \
  "${TMP_DIR}"/client/github.com/megacamelus/qdrant-operator/pkg/client/qdrant/* \
  "${PROJECT_ROOT}"/pkg/client/qdrant