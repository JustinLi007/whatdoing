target = service-progress
main = cmd/app/
src = $(shell find . -path './.*' -prune -o -name '*.go')

default: $(target)

$(target): $(src)
	go build -o ./$(target) ./$(main)

run-service: $(target)
	./$(target) "service"

run-sub: 
	./$(target) "sub"

test: $(target)
	go clean -testcache
	go test ./...

dcup-service:
	docker compose up app -d

dcup-sub:
	docker compose up sub -d

dcdown:
	docker compose down

tidy:
	go mod tidy

clean:
	rm ./$(target) -f

PHONY: clean default test dcup-service dcup-sub dcdown run-service run-sub
