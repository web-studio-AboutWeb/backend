.PHONY:docs
docs:
	swag init -o web/static/apidocs --ot json -q -g cmd/app/main.go

.PHONY:run
run: docs
	go build -o app cmd/app/main.go && ./app -config-path config.default.yml

mcreate:
	migrate create -ext sql -dir migrations -seq $(name) # $ name=<migration_name> make mcreate

migrate-up:
	go run cmd/migrate/migrate.go up

migrate-down:
	go run cmd/migrate/migrate.go down

migrate-drop:
	go run cmd/migrate/migrate.go drop

test:
	go test -v ./internal/...

test-s:
	go test ./internal/...

test-r:
	go test -v -race ./internal/...

gen:
	go generate ./...