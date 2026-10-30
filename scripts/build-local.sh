#!/bin/sh
# Fallback for a prepared laptop when Docker Hub is unavailable.
set -eu
cd "$(dirname "$0")/.."

command -v go >/dev/null 2>&1 || { printf '%s\n' 'Установите Go 1.26 или новее.' >&2; exit 1; }
architecture=$(./scripts/docker.sh info --format '{{.Architecture}}')
case "$architecture" in
    arm64|aarch64) target_arch=arm64 ;;
    amd64|x86_64) target_arch=amd64 ;;
    *) printf 'Неподдерживаемая архитектура Docker: %s\n' "$architecture" >&2; exit 1 ;;
esac

# One shared, stable runtime base for all services. Reuse local layers without registry access.
runtime_image=his-runtime:local
if ! ./scripts/docker.sh image inspect "$runtime_image" >/dev/null 2>&1; then
    if ./scripts/docker.sh image inspect his-staff-runtime:local >/dev/null 2>&1; then
        ./scripts/docker.sh tag his-staff-runtime:local "$runtime_image"
    elif ./scripts/docker.sh image inspect his-staff:local >/dev/null 2>&1; then
        ./scripts/docker.sh tag his-staff:local "$runtime_image"
    else
        printf '%s\n' 'Нет подготовленного runtime. Первый раз выполните make prepare с доступом в интернет.' >&2
        exit 1
    fi
fi
image_arch=$(./scripts/docker.sh image inspect "$runtime_image" --format '{{.Architecture}}')
if [ "$image_arch" != "$target_arch" ]; then
    printf 'Архитектура runtime — %s, Docker — %s. Выполните make prepare на этом компьютере.\n' "$image_arch" "$target_arch" >&2
    exit 1
fi

work_dir=$(mktemp -d "${TMPDIR:-/tmp}/his-local-build.XXXXXX")
trap 'rm -rf "$work_dir"' EXIT HUP INT TERM
for service in identity staff catalog clinical; do
    mkdir -p "$work_dir/$service"
    printf 'Локальная сборка %s для linux/%s\n' "$service" "$target_arch"
    if ! (
        cd "services/$service"
        env CGO_ENABLED=0 GOOS=linux GOARCH="$target_arch" GOPROXY=off GOSUMDB=off GOTOOLCHAIN=local \
            go build -trimpath -o "$work_dir/$service/service" "./cmd/${service}-service"
    ); then
        printf '%s\n' 'Локальная сборка не выполнена. Нужны установленный Go и ранее загруженные зависимости; при доступном интернете используйте make prepare.' >&2
        exit 1
    fi
    ./scripts/docker.sh build --pull=false --network=none \
        --build-arg "RUNTIME_IMAGE=$runtime_image" \
        -f deploy/Dockerfile.local -t "his-${service}:local" "$work_dir/$service"
done
printf '%s\n' 'Образы обновлены из текущих исходников без загрузки зависимостей и образов из интернета.'
