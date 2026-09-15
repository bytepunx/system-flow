#!/usr/bin/env sh
# Build bin/flai from flai/.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "building bin/flai"
go build -o "$ROOT/bin/flai" .
