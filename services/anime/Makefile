target = service-anime
main = cmd/app/

src = $(shell find . -path './.*' -prune -o -name '*.go')

default: $(target)

$(target): $(src)
	go build -o ./$(target) ./$(main)

run-service: $(target)
	./$(target) "service"

run-pub: 
	./$(target) "pub"

test: $(target)
	go clean -testcache
	go test ./...

dcup-service:
	docker compose up app -d

dcup-pub:
	docker compose up pub -d

dcdown:
	docker compose down

clean:
	rm ./$(target) -f

PHONY: clean default test dcup-service dcup-pub dcdown run-service run-pub
