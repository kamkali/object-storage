build-container:
	docker build -t object-storage:latest .

run-container:
	docker run -it --name object-storage object-storage:latest

tools:
	cat tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %

generate:
	go generate ./...

test:
	go test -short -race ./...

lint:
	golangci-lint run ./...

fmt-code:
	gofmt -w .

fmt-imports:
	goimports -w .

format: fmt-code fmt-imports

.PHONY: build-container run-container tools generate lint format