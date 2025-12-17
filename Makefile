# ----------------------------------------------------------------------
# Configuration:
# ----------------------------------------------------------------------

include .docker/.env
include .docker/pgsql/.env

project := ${COMPOSE_PROJECT_NAME}

args := --env-file .docker/.env \
    --env-file .docker/traefik/.env \
    --env-file .docker/nginx/.env \
    --env-file .docker/pgsql/.env \
    --env-file .docker/redis/.env \
   	--env-file .docker/go/.env \
    --env-file .docker/centrifugo/.env

include helpers.mk

pgsql_conn := postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@pgsql:5432/${POSTGRES_DB}?sslmode=disable

## ----------------------------------------------------------------------
## Environment:
## ----------------------------------------------------------------------

.PHONY: setup
setup: ## Environment setup
	make down
	make build
	make up
	make app.migrate.fresh
	make app.ci
	make stop

.PHONY: up
up: ## Environment up
	@$(call compose.use,up,-d --force-recreate --remove-orphans)

.PHONY: stop
stop: ## Environment stop
	@$(call compose.use,stop)

.PHONY: restart
restart: ## Environment restart
	make stop
	make up

.PHONY: down
down: ## Environment down
	make stop
	@$(call compose.use,down,--remove-orphans --volumes)

.PHONY: build
build: ## No cache building containers
	@$(call compose.use,build,--no-cache)

.PHONY: logs
logs: ## Show containers logs
	@$(call compose.use,logs,-f)

## ----------------------------------------------------------------------
## Orchestrator:
## ----------------------------------------------------------------------

.PHONY: server
server: ## Run server container
	@$(call compose.use,up,-d --force-recreate --remove-orphans,server)

.PHONY: server.stop
server.stop: ## Stop server container
	@$(call compose.use,stop,server)

.PHONY: server.logs
server.logs: ## Show server logs
	@$(call compose.use,logs,-f -n 0 server)

.PHONY: scheduler
scheduler: ## Run scheduler container
	@$(call compose.use,up,-d --force-recreate --remove-orphans,scheduler)

.PHONY: scheduler.stop
scheduler.stop: ## Stop scheduler container
	@$(call compose.use,stop,scheduler)

.PHONY: scheduler.logs
scheduler.logs: ## Show scheduler logs
	@$(call compose.use,logs,-f -n 0 scheduler)

.PHONY: centrifugo
centrifugo: ## Run centrifugo container
	@$(call compose.use,up,-d --force-recreate --remove-orphans,centrifugo)

.PHONY: centrifugo.stop
centrifugo.stop: ## Stop centrifugo container
	@$(call compose.use,stop,centrifugo)

.PHONY: centrifugo.logs
centrifugo.logs: ## Show centrifugo logs
	@$(call compose.use,logs,-f centrifugo)

.PHONY: pgsql
pgsql: ## Run pgsql container
	@$(call compose.use,up,-d --force-recreate --remove-orphans,pgsql)

.PHONY: pgsql.stop
pgsql.stop: ## Stop pgsql container
	@$(call compose.use,stop,pgsql)

.PHONY: pgsql.logs
pgsql.logs: ## Show pgsql logs
	@$(call compose.use,logs,-f pgsql)

.PHONY: redis
redis: ## Run redis container
	@$(call compose.use,up,-d --force-recreate --remove-orphans,redis)

.PHONY: redis.stop
redis.stop: ## Stop redis container
	@$(call compose.use,stop,redis)

.PHONY: redis.logs
redis.logs: ## Show redis logs
	@$(call compose.use,logs,-f redis)

.PHONY: nginx
nginx: ## Run nginx container
	@$(call compose.use,up,-d --force-recreate --remove-orphans,nginx)

.PHONY: nginx.stop
nginx.stop: ## Stop nginx container
	@$(call compose.use,stop,nginx)

.PHONY: nginx.logs
nginx.logs: ## Show nginx logs
	@$(call compose.use,logs,-f nginx)

.PHONY: traefik
traefik: ## Run traefik container
	@$(call compose.use,up,-d --force-recreate --remove-orphans,traefik)

.PHONY: traefik.stop
traefik.stop: ## Stop traefik container
	@$(call compose.use,stop,traefik)

.PHONY: traefik.logs
traefik.logs: ## Show traefik logs
	@$(call compose.use,logs,-f traefik)

## ----------------------------------------------------------------------
## App:
## ----------------------------------------------------------------------

.PHONY: app.shell
app.shell: ## Run app shell
	@$(call compose.exec,server,sh)

.PHONY: app.tidy
app.tidy: ## Run app tidy
	@$(call compose.exec,server,go mod tidy)

.PHONY: app.test.static
app.test.static: ## Run app static
	@$(call compose.exec,server,golangci-lint --color always -v run ./...)

.PHONY: app.test.unit
app.test.unit: ## Run app test
	@$(call compose.exec,server,go test -v -count=1 ./...)

.PHONY: app.ci
app.ci: ## Run ci tests
	make app.test.static
	@$(call common.info, "")
	@$(call common.info, "GO static analyze done")
	@$(call common.info, "")
	make app.test.unit
	@$(call common.info, "")
	@$(call common.info, "GO unit tests done")
	@$(call common.info, "")

.PHONY: app.migrate.up
app.migrate.up: ## Run app migrate up
	@$(call compose.exec,server,migrate -source file://db/migrations -database "$(pgsql_conn)" up)

.PHONY: app.migrate.down
app.migrate.down: ## Run app migrate down
	@$(call compose.exec,server,migrate -source file://db/migrations -database "$(pgsql_conn)" down -all)

.PHONY: app.migrate.seed
app.migrate.seed: ## Run app seed
	@$(call compose.exec,server,go run cmd/commands/fixture/main.go)

.PHONY: app.migrate.fresh
app.migrate.fresh: ## Run app migrate fresh
	make app.migrate.down
	make app.migrate.up
	make app.migrate.seed
