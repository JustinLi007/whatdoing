services = gateway users anime progress
subdirs = $(patsubst %, services/%, $(services))

all: $(subdirs)

$(subdirs):
	@$(MAKE) --no-print-directory -C $@

up:
	@docker compose \
	-f ./compose.msgbroker.yaml \
	-f ./compose.gateway.yaml \
	-f ./compose.db.yaml \
	-f ./compose.service.yaml \
	up -d

down:
	@docker compose \
	-f ./compose.msgbroker.yaml \
	-f ./compose.gateway.yaml \
	-f ./compose.db.yaml \
	-f ./compose.service.yaml \
	down

rmi:
	@docker rmi service-user-app:latest
	@docker rmi service-anime-app:latest
	@docker rmi service-anime-progress-app:latest

logs:
	@docker compose \
	-f ./compose.gateway.yaml \
	-f ./compose.db.yaml \
	-f ./compose.service.yaml \
	logs

tidy:
	@for v in $(subdirs); do \
		$(MAKE) --no-print-directory -C $$v tidy; \
	done

clean:
	@for v in $(subdirs); do \
		$(MAKE) --no-print-directory -C $$v clean; \
	done

.PHONY: all clean $(subdirs) up down logs
