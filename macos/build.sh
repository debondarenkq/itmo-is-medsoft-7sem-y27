#!/bin/sh
set -eu
cd "$(dirname "$0")"
case "${1:-release}" in
    debug|Debug) configuration=Debug ;;
    release|Release) configuration=Release ;;
    *) printf '%s\n' 'Usage: ./build.sh [debug|release] [all|HIS|HISAdmin|LIS]' >&2; exit 1 ;;
esac
case "${2:-all}" in
    all) applications='HIS HISAdmin LIS' ;;
    HIS|HISAdmin|LIS) applications="$2" ;;
    *) printf '%s\n' 'Usage: ./build.sh [debug|release] [all|HIS|HISAdmin|LIS]' >&2; exit 1 ;;
esac
for application in $applications; do
    xcodebuild -project HIS.xcodeproj -scheme "$application" -configuration "$configuration" \
        -destination 'platform=macOS' -derivedDataPath build/DerivedData \
        CONFIGURATION_BUILD_DIR="$PWD/build/Products/$configuration" build
    rm -rf "build/$application.app"
    ditto "build/Products/$configuration/$application.app" "build/$application.app"
    printf '%s\n' "$PWD/build/$application.app"
done
