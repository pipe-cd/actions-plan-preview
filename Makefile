.PHONY: build
build:
	go build -o plan-preview .

.PHONY: test
test:
	go test ./...

.PHONY: dep
dep:
	go mod tidy
