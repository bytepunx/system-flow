#!/usr/bin/env sh
# Build the flaiover image locally as flaiover:local (or $1), from the repo root,
# with the flai version metadata of this checkout.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
TAG="${1:-flaiover:local}"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"
VERSION="$(git describe --tags --match 'flai/v*' --abbrev=0 2>/dev/null | sed 's|flai/v||' || echo dev)"
FO_VERSION="$(git describe --tags --match 'flaiover/v*' --abbrev=0 2>/dev/null | sed 's|flaiover/v||' || echo 0.0.0)"
docker build -f flaiover/Dockerfile -t "$TAG" \
  --build-arg "FLAI_VERSION=${VERSION:-dev}" --build-arg "FLAI_COMMIT=$COMMIT" --build-arg "FLAIOVER_VERSION=${FO_VERSION:-0.0.0}" \
  --build-arg "FLAI_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)" .
echo "built $TAG"
