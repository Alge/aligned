.PHONY: fmt test build clean coverage install lint align release help

fmt:
	go fmt ./...

test:
	go test ./...

build:
	go build -o bin/align ./cmd/align

clean:
	rm -rf bin/ coverage.out

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

install:
	go install ./cmd/align

lint:
	golangci-lint run

align:
	./bin/align check spec/

release:
	@if [ -z "$(TYPE)" ]; then \
		echo "Usage: make release TYPE=major|minor|patch"; \
		echo "Example: make release TYPE=patch"; \
		exit 1; \
	fi; \
	CURRENT=$$(cat cmd/align/VERSION); \
	MAJOR=$$(echo $$CURRENT | cut -d. -f1); \
	MINOR=$$(echo $$CURRENT | cut -d. -f2); \
	PATCH=$$(echo $$CURRENT | cut -d. -f3); \
	case "$(TYPE)" in \
		major) NEW="$$((MAJOR+1)).0.0" ;; \
		minor) NEW="$$MAJOR.$$((MINOR+1)).0" ;; \
		patch) NEW="$$MAJOR.$$MINOR.$$((PATCH+1))" ;; \
		*) echo "Invalid TYPE. Use major, minor, or patch"; exit 1 ;; \
	esac; \
	echo "Bumping version $$CURRENT -> $$NEW"; \
	echo "$$NEW" > cmd/align/VERSION; \
	git add cmd/align/VERSION; \
	git commit -m "Bump version to $$NEW"; \
	git tag -a "v$$NEW" -m "Release v$$NEW"; \
	echo ""; \
	echo "Version bumped to v$$NEW and tagged locally."; \
	echo "To push: git push origin master && git push origin v$$NEW"

help:
	@echo "Available targets:"
	@echo "  fmt                  - Format Go code"
	@echo "  test                 - Run tests"
	@echo "  build                - Build align binary"
	@echo "  clean                - Remove build artifacts"
	@echo "  coverage             - Run tests with coverage report"
	@echo "  install              - Install align to GOPATH/bin"
	@echo "  lint                 - Run linter"
	@echo "  align                - Validate specs using align"
	@echo "  release TYPE=<type>  - Create new release (major|minor|patch)"
	@echo "  help                 - Show this help message"