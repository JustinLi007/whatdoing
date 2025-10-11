target = service-progress
main = cmd/app/
src = $(shell find . -path "./.*" -prune -o -name '*.go')

default: $(target)

$(target): $(src)
	go build -o ./$(target) ./$(main)

run: $(target)
	./$(target)

test: $(target)
	go clean -testcache
	go test ./... -v

dcup:
	docker compose --env-file ./.env --profile service-anime-progress up -d
	docker compose --env-file ../service-progress-consumer/consumer-anime/.env --profile consumer-anime up -d

dcdown:
	docker compose --profile service-anime-progress down
	docker compose --profile consumer-anime down

tidy:
	go mod tidy

clean:
	rm ./$(target) -f

PHONY: clean default run test dcup dcdown
