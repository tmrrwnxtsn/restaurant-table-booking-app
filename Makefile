# строка подключения к БД
APP_DSN ?= postgres://127.0.0.1/aero?sslmode=disable&user=postgres&password=qwerty

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build: ## сборка бинарника API сервера
	go build -o apiserver cmd/apiserver/main.go

.PHONY: run
run: build ## запуск API сервера
	./apiserver

.PHONY: migrate-up
migrate-up: ## применение миграции к БД
	echo "Running database migration..."
	@migrate -path ./migrations -database "$(APP_DSN)" up 1

.PHONY: migrate-down
migrate-down: ## откат миграций БД на 1 шаг
	echo "Reverting database to the last migration step..."
	@migrate -path ./migrations -database "$(APP_DSN)" down 1

.PHONY: testdata
testdata: ## заполнить БД тестовыми данными
	echo "Filling database with test data..."
	psql -a -f ./testdata/testdata.sql "$(APP_DSN)"

.PHONY: swag-init
swag-init: ## парсинг комментариев у методов и генерация Swagger-документации
	swag init -g cmd/apiserver/main.go

.PHONY: swag-fmt
swag-fmt: ## форматирование комментариев swag
	swag fmt -g cmd/apiserver/main.go

.PHONY: compose-up
compose-up: ## собирает образы API и БД при необходимости и запускает контейнеры (API сервер и Postgres БД)
	docker-compose up

.PHONY: dump
dump: ## делает дамп текущего состояния БД
	pg_dump -U postgres -W -E UTF8 -d aero -f dump.sql

.DEFAULT_GOAL := run