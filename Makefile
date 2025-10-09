services = gateway users anime
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

logs:
	@docker compose \
	-f ./compose.gateway.yaml \
	-f ./compose.db.yaml \
	-f ./compose.service.yaml \
	logs

clean:
	@for v in $(subdirs); do \
		$(MAKE) --no-print-directory -C $$v clean; \
	done

.PHONY: all clean $(subdirs) up down logs
