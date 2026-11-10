.DEFAULT_GOAL := help

.PHONY: help test integration check certs prepare prepare-local up demo demo-worker demo-admin down status doctor cert-info images-save images-load smoke macos-build generate-contracts reset reset-demo demo-accounts

test:
	@for service in identity staff catalog clinical; do (cd services/$$service && go test -race ./...) || exit 1; done

check:
	@for service in identity staff catalog clinical; do (cd services/$$service && go vet ./... && test -z "$$(gofmt -l .)") || exit 1; done

help:
	@printf '%s\n' \
	  'make prepare       — первоначальная подготовка с интернетом: ключи, образы и приложение' \
	  'make prepare-local — пересборка на подготовленном ноутбуке без интернета (нужен Go)' \
	  'make demo          — запуск backend и обоих приложений без сборки и скачиваний' \
	  'make demo-worker   — запуск backend и рабочего места медработника' \
	  'make demo-admin    — запуск backend и приложения администратора' \
	  'make up / down     — запуск / остановка backend с сохранением данных' \
	  'make reset         — очистить данные и заполнить демо-набором, без аккаунтов' \
	  'make demo-accounts — создать демо-аккаунты через API (нужен вход администратора)' \
	  'make status        — состояние Docker-компонентов' \
	  'make doctor        — проверка готовности, HTTPS и защиты API' \
	  'make cert-info     — сведения о публичных сертификатах' \
	  'make images-save   — резервная копия Docker-образов в build/his-images.tar' \
	  'make images-load   — загрузка образов из build/his-images.tar'

certs:
	./scripts/setup-security.sh

prepare: certs
	./scripts/compose.sh pull postgres nginx
	./scripts/compose.sh build
	./macos/build.sh

prepare-local: certs
	./scripts/build-local.sh
	./macos/build.sh

up: certs
	./scripts/compose.sh up --no-build --pull never -d --wait

demo:
	@test -x macos/build/HIS.app/Contents/MacOS/HIS && test -x macos/build/HISAdmin.app/Contents/MacOS/HISAdmin || { echo 'Сначала выполните make prepare или make prepare-local.'; exit 1; }
	$(MAKE) up
	open macos/build/HIS.app macos/build/HISAdmin.app

demo-worker:
	@test -x macos/build/HIS.app/Contents/MacOS/HIS || { echo 'Сначала выполните make macos-build.'; exit 1; }
	$(MAKE) up
	open macos/build/HIS.app

demo-admin:
	@test -x macos/build/HISAdmin.app/Contents/MacOS/HISAdmin || { echo 'Сначала выполните make macos-build.'; exit 1; }
	$(MAKE) up
	open macos/build/HISAdmin.app

status:
	@./scripts/docker.sh context show
	./scripts/compose.sh ps

doctor:
	python3 scripts/check-demo.py

cert-info:
	openssl x509 -in deploy/security/tls/ca/ca.crt -noout -subject -issuer -dates -fingerprint -sha256
	openssl x509 -in deploy/security/tls/server/server.crt -noout -subject -issuer -dates -fingerprint -sha256
	openssl verify -CAfile deploy/security/tls/ca/ca.crt -purpose sslserver deploy/security/tls/server/server.crt

images-save:
	mkdir -p build
	./scripts/docker.sh image save -o build/his-images.tar his-identity:local his-staff:local his-catalog:local his-clinical:local postgres:17-alpine nginx:1.28-alpine

images-load:
	./scripts/docker.sh image load -i build/his-images.tar

down:
	./scripts/compose.sh down

smoke:
	python3 scripts/smoke.py

integration:
	./scripts/test-integration.sh

macos-build:
	./macos/build.sh

generate-contracts:
	./scripts/generate-contracts.sh

reset reset-demo:
	python3 scripts/reset-demo.py

demo-accounts:
	python3 scripts/create-demo-accounts.py
