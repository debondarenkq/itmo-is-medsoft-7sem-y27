#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
if command -v flutter >/dev/null 2>&1; then
    cd desktop
    flutter pub get
    dart format --output=none --set-exit-if-changed lib test
    flutter analyze
    flutter test
else
    docker run --rm -e CI=true -e PUB_CACHE=/pub-cache \
        -v his-flutter-pub-cache:/pub-cache \
        -v "$PWD/desktop:/app" -w /app \
        ghcr.io/cirruslabs/flutter@sha256:46691e311715845de03a3ba4753a475476936805b29431b1f00f1816981033f8 \
        bash -lc 'flutter pub get && dart format --output=none --set-exit-if-changed lib test && flutter analyze && flutter test'
fi
