target = service-gateway
main = cmd/app
src = $(shell find . -path './.*' -prune -o -name '*.go')

default: $(target)

$(target): $(src)
	go build -o ./$(target) ./$(main)

run: $(target)
	./$(target)

test: $(target)
	go clean -testcache
	go test ./... $(ARGS)

dcup:
	docker compose up -d

dcdown:
	docker compose down

clean:
	rm ./$(target) -f

PHONY: clean default run test dcup dcdown
