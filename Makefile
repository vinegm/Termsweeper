BINARY := termsweeper

.PHONY: build
build:
	go build -o $(BINARY) ./

.PHONY: run
run: build
	./$(BINARY)

.PHONY: have-lint
have-lint:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "error: golangci-lint is not installed"; \
		echo "install: https://golangci-lint.run/welcome/install/"; \
		exit 1; \
	fi

.PHONY: fmt
fmt: have-lint
	golangci-lint fmt

.PHONY: fmt-check
fmt-check: have-lint
	golangci-lint fmt --diff

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: have-lint
	golangci-lint run

.PHONY: ci
ci: fmt-check vet lint

.PHONY: have-vhs
have-vhs:
	@if ! command -v vhs >/dev/null 2>&1; then \
		echo "error: vhs is not installed"; \
		echo "install: https://github.com/charmbracelet/vhs"; \
		exit 1; \
	fi

.PHONY: demo
demo: have-vhs
	vhs vhs/demo.tape

.PHONY: clean
clean:
	rm -f $(BINARY)
