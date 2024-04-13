
dev:
	@command -v gow >/dev/null 2>&1 || { echo >&2 "gow is required but it's not installed. Aborting."; exit 1; }
	gow test ./pkg/utils/...


install:
	cd cmd/gho && go install