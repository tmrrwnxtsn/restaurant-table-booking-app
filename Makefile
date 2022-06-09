.PHONY: build
build: ## сборка бинарника API сервера
	go build -o server cmd/server/main.go

.PHONY: run
run: build ## запуск API сервера
	./server