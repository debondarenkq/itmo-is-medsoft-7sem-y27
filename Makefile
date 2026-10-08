.PHONY: test integration check up down smoke desktop-test

test:
	@for service in staff catalog clinical; do (cd services/$$service && go test -race ./...) || exit 1; done

check:
	@for service in staff catalog clinical; do (cd services/$$service && go vet ./... && test -z "$$(gofmt -l .)") || exit 1; done

up:
	./scripts/compose.sh up --build -d --wait

down:
	./scripts/compose.sh down

smoke:
	python3 scripts/smoke.py

integration:
	./scripts/test-integration.sh

desktop-test:
	./scripts/test-desktop.sh
