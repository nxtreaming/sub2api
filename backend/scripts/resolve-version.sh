#!/usr/bin/env sh
set -eu

# Resolve the build version in the same way across local builds and Docker builds.
# Prefer an exact git tag when available; otherwise fall back to cmd/server/VERSION.
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
BACKEND_DIR=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
VERSION_FILE="$BACKEND_DIR/cmd/server/VERSION"

if command -v git >/dev/null 2>&1; then
    if tag=$(git -C "$BACKEND_DIR" describe --tags --exact-match 2>/dev/null); then
        tag=${tag#v}
        if [ -n "$tag" ]; then
            printf '%s\n' "$tag"
            exit 0
        fi
    fi
fi

tr -d '\r\n' < "$VERSION_FILE"
printf '\n'
