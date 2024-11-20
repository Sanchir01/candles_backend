PHONY:

SILENT:

MIGRATION_NAME ?= new_migration

build:
	go build -o ./.bin/main ./cmd/main/main.go


generate-dataloaders:
	(cd internal/feature/color  && dataloaden  LoaderByID LoaderByIdColor \*github.com/Sanchir01/colors/pkg/feature/color.Color)

gql:
	go get github.com/99designs/gqlgen@latest && go run github.com/99designs/gqlgen generate

run: build
	./.bin/main

migrations-up:
	goose -dir migrations postgres "host=localhost user=postgres password=postgres port=5435 dbname=test sslmode=disable"  up

migrations-down:
	goose -dir migrations postgres  "host=localhost user=postgres password=postgres port=5435 dbname=test sslmode=disable"  down


migrations-status:
	goose -dir migrations postgres  "host=localhost user=postgres password=postgres port=5435 dbname=test sslmode=disable" status

migrations-new:
	goose -dir migrations create $(MIGRATION_NAME) sql

migrations-up-prod:
	goose -dir migrations postgres "host=92.118.114.96 user=gen_user password=lzGFBsM~#Z%8Qv port=5432 dbname=default_db"  up

migrations-down-prod:
	goose -dir migrations postgres  "host=92.118.114.96 user=gen_user password=lzGFBsM~#Z%8Qv port=5432 dbname=default_db"  down


migrations-status-prod:
	goose -dir migrations postgres  "host=92.118.114.96 user=gen_user password=lzGFBsM~#Z%8Qv port=5432 dbname=default_db" status

docker-build:
	docker build -t candles .

docker:
	docker-compose  up -d

docker-app: docker-build docker

seed:
	go run cmd/seed/main.go