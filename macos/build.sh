#!/bin/sh
set -eu
cd "$(dirname "$0")"
case "${1:-release}" in
    debug|Debug) configuration=Debug ;;
    release|Release) configuration=Release ;;
    *) printf '%s\n' 'Usage: ./build.sh [debug|release]' >&2; exit 1 ;;
esac
xcodebuild -project HIS.xcodeproj -scheme HIS -configuration "$configuration" \
    -destination 'platform=macOS' -derivedDataPath build/DerivedData \
    CONFIGURATION_BUILD_DIR="$PWD/build/Products/$configuration" build
rm -rf build/HIS.app
ditto "build/Products/$configuration/HIS.app" build/HIS.app
printf '%s\n' "$PWD/build/HIS.app"
