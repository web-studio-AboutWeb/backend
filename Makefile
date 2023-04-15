DOCKER_COMPOSE := deploy/docker-compose.yml
DOCKER_ENV := deploy/.env
DOCKER_COMPOSE_RUNNER := docker compose
ifneq ($(ENV),)
	DOCKER_COMPOSE := deploy/dev.docker-compose.yml
	DOCKER_ENV := deploy/.env.dev
	DOCKER_COMPOSE_RUNNER := docker compose
	include deploy/.env.dev
	export $(shell sed 's/=.*//' deploy/.env.dev)
endif

run_backend:
	go build -o web-studio cmd/main.go && ./web-studio

compose-build:
	docker compose -f ./deploy/docker-compose.yml --env-file deploy/.env build

compose-up:
	docker compose -f ./deploy/docker-compose.yml --env-file deploy/.env up

docker-rm-volume:
	docker volume rm -f ${PROJECT_NAME}_database_data
