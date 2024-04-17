#!/bin/sh

set -e

make deploy/e2e

sleep 5

kubectl wait \
  --namespace=qdrant-system \
  --for=condition=ready \
  pod \
  --selector=control-plane=qdrant-operator \
  --timeout=90s