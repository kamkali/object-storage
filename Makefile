run:
	docker compose up --build

clean:
	docker compose down

tools:
	cat tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %

generate:
	go generate ./...

test:
	go test -short -race ./...

itest:
	docker-compose -f ./docker-compose.test.yml up --build

lint:
	golangci-lint run ./...

fmt-code:
	gofmt -w .

fmt-imports:
	goimports -w .

format: fmt-code fmt-imports

.PHONY: run clean tools generate test itest lint format