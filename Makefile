
check-gow:
	@command -v gow >/dev/null 2>&1 || { echo >&2 "gow is required but it's not installed. Aborting."; exit 1; }

dev-unit: check-gow
	gow test -failfast -run Unit_ ./...

dev-integration: check-gow
	gow test -failfast -run Integration_ ./...

test:
	go test -failfast ./...

install:
	cd cmd/gho && go install