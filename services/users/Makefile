target = service-users
main = cmd/app/
# src = $(wildcard **/*.go)
src = $(shell find . -path ./.database -prune -o -name '*.go')

default: $(target)

$(target): $(src)
	go build -o ./$(target) ./$(main)

run: $(target)
	./$(target)

test: $(target)
	go clean -testcache
	go test ./internal/setuptest/
	go test ./... $(ARGS)
	go clean -testcache
	go test ./internal/setuptest/

dcup:
	docker compose up -d

dcdown:
	docker compose down

tidy:
	go mod tidy

clean:
	rm ./$(target) -f

PHONY: clean default run test dcup dcdown
